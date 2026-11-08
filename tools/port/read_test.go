package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadExamples(t *testing.T) {
	clone, _ := standaloneWith(t, map[string][]Example{
		"counting.gaps.json":   {example("lamps-1"), example("lamps-2")},
		"fractions.parts.json": {example("share-1")},
	})
	examples, err := readExamples(clone)
	if err != nil {
		t.Fatalf("readExamples() error = %v", err)
	}
	if len(examples) != 3 || examples["lamps-2"].Options["C"] != "14" {
		t.Errorf("readExamples() = %v", examples)
	}
}

func TestReadExamplesRefuses(t *testing.T) {
	if _, err := readExamples(t.TempDir()); err == nil || !strings.Contains(err.Error(), "no content/examples") {
		t.Errorf("an empty folder: error = %v", err)
	}

	twice, _ := standaloneWith(t, map[string][]Example{"a.json": {example("lamps-1")}, "b.json": {example("lamps-1")}})
	if _, err := readExamples(twice); err == nil || !strings.Contains(err.Error(), "there twice") {
		t.Errorf("a task twice: error = %v", err)
	}

	unknown, _ := standaloneWith(t, nil)
	writeFile(t, filepath.Join(unknown, "content", "examples", "a.json"), `[{"id": "x", "source": "a book"}]`)
	if _, err := readExamples(unknown); err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Errorf("a field it does not know: error = %v", err)
	}

	for name, change := range map[string]func(*Example){
		"the ID":             func(e *Example) { e.ID = "Lamps_1" },
		"the question":       func(e *Example) { e.Question = "" },
		"the grade level":    func(e *Example) { e.GradeLevel = "1-7" },
		"the difficulty":     func(e *Example) { e.Difficulty = 6 },
		"an option":          func(e *Example) { delete(e.Options, "E") },
		"an empty option":    func(e *Example) { e.Options["E"] = " " },
		"the right answer":   func(e *Example) { e.CorrectAnswer = "F" },
		"a right distractor": func(e *Example) { e.Distractors["C"] = Distractor{Trap: "x", Text: "y"} },
		"a distractor":       func(e *Example) { delete(e.Distractors, "A") },
		"a sixth":            func(e *Example) { e.Distractors["F"] = Distractor{Trap: "x", Text: "y"} },
	} {
		ex := example("lamps-1")
		change(&ex)
		if err := ex.valid(); err == nil {
			t.Errorf("%s: valid() = nil, want an error", name)
		}
	}
}

func TestReadBatch(t *testing.T) {
	batch, err := readBatch(batchFile(t, testCommit, "lamps-1", "lamps-2"))
	if err != nil {
		t.Fatalf("readBatch() error = %v", err)
	}
	if batch.Standalone != testCommit || len(batch.Tasks) != 2 || batch.Canvas.Grid != "squares" {
		t.Errorf("readBatch() = %+v", batch)
	}
	for text, message := range map[string]string{
		"standalone: 89986bc\ncanvas: {kind: k}\ntasks: {a: {}}\n":              "full SHA",
		"standalone: " + testCommit + "\ntasks: {a: {}}\n":                      "no canvas",
		"standalone: " + testCommit + "\ncanvas: {kind: k}\n":                   "no task",
		"standalone: " + testCommit + "\ncanvas: {kind: k}\ntasks: {A_1: {}}\n": "no task ID",
		"standalone: " + testCommit + "\nsource: a book\n":                      "not found",
		"standalone: " + testCommit + "\n---\nstandalone: x\n":                  "more than one document",
	} {
		file := filepath.Join(t.TempDir(), "batch.yaml")
		writeFile(t, file, text)
		if _, err := readBatch(file); err == nil || !strings.Contains(err.Error(), message) {
			t.Errorf("%q: error = %v, want one that says %q", text, err, message)
		}
	}
	if _, err := readBatch(filepath.Join(t.TempDir(), "none.yaml")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("no file: error = %v", err)
	}
}

