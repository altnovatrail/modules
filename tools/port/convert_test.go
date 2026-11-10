package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func convertOne(t *testing.T, ex Example, task BatchTask, overlay Overlay) (string, Converted) {
	t.Helper()
	batch := testBatch(ex.ID)
	c, err := convert(ex, batch, task, overlay, "content/examples/counting.gaps.json")
	if err != nil {
		t.Fatalf("convert() error = %v", err)
	}
	return string(c.Files["tasks/"+ex.ID+".yaml"]), c
}

var metres = BatchTask{Units: []string{"metre", "metres", "m"}}

// The whole file of a task with a whole-number answer: the table of docs/architecture.md, section
// 3.5, line by line.
func TestConvertWholeNumber(t *testing.T) {
	ex := example("lamps-1")
	ex.SolutionTotal = "2 × 7 = 14"
	ex.Hint = "How many gaps?"
	overlay := Overlay{
		Givens:      []string{"8 lamps are in a row.", "They are 2 metres apart."},
		Unknown:     "How far the first lamp is from the last.",
		SimplerCase: "Try 3 lamps first.",
		LeakPhrases: []string{"fourteen metres"},
	}
	got, c := convertOne(t, ex, metres, overlay)
	want := `# Written by tools/port from mathtrail-standalone at 89986bc, from
# content/examples/counting.gaps.json, and not to be edited here: what a person adds goes in
# tools/port/overlay/lamps-1.yaml, and the port writes this file again (R12).
id: "lamps-1"
statement:
  - text: "A path has 8 lamps in a row, 2 metres apart. How far is it from the first lamp to the last?"
givens: ["8 lamps are in a row.", "They are 2 metres apart."]
unknown: "How far the first lamp is from the last."
leading_question: "How many gaps?"
canvas:
  kind: "handwriting_sheet"
  grid: "squares"
items:
  - answer_slots:
      - binding: "anywhere"
        check:
          kind: "number"
          key: "14"
          units: ["metre", "metres", "m"]
          trap_values:
            - value: "7"
              trap: "stopped-early"
              text: "That is the number of gaps."
            - value: "8"
              trap: "number-from-text"
              text: "That is the number of lamps."
            - value: "10"
              trap: "wrong-operation"
              text: "Every gap is 2 metres, not one of them."
            - value: "16"
              trap: "off-by-one"
              text: "8 lamps have 7 gaps, not 8."
    simpler_case: "Try 3 lamps first."
solution:
  - text: "8 lamps have 7 gaps of 2 metres: 2 × 7 = 14 metres."
  - text: "2 × 7 = 14"
techniques: ["draw", "small-numbers"]
leak_phrases: ["fourteen metres"]
skills:
  primary: ["counting-gaps"]
band: "grades-1-2"
difficulty: 2
`
	if got != want {
		t.Errorf("the task is\n%s\nwant\n%s", got, want)
	}
	if len(c.Notes) != 0 {
		t.Errorf("notes = %v, want none: the hint is the leading question", c.Notes)
	}
	if again, _ := convertOne(t, ex, metres, overlay); again != got {
		t.Error("a second run writes other bytes")
	}
}

// The hint of standalone is the task's leading question, which the overlay's leading_question
// writes anew; an overlay may not give one to a task with no hint. The key of a choice, which the
// leak check of the format does not search, the port searches there itself, and a survey writes
// the question and says so (amendment T51 to R12).
func TestConvertLeadingQuestion(t *testing.T) {
	ex := example("lamps-5")
	got, _ := convertOne(t, ex, metres, Overlay{})
	if strings.Contains(got, "leading_question") {
		t.Errorf("a task with no hint has a leading question:\n%s", got)
	}
	if _, err := convert(ex, testBatch(ex.ID), metres, Overlay{LeadingQuestion: "How many gaps?"}, "content/examples/x.json"); err == nil ||
		!strings.Contains(err.Error(), "its overlay has leading_question, and the task has no hint") {
		t.Errorf("an overlay's question for a task with no hint: error = %v", err)
	}
	ex.Hint = "How many gaps are there between 8 lamps?"
	if got, _ = convertOne(t, ex, metres, Overlay{LeadingQuestion: "How many gaps lie between the lamps?"}); !strings.Contains(got,
		"leading_question: \"How many gaps lie between the lamps?\"\ncanvas:") {
		t.Errorf("the overlay's question does not stand in place of the hint, before the canvas:\n%s", got)
	}

	for name, how := range map[string]BatchTask{"by hand": {}, "tapped": {Tapped: true}} {
		ex := wordExample("race-9")
		ex.Hint = "Who is slower than Kim, and is anyone slower than Ben?"
		if _, err := convert(ex, testBatch(ex.ID), how, Overlay{}, "content/examples/x.json"); err == nil ||
			!strings.Contains(err.Error(), `its leading question names its key "Ben": rewrite it in the leading_question of its overlay`) {
			t.Errorf("%s: error = %v", name, err)
		}
		if _, err := convert(ex, testBatch(ex.ID), how, Overlay{LeadingQuestion: "Who is slower than Kim?"}, "content/examples/x.json"); err != nil {
			t.Errorf("%s, rewritten: error = %v", name, err)
		}
		batch := testBatch(ex.ID)
		batch.survey = true
		c, err := convert(ex, batch, how, Overlay{}, "content/examples/x.json")
		if err != nil || c.KeyNamed != 1 || !strings.Contains(strings.Join(c.Notes, "\n"), `its leading question names its key "Ben"`) {
			t.Errorf("%s: a survey: %v, %d, %v", name, err, c.KeyNamed, c.Notes)
		}
	}

	// A time's key, which the statement may show, is searched with a leading zero too.
	clock := clockExample("clock-8")
	clock.Hint = "Does it end at 07:30, or later?"
	if _, err := convert(clock, testBatch(clock.ID), BatchTask{}, Overlay{}, "content/examples/x.json"); err == nil ||
		!strings.Contains(err.Error(), `its leading question names its key "7:30"`) {
		t.Errorf("a time: error = %v", err)
	}
}

