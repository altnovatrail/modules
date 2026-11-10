package main

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// A survey runs the port over every task of standalone that no batch names yet, before the batches
// of E3 are written: which tasks it writes, by the kind of their answer, which wait for a picture of
// a kind it does not draw yet, and which it refuses (R20). It reads no overlays, so it writes into a
// module of its own, which the recipe port-survey makes of the format's examples, and which never
// enters the repository: its pictures have surveyAlt, an explanation that names the key of a choice
// is counted rather than refused, and a number with words after it, which no batch gives units for,
// is a choice tapped.
type survey struct {
	gitCommand, goCommand string
	standalone, commit    string
	except, module        string
	prove                 bool
}

// runSurvey is the survey, as run calls it, with its exit status.
func runSurvey(s survey, stdout, stderr io.Writer) int {
	fail := func(err error) int {
		fmt.Fprintf(stderr, "port: %v\n", err)
		return statusCannotRun
	}
	if !commit.MatchString(s.commit) {
		return fail(fmt.Errorf("--commit %q is not the full SHA of a commit", s.commit))
	}
	examples, topics, err := readClone(s.gitCommand, s.standalone, s.commit)
	if err != nil {
		return fail(err)
	}
	named, err := batchedTasks(s.except)
	if err != nil {
		return fail(err)
	}
	batch := Batch{
		Standalone: s.commit,
		Canvas:     Canvas{Kind: "handwriting_sheet", Grid: "squares"},
		Skills:     map[string]SkillDefaults{},
		survey:     true,
	}
	for id := range topics {
		batch.Skills[idOf(id)] = SkillDefaults{}
	}

	var written []Converted
	waiting := map[string][]string{}
	var refused []string
	ids := slices.Sorted(maps.Keys(examples))
	ids = slices.DeleteFunc(ids, func(id string) bool { return named[id] })
	for _, id := range ids {
		c, err := convertOf(s.standalone, examples[id], topics, batch, BatchTask{}, Overlay{})
		var notYet notDrawn
		switch {
		case errors.As(err, &notYet):
			waiting[notYet.kind] = append(waiting[notYet.kind], id)
		case err != nil:
			refused = append(refused, err.Error())
		default:
			written = append(written, c)
		}
	}
	if s.prove && len(written) > 0 {
		proved := make([]string, len(written))
		for i, c := range written {
			proved[i] = c.ID
		}
		if err := prove(s.goCommand, s.standalone, proved); err != nil {
			return fail(err)
		}
		fmt.Fprintf(stdout, "port: the bench of standalone proves the answers of all %d tasks written\n", len(proved))
	}
	if err := write(s.module, written); err != nil {
		return fail(err)
	}

	report(stdout, s, len(ids), written, waiting, refused)
	if len(refused) > 0 {
		return statusDiffers
	}
	return statusDone
}

// report writes what became of each task of a survey, and then the sums.
func report(out io.Writer, s survey, total int, written []Converted, waiting map[string][]string, refused []string) {
	answers := map[string]int{}
	keyNamed, inTasks, fromPicture := 0, 0, 0
	for _, c := range written {
		if c.FromPicture {
			fromPicture++
		}
		fmt.Fprintf(out, "port: %s: written, %s\n", c.ID, c.Answer)
		for _, note := range c.Notes {
			fmt.Fprintf(out, "port: %s: %s\n", c.ID, note)
		}
		answers[c.Answer]++
		if c.KeyNamed > 0 {
			keyNamed += c.KeyNamed
			inTasks++
		}
	}
	kinds := slices.SortedFunc(maps.Keys(waiting), func(a, b string) int {
		return cmp.Or(cmp.Compare(len(waiting[b]), len(waiting[a])), cmp.Compare(a, b))
	})
	waits := 0
	var byKind []string
	for _, kind := range kinds {
		for _, id := range waiting[kind] {
			fmt.Fprintf(out, "port: %s: waits for a picture of the kind %s\n", id, kind)
		}
		waits += len(waiting[kind])
		byKind = append(byKind, fmt.Sprintf("%s %d", kind, len(waiting[kind])))
	}
	for _, reason := range refused {
		fmt.Fprintf(out, "port: refused: %s\n", reason)
	}

	fmt.Fprintf(out, "port: survey of the %d tasks of standalone at %s that no batch of %s names\n", total, s.commit[:7], s.except)
	fmt.Fprintf(out, "port: written into %s: %d — %s %d, %s %d, %s %d, %s %d, %s %d\n", s.module, len(written),
		answerNumber, answers[answerNumber], answerFraction, answers[answerFraction], answerTime, answers[answerTime],
		answerByHand, answers[answerByHand], answerTapped, answers[answerTapped])
	fmt.Fprintf(out, "port: explanations and leading questions that name the key of a choice or a time: %d, in %d tasks\n", keyNamed, inTasks)
	fmt.Fprintf(out, "port: written tasks solved from the picture: %d\n", fromPicture)
	fmt.Fprintf(out, "port: waiting for pictures: %d", waits)
	if len(byKind) > 0 {
		fmt.Fprintf(out, " — %s", strings.Join(byKind, ", "))
	}
	fmt.Fprintf(out, "\nport: refused: %d\n", len(refused))
}

// batchedTasks are the tasks that the batches of a folder name, every file <name>.yaml of it.
func batchedTasks(dir string) (map[string]bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("the batches: %w", err)
	}
	named := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			return nil, fmt.Errorf("%s: a batch is a file <name>.yaml, and %s is not one", dir, entry.Name())
		}
		batch, err := readBatch(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		for id := range batch.Tasks {
			named[id] = true
		}
	}
	return named, nil
}
