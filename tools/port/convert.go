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
	// notation is an option in the notation of the checks of numbers, then, after a space or none,
	// a unit, which must be one of the task's units.
	notation = regexp.MustCompile(`^(-?(?:0|[1-9][0-9]*)(?:\.[0-9]+|/[1-9][0-9]*)?) ?([^0-9\s./-].*)?$`)
	// pointsAtOptions is a question that is not asked without its options, "Which of these
	// numbers…": it keeps them, and the learner taps one.
	pointsAtOptions = regexp.MustCompile(`(?i)\bwhich of these\b`)
	// clockTime and oneWord are the answers that are not numbers and are written by hand; one word
	// only when the question holds it.
	clockTime = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	oneWord   = regexp.MustCompile(`^\p{L}+$`)
	// yesOrNo is tapped, though it is one word, and the question may hold it: written by hand, a
	// learner's "yes, always" would be no option at all.
	yesOrNo = regexp.MustCompile(`(?i)^(yes|no)$`)
	// leadingDigit is an answer that looks like a number: as a choice tapped, the report names it,
	// since the batch may have forgotten its units.
	leadingDigit = regexp.MustCompile(`^-?[0-9]`)
)

// maxValue is the longest a key or a trap value of a check number or fraction may be written: the
// format reads a longer value as out of notation.
const maxValue = 16

// The kinds of answer, as the survey counts them.
const (
	answerNumber   = "number"
	answerFraction = "fraction"
	answerTime     = "time"
	answerByHand   = "choice by hand"
	answerTapped   = "choice tapped"
)

// twelveHoursText is the text of the trap value of twelveHoursOff: it names no number, as it is the
// same in every task.
const twelveHoursText = "That is the time on a clock face, which starts its hours again at noon. Count the hours of the afternoon on from noon."

// Converted is a task of the module as the port writes it: the file tasks/<id>.yaml and the
// pictures it shows, by their paths from the module's folder.
type Converted struct {
	ID    string
	Files map[string][]byte
	// Answer is the kind of its answer, one of the constants answerNumber to answerTapped.
	Answer string
	// Notes are what the port left out, for the person who reads its report.
	Notes []string
	// KeyNamed counts the explanations and the leading question that name the key of a choice or a
	// time, which only a survey writes.
	KeyNamed int
	// FromPicture says that the task's solver reads its picture.
	FromPicture bool
}

// surveyAlt is the alt of every picture of a survey, whose tasks have no overlays: the module of a
// survey is a folder of its own, which no batch writes and the repository never holds.
const surveyAlt = "Survey: a picture of standalone, with no alt written."

