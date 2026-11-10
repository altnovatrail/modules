package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"go.yaml.in/yaml/v4"
)

// Batch is a batch of tasks to carry, a file of tools/port/batches/: the commit of standalone the
// tasks come from, the canvas of every task, what each skill gives its tasks, and the tasks by ID
// with what each needs that standalone does not say (R07; docs/architecture.md, section 3.6).
type Batch struct {
	Standalone string                   `yaml:"standalone"`
	Canvas     Canvas                   `yaml:"canvas"`
	Skills     map[string]SkillDefaults `yaml:"skills"`
	Tasks      map[string]BatchTask     `yaml:"tasks"`
	// overlays is the folder of the batch's overlays, as the command was given it, which the header of
	// every task it writes names.
	overlays string
	// survey marks the batch that a survey makes, of every task no batch names: it is never read
	// from a file, its pictures get surveyAlt, and an explanation that names the key of a choice is
	// counted, not refused.
	survey bool
}

// Canvas is the canvas of a task (format/README.md, section 14).
type Canvas struct {
	Kind string `yaml:"kind"`
	Grid string `yaml:"grid,omitempty"`
}

// SkillDefaults is what a skill gives each task of it: the techniques its site page links to the
// topic, and the anchor skills that support it.
type SkillDefaults struct {
	Techniques []string `yaml:"techniques"`
	Supporting []string `yaml:"supporting,omitempty"`
}

// BatchTask is what a task of a batch needs beside standalone: the units of its answer, the word
// its question counts in the singular and the plural and, for a measure, its abbreviation
// (format/README.md, section 15); and whether its choice is tapped though its key is a word of the
// question, because its wrong options are reasons a learner would not write by hand (R20).
type BatchTask struct {
	Units  []string `yaml:"units,omitempty"`
	Tapped bool     `yaml:"tapped,omitempty"`
}

// Overlay is what a person writes for a task, tools/port/overlay/<batch>/<id>.yaml: none of it comes from
// standalone, and the port puts it into the task as it is (R12).
type Overlay struct {
	Givens      []string `yaml:"givens,omitempty"`
	Unknown     string   `yaml:"unknown,omitempty"`
	SimplerCase string   `yaml:"simpler_case,omitempty"`
	// Techniques, when given, stand in place of those of the task's skill.
	Techniques []string `yaml:"techniques,omitempty"`
	// Traps are added to those of the wrong options that become no trap value.
	Traps       []string `yaml:"traps,omitempty"`
	LeakPhrases []string `yaml:"leak_phrases,omitempty"`
	// PictureAlt and SolutionPictureAlt are the alt of the pictures of the statement and of the
	// solution: what the picture shows, in words, with no answer in them.
	PictureAlt         string `yaml:"picture_alt,omitempty"`
	SolutionPictureAlt string `yaml:"solution_picture_alt,omitempty"`
	// AnswerForm is a paragraph of the statement after the question and its picture, which says what
	// form the answer takes, such as "Write the time as hours and minutes, like 9:15.", as a statement
	// SHOULD say (format/README.md, section 15).
	AnswerForm string `yaml:"answer_form,omitempty"`
	// TrapTexts are explanations written anew, by the trap value they explain, as the module
	// writes it: "5", not "5 m".
	TrapTexts map[string]string `yaml:"trap_texts,omitempty"`
}

var commit = regexp.MustCompile(`^[0-9a-f]{40}$`)

// readBatch reads a batch file strictly: a field it does not know is an error.
func readBatch(file string) (Batch, error) {
	var batch Batch
	if err := readYAML(file, &batch); err != nil {
		return Batch{}, err
	}
	var problems []string
	if !commit.MatchString(batch.Standalone) {
		problems = append(problems, "standalone is not the full SHA of a commit")
	}
	if batch.Canvas.Kind == "" {
		problems = append(problems, "it names no canvas")
	}
	if len(batch.Tasks) == 0 {
		problems = append(problems, "it names no task")
	}
	for id := range batch.Tasks {
		if !taskID.MatchString(id) {
			problems = append(problems, fmt.Sprintf("%q is no task ID", id))
		}
	}
	if len(problems) > 0 {
		slices.Sort(problems)
		return Batch{}, fmt.Errorf("%s: %s", file, strings.Join(problems, "; "))
	}
	return batch, nil
}

// readOverlays reads the overlay of every task of a batch that has one, from the batch's own
// folder, tools/port/overlay/<batch>/ (T22). The folder must be there: a mistyped path would drop
// every overlay unsaid. An overlay of a task the batch does not name is an error too, since it would
// be left out, and so is a folder in it.
func readOverlays(dir string, batch Batch) (map[string]Overlay, error) {
	overlays := map[string]Overlay{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("the overlays: %w", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		id, isYAML := strings.CutSuffix(name, ".yaml")
		if entry.IsDir() || !isYAML {
			return nil, fmt.Errorf("%s: an overlay is a file <id>.yaml, and %s is not one", dir, name)
		}
		if _, named := batch.Tasks[id]; !named {
			return nil, fmt.Errorf("%s: the batch names no task %s", filepath.Join(dir, name), id)
		}
		var overlay Overlay
		if err := readYAML(filepath.Join(dir, name), &overlay); err != nil {
			return nil, err
		}
		overlays[id] = overlay
	}
	return overlays, nil
}

// readYAML decodes one YAML document of a file into value, refusing fields it does not know.
func readYAML(file string, value any) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}
	var more any
	if err := decoder.Decode(&more); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%s: it holds more than one document", file)
	}
	return nil
}
