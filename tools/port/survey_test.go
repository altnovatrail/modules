package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// surveySetUp makes up a standalone with a task of each fate in a survey, a folder of batches that
// names one of them, and an empty module, and gives the arguments of a survey.
func surveySetUp(t *testing.T, refused bool) []string {
	t.Helper()
	keyNamed := wordExample("race-2")
	keyNamed.Question = "Which of these children is the slowest?"
	keyNamed.Distractors["B"] = Distractor{Trap: "stopped_early", Text: "Ben is slower still."}
	clock := example("lamps-4")
	clock.Picture = json.RawMessage(`{"kind": "flags", "groups": []}`)
	tasks := []Example{example("lamps-1"), example("lamps-2"), wordExample("race-1"), keyNamed, clock}
	if refused {
		twice := example("lamps-3")
		twice.Options["D"] = "14.0"
		tasks = append(tasks, twice)
	}
	clone, commit := standaloneWith(t, map[string][]Example{"counting.gaps.json": tasks})
	batches := t.TempDir()
	writeFile(t, filepath.Join(batches, "first.yaml"), "standalone: "+commit+"\ncanvas: {kind: handwriting_sheet}\ntasks: {lamps-1: {}}\n")
	return []string{"--standalone", clone, "--survey", "--commit", commit, "--except", batches, "--module", t.TempDir()}
}

func TestRunSurvey(t *testing.T) {
	args := surveySetUp(t, false)
	status, stdout, stderr := runPort(args...)
	if status != statusDone {
		t.Fatalf("run() = %d, %s%s", status, stdout, stderr)
	}
	for _, line := range []string{
		"port: lamps-2: written, number\n",
		"port: race-1: written, choice by hand\n",
		"port: race-2: written, choice tapped\n",
		`port: race-2: the explanation of its option B names its key "Ben"`,
		"port: lamps-4: waits for a picture of the kind flags\n",
		"port: survey of the 4 tasks of standalone at ",
		": 3 — number 1, fraction 0, choice by hand 1, choice tapped 1\n",
		"port: explanations that name the key of a choice: 1, in 1 tasks\n",
		"port: waiting for pictures: 1 — flags 1\n",
		"port: refused: 0\n",
	} {
		if !strings.Contains(stdout, line) {
			t.Errorf("the report lacks %q:\n%s", line, stdout)
		}
	}
	module := args[len(args)-1]
	for _, file := range []string{"tasks/lamps-2.yaml", "tasks/race-1.yaml", "tasks/race-2.yaml"} {
		if _, err := os.Stat(filepath.Join(module, file)); err != nil {
			t.Errorf("%s is not written: %v", file, err)
		}
	}
	for _, file := range []string{"tasks/lamps-1.yaml", "tasks/lamps-4.yaml"} {
		if _, err := os.Stat(filepath.Join(module, file)); !os.IsNotExist(err) {
			t.Errorf("%s is written: %v", file, err)
		}
	}
}

func TestRunSurveyRefusesATask(t *testing.T) {
	status, stdout, _ := runPort(surveySetUp(t, true)...)
	if status != statusDiffers || !strings.Contains(stdout, "port: refused: lamps-3: its option D, 14.0, is the value of option C too\n") ||
		!strings.Contains(stdout, "port: refused: 1\n") || !strings.Contains(stdout, "waiting for pictures: 1 — flags 1") {
		t.Errorf("run() = %d, %s", status, stdout)
	}
}

