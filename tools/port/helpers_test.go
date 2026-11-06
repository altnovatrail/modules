package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The tests make up their own standalone: a folder with content/examples/<topic>.json, and tasks
// in its form with texts of their own, so that they read nothing of a clone, and nothing of spec,
// which the shared go.yml does not check out.

const testCommit = "89986bc01a12ff3188cf08194677bfae0464568a"

// example is a task in standalone's form, with a whole-number answer 14 and four wrong options.
func example(id string) Example {
	return Example{
		ID:            id,
		Topic:         "counting.gaps",
		GradeLevel:    "1-2",
		Difficulty:    2,
		Question:      "A path has 8 lamps in a row, 2 metres apart. How far is it from the first lamp to the last?",
		Options:       map[string]string{"A": "7", "B": "8", "C": "14", "D": "10", "E": "16"},
		CorrectAnswer: "C",
		Solution:      "8 lamps have 7 gaps of 2 metres: 2 × 7 = 14 metres.",
		Distractors: map[string]Distractor{
			"A": {Trap: "stopped_early", Text: "That is the number of gaps."},
			"B": {Trap: "number_from_text", Text: "That is the number of lamps."},
			"D": {Trap: "wrong_operation", Text: "Every gap is 2 metres, not one of them."},
			"E": {Trap: "off_by_one", Text: "8 lamps have 7 gaps, not 8."},
		},
	}
}

// git runs git in a folder with none of the caller's settings, and gives what it printed.
func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=Test", "-c", "user.email=test@example.com"}, args...)...)
	command.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	out, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// standaloneWith writes a made-up clone of standalone holding the tasks given, by topic file, and
// commits them, and gives the clone and its commit.
func standaloneWith(t *testing.T, files map[string][]Example) (string, string) {
	t.Helper()
	clone := t.TempDir()
	dir := filepath.Join(clone, "content", "examples")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, tasks := range files {
		data, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(dir, name), string(data))
	}
	git(t, clone, "init", "-q")
	git(t, clone, "add", "-A")
	git(t, clone, "commit", "-q", "--allow-empty", "-m", "Tasks")
	return clone, git(t, clone, "rev-parse", "HEAD")
}

// testBatch is a batch of the tasks given, each counted in metres.
func testBatch(ids ...string) Batch {
	batch := Batch{
		Standalone: testCommit,
		Canvas:     Canvas{Kind: "handwriting_sheet", Grid: "squares"},
		Skills: map[string]SkillDefaults{
			"counting-gaps":   {Techniques: []string{"draw", "small-numbers"}},
			"fractions-parts": {Techniques: []string{"draw"}, Supporting: []string{"core.fractions-of-a-whole"}},
		},
		Tasks: map[string]BatchTask{},
	}
	for _, id := range ids {
		batch.Tasks[id] = BatchTask{Units: []string{"metre", "metres", "m"}}
	}
	return batch
}

// batchFile writes a batch file of the tasks given, of the commit given.
func batchFile(t *testing.T, commit string, ids ...string) string {
	t.Helper()
	text := "standalone: " + commit + "\ncanvas: {kind: handwriting_sheet, grid: squares}\n" +
		"skills:\n  counting-gaps: {techniques: [draw, small-numbers]}\ntasks:\n"
	for _, id := range ids {
		text += "  " + id + ": {units: [metre, metres, m]}\n"
	}
	file := filepath.Join(t.TempDir(), "batch.yaml")
	writeFile(t, file, text)
	return file
}

func writeFile(t *testing.T, file, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// rowJSON is the picture of a row as standalone writes it.
const rowJSON = `{"kind": "row", "items": [{}, {}, {}, {"skip": true}, {}], "gaps": "2"}`