func TestConvertNumberWithUnit(t *testing.T) {
	ex := example("lamps-2")
	ex.Options = map[string]string{"A": "7 m", "B": "8 m", "C": "14 m", "D": "10 m", "E": "16"}
	got, _ := convertOne(t, ex, metres, Overlay{})
	for _, line := range []string{`key: "14"`, `value: "7"`, `value: "16"`, `units: ["metre", "metres", "m"]`} {
		if !strings.Contains(got, line) {
			t.Errorf("the task lacks %s:\n%s", line, got)
		}
	}
	if strings.Contains(got, `"14 m"`) {
		t.Error("the unit stays in the key")
	}
}

func TestConvertFraction(t *testing.T) {
	ex := example("share-1")
	ex.Topic = "fractions.parts"
	ex.Options = map[string]string{"A": "1/2", "B": "1/3", "C": "3/4", "D": "2", "E": "4/3"}
	got, _ := convertOne(t, ex, BatchTask{}, Overlay{})
	for _, line := range []string{`kind: "fraction"`, `key: "3/4"`, `value: "1/2"`, `value: "4/3"`,
		`traps: ["wrong-operation"]`, `supporting: ["core.fractions-of-a-whole"]`, `primary: ["fractions-parts"]`} {
		if !strings.Contains(got, line) {
			t.Errorf("the task lacks %s:\n%s", line, got)
		}
	}
	if strings.Contains(got, "units") {
		t.Error("a fraction has units")
	}
}

func TestConvertWhichOfThese(t *testing.T) {
	ex := example("odd-one-1")
	ex.Question = "Which of these numbers is the distance from the first lamp to the last?"
	overlay := Overlay{TrapTexts: map[string]string{"8": "Count the gaps, not the lamps."}}
	got, c := convertOne(t, ex, BatchTask{}, overlay)
	for _, line := range []string{`binding: "choice"`, `options: ["7", "8", "14", "10", "16"]`, `kind: "choice"`,
		`key: "14"`, `text: "Count the gaps, not the lamps."`} {
		if !strings.Contains(got, line) {
			t.Errorf("the task lacks %s:\n%s", line, got)
		}
	}
	if strings.Contains(got, "anywhere") || len(c.Notes) != 0 {
		t.Errorf("a question that points at its options is written by hand:\n%s", got)
	}
}

func TestConvertNonNumericWrongOption(t *testing.T) {
	ex := example("lamps-3")
	ex.Options["B"] = "none"
	ex.Options["D"] = "10 km"
	overlay := Overlay{Traps: []string{"number-from-text", "ignored-condition"}, Techniques: []string{"systematic"}}
	got, c := convertOne(t, ex, metres, overlay)
	if strings.Contains(got, `value: "none"`) || strings.Contains(got, `value: "10"`) {
		t.Errorf("a word, or a unit the task does not count in, is a trap value of a number:\n%s", got)
	}
	if len(c.Notes) != 2 || !strings.Contains(c.Notes[0], `its option B, "none", is no value of the check number`) ||
		!strings.Contains(c.Notes[1], "its option D") {
		t.Errorf("notes = %v, want the options left out named", c.Notes)
	}
	for _, line := range []string{`traps: ["number-from-text", "wrong-operation", "ignored-condition"]`, `techniques: ["systematic"]`} {
		if !strings.Contains(got, line) {
			t.Errorf("the task lacks %s:\n%s", line, got)
		}
	}
}