func TestRunSurveyProves(t *testing.T) {
	args := surveySetUp(t, false)
	command, called := fakeGo(t, []string{event("pass", "lamps-2"), event("pass", "race-1"), event("pass", "race-2")}, 0)
	status, stdout, stderr := runPort(append(args, "--prove", "--go", command)...)
	if status != statusDone || !strings.Contains(stdout, "proves the answers of all 3 tasks written") {
		t.Errorf("run() = %d, %s%s", status, stdout, stderr)
	}
	if got, _ := os.ReadFile(called); !strings.Contains(string(got), "/^(lamps-2|race-1|race-2)$") {
		t.Errorf("the bench ran on\n%s", got)
	}

	failing, _ := fakeGo(t, []string{event("pass", "lamps-2"), event("fail", "race-1"), event("pass", "race-2")}, 1)
	status, stdout, stderr = runPort(append(args, "--prove", "--go", failing)...)
	if status != statusCannotRun || !strings.Contains(stderr, "race-1: the solver does not prove") || strings.Contains(stdout, "written,") {
		t.Errorf("a failing bench: %d, %s%s", status, stdout, stderr)
	}
}

func TestRunSurveyRefuses(t *testing.T) {
	args := surveySetUp(t, false)
	with := func(flag, value string) []string {
		out := append([]string(nil), args...)
		for i := range out {
			if out[i] == flag {
				out[i+1] = value
			}
		}
		return out
	}
	notFolder := filepath.Join(t.TempDir(), "file")
	writeFile(t, notFolder, "")
	strayFile := t.TempDir()
	writeFile(t, filepath.Join(strayFile, "notes.txt"), "")
	badBatch := t.TempDir()
	writeFile(t, filepath.Join(badBatch, "first.yaml"), "standalone: x\n")
	for name, test := range map[string]struct {
		args    []string
		message string
	}{
		"with a batch":            {append(args, "--batch", "b.yaml"), "usage"},
		"with an overlay":         {append(args, "--overlay", "o"), "usage"},
		"with --check":            {append(args, "--check"), "usage"},
		"without --commit":        {with("--commit", ""), "usage"},
		"without --except":        {with("--except", ""), "usage"},
		"--commit and a batch":    {[]string{"--standalone", args[1], "--commit", args[4], "--batch", "b", "--overlay", "o", "--module", "m"}, "usage"},
		"a short commit":          {with("--commit", "89986bc"), "is not the full SHA"},
		"another commit":          {with("--commit", testCommit), "and the batch is of " + testCommit},
		"no batches":              {with("--except", filepath.Join(t.TempDir(), "none")), "the batches"},
		"a file of no batch":      {with("--except", strayFile), "notes.txt is not one"},
		"a bad batch":             {with("--except", badBatch), "full SHA"},
		"a module that is a file": {with("--module", notFolder), "is no folder"},
	} {
		status, _, stderr := runPort(test.args...)
		if status != statusCannotRun || !strings.Contains(stderr, test.message) {
			t.Errorf("%s: %d, %s", name, status, stderr)
		}
	}
}

// A task of a topic the catalog lacks is refused by a run as by a survey.
func TestRunTopicNotInTheCatalog(t *testing.T) {
	args, _ := setUp(t)
	writeFile(t, filepath.Join(args[1], "content", "catalogs", "topics.json"), "[]")
	git(t, args[1], "commit", "-q", "-am", "No topics")
	batch := batchFile(t, git(t, args[1], "rev-parse", "HEAD"), "lamps-1")
	status, _, stderr := runPort("--standalone", args[1], "--batch", batch, "--overlay", t.TempDir(), "--module", t.TempDir())
	if status != statusCannotRun || !strings.Contains(stderr, "lamps-1: its topic counting.gaps is no topic") {
		t.Errorf("run() = %d, %s", status, stderr)
	}
	writeFile(t, filepath.Join(args[1], "content", "catalogs", "topics.json"), "{")
	git(t, args[1], "commit", "-q", "-am", "A broken catalog")
	batch = batchFile(t, git(t, args[1], "rev-parse", "HEAD"), "lamps-1")
	if status, _, stderr := runPort("--standalone", args[1], "--batch", batch, "--overlay", t.TempDir(), "--module", t.TempDir()); status != statusCannotRun ||
		!strings.Contains(stderr, "topics.json") {
		t.Errorf("a broken catalog: %d, %s", status, stderr)
	}
}
