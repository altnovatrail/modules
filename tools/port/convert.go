package main

import (
	"fmt"
	"maps"
	"math/big"
	"regexp"
	"slices"
	"strings"
)

var (
	// notation is an option in the notation of the checks of numbers, as the port writes a key or a
	// trap value: an integer, a decimal or a simple fraction (format/README.md, section 15), and
	// then, after a space or none, a unit, which must be one of the task's units.
	notation = regexp.MustCompile(`^(-?(?:0|[1-9][0-9]*)(?:\.[0-9]+|/[1-9][0-9]*)?) ?([^0-9\s./-].*)?$`)
	// pointsAtOptions is a question that is not asked without its options, "Which of these
	// numbers…": it keeps them, and the learner taps one (docs/architecture.md, section 3.5).
	pointsAtOptions = regexp.MustCompile(`(?i)\bwhich of these\b`)
)

// maxValue is the longest a key or a trap value of a check number or fraction may be written: a
// value read is in notation only when it is at most 16 characters long (format/README.md, section
// 15).
const maxValue = 16

// Converted is a task of the module as the port writes it: the file tasks/<id>.yaml and the
// pictures it shows, by their paths from the module's folder.
type Converted struct {
	ID    string
	Files map[string][]byte
	// Notes are what the port left out, for the person who reads its report.
	Notes []string
}

// convert carries one task of standalone into a task of the module, by the table «что во что»
// (docs/architecture.md, section 3.5).
func convert(ex Example, batch Batch, task BatchTask, overlay Overlay, topicFile string) (Converted, error) {
	skill := idOf(ex.Topic)
	defaults, known := batch.Skills[skill]
	if !known {
		return Converted{}, fmt.Errorf("%s: the batch says nothing of its skill %s", ex.ID, skill)
	}
	files := map[string][]byte{}
	var notes []string
	if ex.Hint != "" {
		notes = append(notes, "its hint is not carried: the format has no field for it")
	}

	if ex.Picture == nil && overlay.PictureAlt != "" {
		return Converted{}, fmt.Errorf("%s: its overlay has picture_alt, and the task has no picture", ex.ID)
	}
	if ex.SolutionPicture == nil && overlay.SolutionPictureAlt != "" {
		return Converted{}, fmt.Errorf("%s: its overlay has solution_picture_alt, and its solution has no picture", ex.ID)
	}
	statement := sequence{mapping{{"text", ex.Question}}}
	if ex.Picture != nil {
		part, png, err := pictureOf(ex.ID, ex.Picture, overlay.PictureAlt, "")
		if err != nil {
			return Converted{}, err
		}
		statement = append(statement, part)
		files[pathOfPicture(part)] = png
	}

	slot, traps, dropped, err := answerOf(ex, task, overlay)
	if err != nil {
		return Converted{}, err
	}
	notes = append(notes, dropped...)
	item := mapping{{"answer_slots", sequence{slot}}}
	if overlay.SimplerCase != "" {
		item = append(item, pair{"simpler_case", overlay.SimplerCase})
	}

	solution := sequence{mapping{{"text", ex.Solution}}}
	if ex.SolutionTotal != "" {
		solution = append(solution, mapping{{"text", ex.SolutionTotal}})
	}
	if ex.SolutionPicture != nil {
		part, png, err := pictureOf(ex.ID, ex.SolutionPicture, overlay.SolutionPictureAlt, "-solution")
		if err != nil {
			return Converted{}, err
		}
		solution = append(solution, part)
		files[pathOfPicture(part)] = png
	}

	root := mapping{{"id", ex.ID}, {"statement", statement}}
	if len(overlay.Givens) > 0 {
		root = append(root, pair{"givens", strings_(overlay.Givens)})
	}
	if overlay.Unknown != "" {
		root = append(root, pair{"unknown", overlay.Unknown})
	}
	canvas := mapping{{"kind", batch.Canvas.Kind}}
	if batch.Canvas.Grid != "" {
		canvas = append(canvas, pair{"grid", batch.Canvas.Grid})
	}
	root = append(root, pair{"canvas", canvas}, pair{"items", sequence{item}}, pair{"solution", solution})
	techniques := defaults.Techniques
	if len(overlay.Techniques) > 0 {
		techniques = overlay.Techniques
	}
	if len(techniques) > 0 {
		root = append(root, pair{"techniques", strings_(techniques)})
	}
	for _, trap := range overlay.Traps {
		if !slices.Contains(traps, trap) {
			traps = append(traps, trap)
		}
	}
	if len(traps) > 0 {
		root = append(root, pair{"traps", strings_(traps)})
	}
	if len(overlay.LeakPhrases) > 0 {
		root = append(root, pair{"leak_phrases", strings_(overlay.LeakPhrases)})
	}
	skills := mapping{{"primary", strings_([]string{skill})}}
	if len(defaults.Supporting) > 0 {
		skills = append(skills, pair{"supporting", strings_(defaults.Supporting)})
	}
	root = append(root,
		pair{"skills", skills},
		pair{"band", "grades-" + ex.GradeLevel},
		pair{"difficulty", ex.Difficulty},
	)

	comment := []string{
		fmt.Sprintf("Written by tools/port from mathtrail-standalone at %s, from", batch.Standalone[:7]),
		topicFile + ", and not to be edited here: what a person adds goes in",
		fmt.Sprintf("tools/port/overlay/%s.yaml, and the port writes this file again (R12).", ex.ID),
	}
	files["tasks/"+ex.ID+".yaml"] = emit(comment, root)
	return Converted{ID: ex.ID, Files: files, Notes: notes}, nil
}