// convert carries one task of standalone into a task of the module. It refuses a skill
// the batch does not know, an overlay field that has nothing to stand for, and a key that the
// answer form, the alt or a text the learner sees before solving gives away.
func convert(ex Example, batch Batch, task BatchTask, overlay Overlay, topicFile string) (Converted, error) {
	skill := idOf(ex.Topic)
	defaults, known := batch.Skills[skill]
	if !known {
		return Converted{}, fmt.Errorf("%s: the batch says nothing of its skill %s", ex.ID, skill)
	}
	files := map[string][]byte{}
	var notes []string
	if ex.Hint == "" && overlay.LeadingQuestion != "" {
		return Converted{}, fmt.Errorf("%s: its overlay has leading_question, and the task has no hint", ex.ID)
	}
	leading := leadingOf(ex, overlay)

	if ex.Picture == nil && overlay.PictureAlt != "" {
		return Converted{}, fmt.Errorf("%s: its overlay has picture_alt, and the task has no picture", ex.ID)
	}
	if ex.SolutionPicture == nil && overlay.SolutionPictureAlt != "" {
		return Converted{}, fmt.Errorf("%s: its overlay has solution_picture_alt, and its solution has no picture", ex.ID)
	}
	pictureAlt, solutionAlt := overlay.PictureAlt, overlay.SolutionPictureAlt
	if batch.survey {
		pictureAlt, solutionAlt = surveyAlt, surveyAlt
	}
	statement := sequence{mapping{{"text", ex.Question}}}
	if ex.Picture != nil {
		part, png, err := pictureOf(ex.ID, ex.Picture, pictureAlt, "")
		if err != nil {
			return Converted{}, err
		}
		statement = append(statement, part)
		files[pathOfPicture(part)] = png
	}
	if overlay.AnswerForm != "" {
		statement = append(statement, mapping{{"text", overlay.AnswerForm}})
	}

	given, err := answerOf(ex, task, overlay, batch.survey)
	if err != nil {
		return Converted{}, err
	}
	// The statement shows answer_form, so the leak check of the format would take a key in it for
	// one the task shows, and search for that key nowhere else.
	if names(overlay.AnswerForm, given.key) {
		return Converted{}, fmt.Errorf("%s: its overlay's answer_form names its key %q", ex.ID, given.key)
	}
	written, err := altOf(ex, overlay, given.key, batch.survey)
	if err != nil {
		return Converted{}, err
	}
	if written != "" {
		notes = append(notes, written)
	}
	notes = append(notes, given.notes...)
	traps := given.traps
	item := mapping{{"answer_slots", sequence{given.slot}}}
	if overlay.SimplerCase != "" {
		item = append(item, pair{"simpler_case", overlay.SimplerCase})
	}

	solution := sequence{mapping{{"text", ex.Solution}}}
	if ex.SolutionTotal != "" {
		solution = append(solution, mapping{{"text", ex.SolutionTotal}})
	}
	if ex.SolutionPicture != nil {
		part, png, err := pictureOf(ex.ID, ex.SolutionPicture, solutionAlt, "-solution")
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
	if leading != "" {
		root = append(root, pair{"leading_question", leading})
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

	overlays := batch.overlays
	if overlays == "" {
		overlays = "tools/port/overlay"
	}
	comment := []string{
		fmt.Sprintf("Written by tools/port from mathtrail-standalone at %s, from", batch.Standalone[:7]),
		topicFile + ", and not to be edited here: what a person adds goes in",
		fmt.Sprintf("%s/%s.yaml, and the port writes this file again (R12).", overlays, ex.ID),
	}
	files["tasks/"+ex.ID+".yaml"] = emit(comment, root)
	return Converted{ID: ex.ID, Files: files, Answer: given.kind, Notes: notes, KeyNamed: given.keyNamed, FromPicture: ex.fromPicture}, nil
}

// altOf holds the alt of a task's picture to its rule (R21). The alt of a task solved from the
// picture names every text the picture writes, since its facts stand on the picture alone and a
// learner who cannot see it reads the alt; what the picture draws only proofreading holds. No alt
// names the key unless the question does: the leak check takes a key in the alt for one the task
// shows, and would search for it nowhere else. A survey, whose alts stand in for none, is held to
// neither. A picture that writes the key shows it, and its alt names it: the leak check then
// searches the key nowhere, and the note asks the overlay's leak_phrases to carry what gives it
// away.
func altOf(ex Example, overlay Overlay, key string, survey bool) (string, error) {
	if ex.Picture == nil {
		return "", nil
	}
	written, err := writtenLabels(ex.ID, ex.Picture)
	if err != nil {
		return "", err
	}
	note := ""
	if slices.Contains(written, key) {
		note = fmt.Sprintf("its picture writes its key %q, so the leak check takes the key for shown and searches it nowhere: leak_phrases carry what gives it away", key)
	}
	if survey {
		return note, nil
	}
	if names(overlay.PictureAlt, key) && !names(ex.Question, key) && !slices.Contains(written, key) {
		return "", fmt.Errorf("%s: its picture_alt names its key %q, which neither its question nor its picture does", ex.ID, key)
	}
	if !ex.fromPicture {
		return note, nil
	}
	for _, text := range written {
		found := names(overlay.PictureAlt, text)
		if capitals.MatchString(text) {
			found = namesCased(overlay.PictureAlt, text)
		}
		if !found {
			return "", fmt.Errorf("%s: it is solved from the picture, and its picture_alt does not name %q, which the picture writes", ex.ID, text)
		}
	}
	return note, nil
}

// capitals is a label of Latin capitals, which the alt names as it is written: A, never a. An
// article that opens a sentence, A, still passes for the label, and proofreading holds that.
var capitals = regexp.MustCompile(`^[A-Z]+$`)

// answer is the answer slot of a task, as answerOf gives it. traps are those of the wrong options
// that become no trap value; keyNamed counts the explanations and leading questions that name the
// key of a choice or a time, which only a survey lets through.
type answer struct {
	slot     mapping
	key      string
	traps    []string
	notes    []string
	kind     string
	keyNamed int
}

// answerOf is the answer slot of a task, by the kind of its right option (R20):
//   - a question that points at its options keeps them, and the learner taps one;
//   - an answer in the notation of the checks of numbers is written by hand anywhere on the sheet,
//     a simple fraction checked as a fraction, any other as a number with the batch's units;
//   - a time on a clock is written by hand and checked as a time, with no options, each wrong
//     option a trap value, and a time from 13:00 or before 1:00 one more, as a clock face shows it;
//   - one word other than yes and no that the question holds is a choice written by hand with
//     notation word, so that a word none of the options is a wrong answer, not one left unread. A
//     batch may have such a word tapped, where its wrong options are reasons a learner would not
//     write;
//   - any other answer is a choice tapped, and the report names such an answer that begins with a
//     digit, since the batch may have forgotten its units.
//
// In a survey, lenient, an explanation that names the key of a choice is counted, not refused.
func answerOf(ex Example, task BatchTask, overlay Overlay, lenient bool) (answer, error) {
	right := ex.Options[ex.CorrectAnswer]
	switch {
	case pointsAtOptions.MatchString(ex.Question):
		return chosen(ex, overlay, lenient, tapping)
	case clockTime.MatchString(right):
		// A time reads as a number with words after it, 7 and :30, so it goes before valueOf.
		return chosen(ex, overlay, lenient, timeByHand)
	}
	kind, key, err := valueOf(right, task.Units)
	if err != nil {
		return answer{}, fmt.Errorf("%s: its answer %q: %w", ex.ID, right, err)
	}
	if kind == "" {
		how := tapping
		if !task.Tapped && oneWord.MatchString(right) && !yesOrNo.MatchString(right) && names(ex.Question, right) {
			how = wordByHand
		}
		a, err := chosen(ex, overlay, lenient, how)
		if err == nil && how == tapping && leadingDigit.MatchString(right) {
			a.notes = append(a.notes, fmt.Sprintf("its answer %q is no number in the notation of the checks, or one with words after it that the batch gives no units for: it is a choice tapped", right))
		}
		return a, err
	}
	if len(key) > maxValue {
		return answer{}, fmt.Errorf("%s: its answer %s is longer than %d characters", ex.ID, key, maxValue)
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
			var note string
			traps, note = notCarried(traps, letter, ex.Options[letter], kind, trap)
			dropped = append(dropped, note)
			continue
		}
		if len(value) > maxValue {
			return answer{}, fmt.Errorf("%s: its option %s, %s, is longer than %d characters", ex.ID, letter, value, maxValue)
		}
		for _, other := range slices.Sorted(maps.Keys(seen)) {
			if same(value, other) {
				return answer{}, fmt.Errorf("%s: its option %s, %s, is the value of option %s too", ex.ID, letter, value, seen[other])
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
		return answer{}, err
	}
	if len(values) > 0 {
		check = append(check, pair{"trap_values", values})
	}
	return answer{slot: mapping{{"binding", "anywhere"}, {"check", check}}, key: key, traps: traps, notes: dropped, kind: kind}, nil
}

// writing is how an answer that is no number is given.
type writing int

const (
	tapping writing = iota
	// timeByHand is a time on a clock, checked as a time with no options, so that a time none of
	// the options is a wrong answer.
	timeByHand
	// wordByHand is one word, with notation word, so that a word none of the options is a wrong
	// answer.
	wordByHand
)

// chosen is the answer slot of an answer given among five options: the right one the key and each
// other a trap value. Tapped, the options are the answer slot's; a word written by hand has them in
// the check, which the learner does not see; a time has none. It refuses two options alike as the
// check compares them. The leak check of the format searches neither the key of a tapped choice
// nor one that the statement holds (R20), so the port searches it in the texts the learner
// sees before solving and refuses it there: in an explanation, to be rewritten in the overlay, in
// the leading question, in the simpler case, and in the givens and unknown but where they restate
// a question that names it. A survey, lenient, has no overlay, and writes such a text and counts
// it.
func chosen(ex Example, overlay Overlay, lenient bool, how writing) (answer, error) {
	key := ex.Options[ex.CorrectAnswer]
	namesKey := func(text string) bool { return names(text, key) }
	if how == timeByHand {
		// The key of a time has hours with no leading zero, and the leak check finds it with one
		// too, 07:30 for 7:30.
		key = timeForm(key)
		namesKey = func(text string) bool { return names(text, key) || len(key) == len("7:30") && names(text, "0"+key) }
	}
	var options, values sequence
	var traps, notes []string
	keyNamed := 0
	used := map[string]bool{}
	seen := map[string]string{}
	for _, letter := range letters {
		option := ex.Options[letter]
		// A word's options are compared as notation word reads them, with their end cut, and none
		// may be left blank (task.choice-options); a time's, as the kind time reads them, by hours
		// and minutes.
		form, compared := choiceForm(option), "a choice"
		switch how {
		case wordByHand:
			form = wordForm(option)
			if form == "" {
				return answer{}, fmt.Errorf("%s: its option %s, %q, is blank as notation word reads it", ex.ID, letter, option)
			}
		case timeByHand:
			form, compared = timeForm(option), "the kind time"
			if form == "" {
				var note string
				traps, note = notCarried(traps, letter, option, "time", idOf(ex.Distractors[letter].Trap))
				notes = append(notes, note)
				continue
			}
		}
		if other, twice := seen[form]; twice {
			return answer{}, fmt.Errorf("%s: its option %s, %q, is option %s too, as %s compares them", ex.ID, letter, option, other, compared)
		}
		seen[form] = letter
		if how != timeByHand {
			options = append(options, option)
		}
		if letter == ex.CorrectAnswer {
			continue
		}
		distractor := ex.Distractors[letter]
		text := distractor.Text
		if rewritten, has := overlay.TrapTexts[option]; has {
			text = rewritten
			used[option] = true
		}
		if namesKey(text) {
			if !lenient {
				return answer{}, fmt.Errorf("%s: the explanation of its option %s names its key %q: rewrite it in the trap_texts of its overlay (R05)", ex.ID, letter, key)
			}
			keyNamed++
			notes = append(notes, fmt.Sprintf("the explanation of its option %s names its key %q", letter, key))
		}
		values = append(values, mapping{{"value", option}, {"trap", idOf(distractor.Trap)}, {"text", text}})
	}
	if err := unusedTexts(ex.ID, overlay, used); err != nil {
		return answer{}, err
	}
	// The leading question comes before the first wrong answer and hints more than a technique, so,
	// like the simpler case, it never names the key.
	if leading := leadingOf(ex, overlay); namesKey(leading) {
		if !lenient {
			return answer{}, fmt.Errorf("%s: its leading question names its key %q: rewrite it in the leading_question of its overlay (R12)", ex.ID, key)
		}
		keyNamed++
		notes = append(notes, fmt.Sprintf("its leading question names its key %q", key))
	}
	// The givens and the unknown restate the statement, so they may name a key that the question
	// names too, as a word written by hand always is, but only beside another option or where they
	// quote the question: else the key there is the answer. Only the question counts, not the
	// picture, which is stricter than the leak check's statement. The simpler case never names it.
	restated := append([]string{overlay.Unknown}, overlay.Givens...)
	for i, text := range append(restated, overlay.SimplerCase) {
		if !namesKey(text) || i < len(restated) && namesKey(ex.Question) && (namesAnother(text, ex) || quotes(ex.Question, text, key)) {
			continue
		}
		return answer{}, fmt.Errorf("%s: its overlay's text %q names its key %q", ex.ID, text, key)
	}
	switch how {
	case timeByHand:
		// The kind time counts hours to 24 and leaves the slip of a clock face to a trap value, unless
		// an option is that time already.
		if twelve := twelveHoursOff(key); twelve != "" && seen[twelve] == "" {
			values = append(values, mapping{{"value", twelve}, {"trap", "twelve-hours-off"}, {"text", twelveHoursText}})
		}
		check := mapping{{"kind", "time"}, {"key", key}}
		if len(values) > 0 {
			check = append(check, pair{"trap_values", values})
		}
		return answer{slot: mapping{{"binding", "anywhere"}, {"check", check}}, key: key, traps: traps, notes: notes, kind: answerTime, keyNamed: keyNamed}, nil
	case wordByHand:
		check := mapping{{"kind", "choice"}, {"notation", "word"}, {"key", key}, {"options", options}, {"trap_values", values}}
		return answer{slot: mapping{{"binding", "anywhere"}, {"check", check}}, key: key, notes: notes, kind: answerByHand, keyNamed: keyNamed}, nil
	}
	check := mapping{{"kind", "choice"}, {"key", key}, {"trap_values", values}}
	return answer{slot: mapping{{"binding", "choice"}, {"options", options}, {"check", check}}, key: key, notes: notes, kind: answerTapped, keyNamed: keyNamed}, nil
}

// quotes tells whether a text is words of a question as they stand there, but for the punctuation
// that ends a sentence, and holds at least three words besides those of a key: a fragment as short
// as "a Friday" only states the key.
func quotes(question, text, key string) bool {
	text = strings.TrimRight(text, " .!?;:")
	keyWords := map[string]bool{}
	for _, word := range strings.FieldsFunc(leakForm(key), notWordChar) {
		keyWords[word] = true
	}
	others := 0
	for _, word := range strings.FieldsFunc(leakForm(text), notWordChar) {
		if !keyWords[word] {
			others++
		}
	}
	return others >= 3 && names(question, text)
}

func notWordChar(r rune) bool { return !wordChar(r) }

// namesAnother tells whether a text names an option of a choice other than its key.
func namesAnother(text string, ex Example) bool {
	for _, letter := range letters {
		if letter != ex.CorrectAnswer && names(text, ex.Options[letter]) {
			return true
		}
	}
	return false
}

// notCarried is the rule of a wrong option that is no value of its check: it cannot be a trap
// value, so its trap goes to the task's traps, for AI to suggest, and its text is left out, which
// the note says.
func notCarried(traps []string, letter, option, kind, trap string) ([]string, string) {
	if !slices.Contains(traps, trap) {
		traps = append(traps, trap)
	}
	return traps, fmt.Sprintf("its option %s, %q, is no value of the check %s: its trap %s goes to the task's traps, and its text is not carried", letter, option, kind, trap)
}

// leadingOf is the leading question of a task: the hint of standalone, or the overlay's question
// written in its place.
func leadingOf(ex Example, overlay Overlay) string {
	if overlay.LeadingQuestion != "" {
		return overlay.LeadingQuestion
	}
	return ex.Hint
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

// valueOf reads an option in the notation of the checks of numbers: "fraction" or "number", with
// the value as it is written, without its unit; or no kind for any other option, such as 07, a
// word or a time. A number with words after it is a number only when the batch gives the task
// units: then words that are none of them are an error, since as a key their unit would be lost,
// and without units the option is no number, such as 3 and 6 or 19 May.
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
			return "", "", nil
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

// idOf is a topic or a trap of standalone as an ID of the format, counting.gaps → counting-gaps.
func idOf(name string) string {
	return strings.NewReplacer(".", "-", "_", "-").Replace(name)
}

func strings_(list []string) sequence {
	s := make(sequence, len(list))
	for i, item := range list {
		s[i] = item
	}
	return s
}