// A check of a number takes a simple fraction too, and a unit with a space before it or none
// (format/README.md, section 15).
func TestConvertNumberNotations(t *testing.T) {
	ex := example("lamps-7")
	ex.Options = map[string]string{"A": "7m", "B": "15/2", "C": "14", "D": "10.5 m", "E": "-16"}
	got, _ := convertOne(t, ex, metres, Overlay{})
	for _, line := range []string{`value: "7"`, `value: "15/2"`, `value: "10.5"`, `value: "-16"`, `kind: "number"`} {
		if !strings.Contains(got, line) {
			t.Errorf("the task lacks %s:\n%s", line, got)
		}
	}
}

func TestConvertTrapTextOverlay(t *testing.T) {
	overlay := Overlay{TrapTexts: map[string]string{"16": "Count the gaps between the lamps once more."}}
	got, _ := convertOne(t, example("lamps-4"), metres, overlay)
	if !strings.Contains(got, `text: "Count the gaps between the lamps once more."`) || strings.Contains(got, "not 8.") {
		t.Errorf("the explanation of 16 is not the overlay's:\n%s", got)
	}
}

func TestConvertRefuses(t *testing.T) {
	for _, test := range []struct {
		name    string
		change  func(*Example, *BatchTask, *Overlay)
		message string
	}{
		{"a unit the task does not count in", func(e *Example, _ *BatchTask, _ *Overlay) { e.Options["C"] = "14 km" },
			`unit "km" is not among`},
		{"an answer too long", func(e *Example, _ *BatchTask, _ *Overlay) { e.Options["C"] = "12345678901234567" },
			"its answer 12345678901234567 is longer than 16"},
		{"an alt of no picture", func(_ *Example, _ *BatchTask, o *Overlay) { o.PictureAlt = "Lamps." },
			"has picture_alt, and the task has no picture"},
		{"an alt of no solution picture", func(_ *Example, _ *BatchTask, o *Overlay) { o.SolutionPictureAlt = "Lamps." },
			"has solution_picture_alt, and its solution has no picture"},
		{"a wrong option equal to the key", func(e *Example, _ *BatchTask, _ *Overlay) { e.Options["D"] = "14.0" },
			"is the value of option C too"},
		{"a value too long", func(e *Example, _ *BatchTask, _ *Overlay) { e.Options["E"] = "12345678901234567" },
			"longer than 16"},
		{"a text rewritten of no value", func(_ *Example, _ *BatchTask, o *Overlay) { o.TrapTexts = map[string]string{"99": "x"} },
			"rewrites the text of 99"},
		{"a skill the batch does not know", func(e *Example, _ *BatchTask, _ *Overlay) { e.Topic = "logic.sets" },
			"says nothing of its skill logic-sets"},
		{"two options alike in a choice", func(e *Example, _ *BatchTask, _ *Overlay) {
			e.Question = "Which of these is it?"
			e.Options["D"] = "7"
		}, "is option A too, as a choice compares them"},
		{"a choice text rewritten of no option", func(e *Example, _ *BatchTask, o *Overlay) {
			e.Question = "Which of these is it?"
			o.TrapTexts = map[string]string{"99": "x"}
		}, "rewrites the text of 99"},
		{"a picture without alt", func(e *Example, _ *BatchTask, _ *Overlay) { e.Picture = json.RawMessage(rowJSON) },
			"has no alt"},
		{"a solution picture without alt", func(e *Example, _ *BatchTask, _ *Overlay) { e.SolutionPicture = json.RawMessage(rowJSON) },
			"has no alt"},
	} {
		t.Run(test.name, func(t *testing.T) {
			ex, task, overlay := example("lamps-5"), metres, Overlay{}
			task.Units = append([]string(nil), metres.Units...)
			test.change(&ex, &task, &overlay)
			_, err := convert(ex, testBatch(ex.ID), task, overlay, "content/examples/x.json")
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Errorf("convert() error = %v, want one that says %q", err, test.message)
			}
		})
	}
}

func TestConvertPictures(t *testing.T) {
	ex := example("lamps-6")
	ex.Picture = json.RawMessage(rowJSON)
	ex.SolutionPicture = json.RawMessage(`{"kind": "row", "items": [{"label": "1"}, {"label": "8"}], "span": "14"}`)
	overlay := Overlay{PictureAlt: "Lamps in a row.", SolutionPictureAlt: "The first and the last lamp."}
	got, c := convertOne(t, ex, metres, overlay)
	for _, line := range []string{`picture: "pictures/lamps-6.png"`, `alt: "Lamps in a row."`,
		`picture: "pictures/lamps-6-solution.png"`, `alt: "The first and the last lamp."`, "width: 272", "height: 74"} {
		if !strings.Contains(got, line) {
			t.Errorf("the task lacks %s:\n%s", line, got)
		}
	}
	for _, file := range []string{"media/pictures/lamps-6.png", "media/pictures/lamps-6-solution.png"} {
		if len(c.Files[file]) == 0 {
			t.Errorf("no file %s", file)
		}
	}
}