func TestReadOverlays(t *testing.T) {
	batch := testBatch("lamps-1", "lamps-2")
	if _, err := readOverlays(filepath.Join(t.TempDir(), "none"), batch); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("no folder: error = %v, want that it is not there", err)
	}
	if overlays, err := readOverlays(t.TempDir(), batch); err != nil || len(overlays) != 0 {
		t.Errorf("an empty folder: %v, %v", overlays, err)
	}

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "lamps-1.yaml"), "unknown: How far.\ntrap_texts: {\"7\": Count again.}\n")
	overlays, err := readOverlays(dir, batch)
	if err != nil || overlays["lamps-1"].Unknown != "How far." || overlays["lamps-1"].TrapTexts["7"] != "Count again." {
		t.Errorf("readOverlays() = %v, %v", overlays, err)
	}

	for name, message := range map[string]string{
		"notes.txt":    "is not one",
		"lamps-9.yaml": "names no task lamps-9",
	} {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, name), "unknown: x\n")
		if _, err := readOverlays(dir, batch); err == nil || !strings.Contains(err.Error(), message) {
			t.Errorf("%s: error = %v, want one that says %q", name, err, message)
		}
	}
	strict := t.TempDir()
	writeFile(t, filepath.Join(strict, "lamps-2.yaml"), "hint: x\n")
	if _, err := readOverlays(strict, batch); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("a field it does not know: error = %v", err)
	}
	notDir := filepath.Join(t.TempDir(), "file")
	writeFile(t, notDir, "")
	if _, err := readOverlays(notDir, batch); err == nil {
		t.Error("a file for a folder: no error")
	}
}

func TestEmitScalars(t *testing.T) {
	for value, want := range map[any]string{
		`He said "no" \ then`: `"He said \"no\" \\ then"`,
		"a < b & c > d":       `"a < b & c > d"`,
		"line\nbreak":         `"line\nbreak"`,
		"2 × 7 ÷ 1 — ½":       `"2 × 7 ÷ 1 — ½"`,
		"no":                  `"no"`,
		7:                     "7",
	} {
		if got := scalar(value); got != want {
			t.Errorf("scalar(%q) = %s, want %s", value, got, want)
		}
	}
	defer func() {
		if recover() == nil {
			t.Error("scalar(1.5) does not panic")
		}
	}()
	scalar(1.5)
}

func TestReadIdea(t *testing.T) {
	clone, _ := standaloneWith(t, nil)
	writeFile(t, filepath.Join(clone, "content", "examples", "a.json"), `[`+exampleJSON(t, example("lamps-1"), `"idea": 14`)+`]`)
	examples, err := readExamples(clone)
	if err != nil || examples["lamps-1"].Idea != 14 {
		t.Errorf("readExamples() = %v, %v", examples, err)
	}
	ex := example("lamps-1")
	ex.Idea = -1
	if err := ex.valid(); err == nil || !strings.Contains(err.Error(), "no place in a list") {
		t.Errorf("a negative idea: error = %v", err)
	}
}

// exampleJSON is a task as standalone writes it, with fields of JSON added.
func exampleJSON(t *testing.T, ex Example, fields string) string {
	t.Helper()
	data, err := json.Marshal(ex)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSuffix(string(data), "}") + ", " + fields + "}"
}

func TestReadTopics(t *testing.T) {
	clone, _ := standaloneWith(t, map[string][]Example{"counting.gaps.json": {example("lamps-1")}})
	topics, err := readTopics(clone)
	if err != nil || len(topics) != 1 || topics["counting.gaps"].Slug != "counting.gaps" {
		t.Errorf("readTopics() = %v, %v", topics, err)
	}
	file := filepath.Join(clone, "content", "catalogs", "topics.json")
	for text, message := range map[string]string{
		`[{"id": "a", "pages": 2}]`:  "unknown field",
		`[{"id": "a"}, {"id": "a"}]`: "is there twice",
		`[{"name": "A"}]`:            "has no ID",
	} {
		writeFile(t, file, text)
		if _, err := readTopics(clone); err == nil || !strings.Contains(err.Error(), message) {
			t.Errorf("%s: error = %v, want one that says %q", text, err, message)
		}
	}
	if _, err := readTopics(t.TempDir()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("no catalog: error = %v", err)
	}
}

func TestOfTopic(t *testing.T) {
	pictured := topicOf("counting.figures")
	pictured.SolvedFromPicture = true
	topics := map[string]Topic{"counting.gaps": topicOf("counting.gaps"), "counting.figures": pictured}
	ex := example("lamps-1")
	if err := ofTopic(ex, topics); err != nil {
		t.Errorf("a task of its topic: %v", err)
	}
	ex.Topic = "logic.sets"
	if err := ofTopic(ex, topics); err == nil || !strings.Contains(err.Error(), "is no topic of content/catalogs/topics.json") {
		t.Errorf("a topic not in the catalog: error = %v", err)
	}
	ex.Topic = "counting.figures"
	if err := ofTopic(ex, topics); err == nil || !strings.Contains(err.Error(), "is solved from the picture, and the task has no picture") {
		t.Errorf("no picture: error = %v", err)
	}
	ex.Picture = json.RawMessage(rowJSON)
	if err := ofTopic(ex, topics); err != nil {
		t.Errorf("a picture: %v", err)
	}
}