// answerOf is the answer slot of a task, the traps of the wrong options that become no trap value,
// and a note of each such option. A question that points at its options keeps them and is tapped;
// any other answer is written by hand anywhere on the sheet: a simple fraction is checked as a
// fraction, and an integer or a decimal, with units, as a number, which takes a fraction too. An
// answer that is neither, a word, a time or a date, is the rest of the port's to write (E3).
func answerOf(ex Example, task BatchTask, overlay Overlay) (mapping, []string, []string, error) {
	right := ex.Options[ex.CorrectAnswer]
	if pointsAtOptions.MatchString(ex.Question) {
		slot, err := chosen(ex, overlay)
		return slot, nil, nil, err
	}
	kind, key, err := valueOf(right, task.Units)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%s: its answer %q: %w", ex.ID, right, err)
	}
	if kind == "" {
		return nil, nil, nil, fmt.Errorf("%s: its answer %q is no number and no fraction: how such an answer is given is for the rest of the port to decide (docs/architecture.md, section 3.5)", ex.ID, right)
	}
	if len(key) > maxValue {
		return nil, nil, nil, fmt.Errorf("%s: its answer %s is longer than %d characters", ex.ID, key, maxValue)
	}
	check := mapping{{"kind", kind}, {"key", key}}
	if kind == "number" && len(task.Units) > 0 {
		check = append(check, pair{"units", strings_(task.Units)})
	}
	var values sequence
	var traps, dropped []string
	seen := map[string]string{key: ex.CorrectAnswer}
	used := map[string]bool{}
	for _, letter := range letters {
		if letter == ex.CorrectAnswer {
			continue
		}
		distractor := ex.Distractors[letter]
		trap := idOf(distractor.Trap)
		// A wrong option with a unit the task does not count in is no value of the check either: the
		// learner who writes it gets no explanation of it, rather than the port one of another value.
		wrongKind, value, unitErr := valueOf(ex.Options[letter], task.Units)
		if unitErr != nil || wrongKind == "" || (kind == "fraction" && wrongKind != "fraction") {
			// A wrong option that is no value of the check cannot be a trap value: its trap goes
			// to the task's traps, for AI to suggest, and its text is left out (docs/architecture.md,
			// section 3.5).
			if !slices.Contains(traps, trap) {
				traps = append(traps, trap)
			}
			dropped = append(dropped, fmt.Sprintf("its option %s, %q, is no value of the check %s: its trap %s goes to the task's traps, and its text is not carried", letter, ex.Options[letter], kind, trap))
			continue
		}
		if len(value) > maxValue {
			return nil, nil, nil, fmt.Errorf("%s: its option %s, %s, is longer than %d characters", ex.ID, letter, value, maxValue)
		}
		for _, other := range slices.Sorted(maps.Keys(seen)) {
			if same(value, other) {
				return nil, nil, nil, fmt.Errorf("%s: its option %s, %s, is the value of option %s too", ex.ID, letter, value, seen[other])
			}
		}
		seen[value] = letter
		text := distractor.Text
		if rewritten, has := overlay.TrapTexts[value]; has {
			text = rewritten
			used[value] = true
		}
		values = append(values, mapping{{"value", value}, {"trap", trap}, {"text", text}})
	}
	if err := unusedTexts(ex.ID, overlay, used); err != nil {
		return nil, nil, nil, err
	}
	if len(values) > 0 {
		check = append(check, pair{"trap_values", values})
	}
	return mapping{{"binding", "anywhere"}, {"check", check}}, traps, dropped, nil
}

