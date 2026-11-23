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

// Batch is a file of tools/port/batches/: the commit of standalone the tasks come from, the canvas
// of every task, what each skill gives its tasks, and what each task needs that standalone does not
// say.
type Batch struct {
	Standalone string                   `yaml:"standalone"`
	Canvas     Canvas                   `yaml:"canvas"`
	Skills     map[string]SkillDefaults `yaml:"skills"`
	Tasks      map[string]BatchTask     `yaml:"tasks"`
	// overlays is the folder of the overlays as the command was given it, which every task's header
	// names.
	overlays string
	// survey marks the batch a survey makes, never read from a file: its pictures get surveyAlt, and
	// a text that names the key of a choice is counted, not refused.
	survey bool
}

// Canvas is the canvas of a task, as the format has it.
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

// BatchTask is what a task needs beside standalone: the units of its answer — the word its question
// counts in the singular and the plural and, for a measure, its abbreviation — and whether its
// choice is tapped though its key is a word of the question, its wrong options being reasons a
// learner would not write by hand.
type BatchTask struct {
	Units  []string `yaml:"units,omitempty"`
	Tapped bool     `yaml:"tapped,omitempty"`
}

// Overlay is what a person writes for a task, tools/port/overlay/<batch>/<id>.yaml, which the port
// puts into the task as it is.
type Overlay struct {
	Givens  []string `yaml:"givens,omitempty"`
	Unknown string   `yaml:"unknown,omitempty"`
	// LeadingQuestion stands in place of standalone's hint where the hint names the key.
	LeadingQuestion string `yaml:"leading_question,omitempty"`
	SimplerCase     string `yaml:"simpler_case,omitempty"`
	// Techniques, when given, stand in place of those of the task's skill.
	Techniques []string `yaml:"techniques,omitempty"`
	// Traps are added to those of the wrong options that become no trap value.
	Traps       []string `yaml:"traps,omitempty"`
	LeakPhrases []string `yaml:"leak_phrases,omitempty"`
	// PictureAlt and SolutionPictureAlt say what the pictures show, with no answer in them.
	PictureAlt         string `yaml:"picture_alt,omitempty"`
	SolutionPictureAlt string `yaml:"solution_picture_alt,omitempty"`
	// AnswerForm is a paragraph after the question and its picture that says what form the answer
	// takes, such as "Write the time as hours and minutes, like 9:15.", as the format would have a
	// statement say.
	AnswerForm string `yaml:"answer_form,omitempty"`
	// TrapTexts are keyed by the trap value as the module writes it: "5", not "5 m".
	TrapTexts map[string]string `yaml:"trap_texts,omitempty"`
}

var commit = regexp.MustCompile(`^[0-9a-f]{40}$`)

// readBatch refuses a field it does not know.
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

// readOverlays reads the overlays of a batch from its own folder. The folder must be there, or a
// mistyped path would drop every overlay unsaid; an overlay of a task the batch does not name would
// be left out unsaid, so it is an error, and so is a folder in it.
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

// readYAML decodes the one YAML document of a file, refusing fields it does not know.
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