// chosen is the answer slot of a question that points at its options: the five of them, tapped,
// the right one the key and each other a trap value, as written.
func chosen(ex Example, overlay Overlay) (mapping, error) {
	var options, values sequence
	used := map[string]bool{}
	for _, letter := range letters {
		option := ex.Options[letter]
		if slices.Contains(options, any(option)) {
			return nil, fmt.Errorf("%s: its option %s, %q, is another option too", ex.ID, letter, option)
		}
		options = append(options, option)
		if letter == ex.CorrectAnswer {
			continue
		}
		distractor := ex.Distractors[letter]
		text := distractor.Text
		if rewritten, has := overlay.TrapTexts[option]; has {
			text = rewritten
			used[option] = true
		}
		values = append(values, mapping{{"value", option}, {"trap", idOf(distractor.Trap)}, {"text", text}})
	}
	if err := unusedTexts(ex.ID, overlay, used); err != nil {
		return nil, err
	}
	check := mapping{{"kind", "choice"}, {"key", ex.Options[ex.CorrectAnswer]}, {"trap_values", values}}
	return mapping{{"binding", "choice"}, {"options", options}, {"check", check}}, nil
}

// unusedTexts refuses an overlay that rewrites the explanation of a value the task does not have:
// the text would be lost unsaid.
func unusedTexts(id string, overlay Overlay, used map[string]bool) error {
	var unused []string
	for value := range overlay.TrapTexts {
		if !used[value] {
			unused = append(unused, value)
		}
	}
	if len(unused) > 0 {
		slices.Sort(unused)
		return fmt.Errorf("%s: its overlay rewrites the text of %s, which is no trap value of it", id, strings.Join(unused, ", "))
	}
	return nil
}

// valueOf reads an option in the notation of the checks of numbers: "fraction" for a simple
// fraction, "number" for an integer or a decimal, each with the value as it is written, without its
// unit; or no kind for any other option, such as 07, a word or a time. A number with a unit that is
// not among the task's units is an error: as a key, its unit would be lost.
func valueOf(option string, units []string) (kind, value string, err error) {
	m := notation.FindStringSubmatch(option)
	if m == nil {
		return "", "", nil
	}
	kind = "number"
	if strings.Contains(m[1], "/") {
		kind = "fraction"
	}
	if m[2] != "" && !slices.Contains(units, m[2]) {
		if len(units) == 0 {
			return "", "", fmt.Errorf("its unit %q needs the units of the task in the batch", m[2])
		}
		return "", "", fmt.Errorf("its unit %q is not among the units %s", m[2], strings.Join(units, ", "))
	}
	return kind, m[1], nil
}

// same tells whether two values are one value, as the check compares them: 14 and 14.0, or 1/2
// and 2/4. Every value that valueOf gives is a number or a fraction that big.Rat reads.
func same(a, b string) bool {
	x, _ := new(big.Rat).SetString(a)
	y, _ := new(big.Rat).SetString(b)
	return x.Cmp(y) == 0
}

// idOf is a topic or a trap of standalone as an ID of the format: a dot and an underscore become a
// hyphen, counting.gaps → counting-gaps (docs/architecture.md, section 3.5).
func idOf(name string) string {
	return strings.NewReplacer(".", "-", "_", "-").Replace(name)
}

// strings_ is a list of strings as a sequence.
func strings_(list []string) sequence {
	s := make(sequence, len(list))
	for i, item := range list {
		s[i] = item
	}
	return s
}
