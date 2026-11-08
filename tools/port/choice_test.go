package main

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// The answers that are not numbers (R20): a time on a clock and one word that the question holds
// are written by hand, a choice whose options are the check's; any other is a choice tapped.

// clockExample is a task whose answer is a time on a clock.
func clockExample(id string) Example {
	ex := example(id)
	ex.Question = "A lesson starts at 7:00 and lasts half an hour. At what time does it end?"
	ex.Options = map[string]string{"A": "7:30", "B": "7:00", "C": "8:00", "D": "6:30", "E": "7:20"}
	ex.CorrectAnswer = "A"
	ex.Solution = "Half an hour is 30 minutes after 7:00."
	ex.Distractors = map[string]Distractor{
		"B": {Trap: "stopped_early", Text: "That is when it starts."},
		"C": {Trap: "wrong_operation", Text: "Half an hour is not a whole hour."},
		"D": {Trap: "wrong_operation", Text: "It ends later than it starts."},
		"E": {Trap: "off_by_one", Text: "Half an hour has 30 minutes, not 20."},
	}
	return ex
}

// wordExample is a task whose answer is a word: one the question holds, unless it is changed.
func wordExample(id string) Example {
	ex := example(id)
	ex.Question = "Ann runs faster than Kim, and Kim runs faster than Ben. Who is the slowest?"
	ex.Options = map[string]string{"A": "Ann", "B": "Kim", "C": "Ben", "D": "Lea", "E": "Max"}
	ex.CorrectAnswer = "C"
	ex.Solution = "Ann is faster than Kim, who is faster than the last of them."
	ex.Distractors = map[string]Distractor{
		"A": {Trap: "ignored_condition", Text: "That is the fastest one."},
		"B": {Trap: "stopped_early", Text: "Someone is slower still."},
		"D": {Trap: "number_from_text", Text: "Nobody of that name runs."},
		"E": {Trap: "number_from_text", Text: "Nobody of that name runs."},
	}
	return ex
}

func TestConvertTimeByHand(t *testing.T) {
	ex := clockExample("clock-1")
	overlay := Overlay{AnswerForm: "Write the time as on the clock, like 9:15."}
	got, c := convertOne(t, ex, BatchTask{}, overlay)
	want := `statement:
  - text: "A lesson starts at 7:00 and lasts half an hour. At what time does it end?"
  - text: "Write the time as on the clock, like 9:15."
`
	if !strings.Contains(got, want) {
		t.Errorf("the statement does not end with the form of the answer:\n%s", got)
	}
	want = `      - binding: "anywhere"
        check:
          kind: "choice"
          key: "7:30"
          options: ["7:30", "7:00", "8:00", "6:30", "7:20"]
          trap_values:
            - value: "7:00"
              trap: "stopped-early"
              text: "That is when it starts."
`
	if !strings.Contains(got, want) {
		t.Errorf("the time is not a choice written by hand:\n%s", got)
	}
	if c.Answer != answerByHand || strings.Contains(got, `binding: "choice"`) {
		t.Errorf("answer = %q:\n%s", c.Answer, got)
	}
}

func TestConvertWordByHand(t *testing.T) {
	got, c := convertOne(t, wordExample("race-1"), BatchTask{}, Overlay{})
	for _, line := range []string{`binding: "anywhere"`, `kind: "choice"`, `key: "Ben"`, `options: ["Ann", "Kim", "Ben", "Lea", "Max"]`,
		`value: "Max"`} {
		if !strings.Contains(got, line) {
			t.Errorf("the task lacks %s:\n%s", line, got)
		}
	}
	if c.Answer != answerByHand {
		t.Errorf("answer = %q, want %q", c.Answer, answerByHand)
	}
}

func TestConvertTapped(t *testing.T) {
	for name, change := range map[string]func(*Example){
		"a word the question does not hold": func(e *Example) { e.Question = "Three children run. Who is the slowest?" },
		"a word the question holds in another": func(e *Example) {
			e.Question = "Ann runs faster than Kim, and Kim faster than Bens. Who is the slowest?"
		},
		"a phrase":           func(e *Example) { e.Options["C"] = "Ben and Max" },
		"which of these":     func(e *Example) { e.Question = "Which of these children is the slowest?" },
		"a time of a phrase": func(e *Example) { e.Options["C"] = "7:30 on Friday" },
	} {
		t.Run(name, func(t *testing.T) {
			ex := wordExample("race-2")
			change(&ex)
			got, c := convertOne(t, ex, BatchTask{}, Overlay{})
			if c.Answer != answerTapped || !strings.Contains(got, `binding: "choice"`) ||
				!strings.Contains(got, `options: ["Ann", "Kim", "`+ex.Options["C"]+`", "Lea", "Max"]`) {
				t.Errorf("answer = %q:\n%s", c.Answer, got)
			}
		})
	}
}

// A number with words after it, which the batch gives no units for, is no number: a date, a
// quantity of two parts. The report names it, since a batch that forgot its units would make it so.
func TestConvertNumberWithoutUnits(t *testing.T) {
	ex := example("date-1")
	ex.Options = map[string]string{"A": "18 May", "B": "17 May", "C": "19 May", "D": "20 May", "E": "19 June"}
	got, c := convertOne(t, ex, BatchTask{}, Overlay{})
	if c.Answer != answerTapped || !strings.Contains(got, `key: "19 May"`) {
		t.Errorf("answer = %q:\n%s", c.Answer, got)
	}
	if len(c.Notes) != 1 || !strings.Contains(c.Notes[0], `"19 May" is no number in the notation of the checks`) {
		t.Errorf("notes = %v", c.Notes)
	}
}

func TestConvertPercent(t *testing.T) {
	ex := example("pct-1")
	ex.Options = map[string]string{"A": "40%", "B": "4%", "C": "44%", "D": "50 %", "E": "88%"}
	got, c := convertOne(t, ex, BatchTask{Units: []string{"%"}}, Overlay{})
	for _, line := range []string{`kind: "number"`, `key: "44"`, `units: ["%"]`, `value: "40"`, `value: "50"`} {
		if !strings.Contains(got, line) {
			t.Errorf("the task lacks %s:\n%s", line, got)
		}
	}
	if c.Answer != answerNumber {
		t.Errorf("answer = %q", c.Answer)
	}
}

// An explanation that names the key of a choice is rewritten in the overlay, or the port refuses
// the task; a survey writes it, and counts it (R05; R20).
func TestConvertKeyNamed(t *testing.T) {
	for name, ex := range map[string]Example{"by hand": wordExample("race-3"), "tapped": wordExample("race-3")} {
		if name == "tapped" {
			ex.Question = "Which of these children is the slowest?"
		}
		ex.Distractors["B"] = Distractor{Trap: "stopped_early", Text: "BEN is slower still."}
		_, err := convert(ex, testBatch(ex.ID), BatchTask{}, Overlay{}, "content/examples/x.json")
		if err == nil || !strings.Contains(err.Error(), `the explanation of its option B names its key "Ben"`) {
			t.Errorf("%s: error = %v", name, err)
		}

		got, c := convertOne(t, ex, BatchTask{}, Overlay{TrapTexts: map[string]string{"Kim": "Someone is slower still."}})
		if strings.Contains(got, "BEN") || c.KeyNamed != 0 {
			t.Errorf("%s: the rewritten text is not the overlay's:\n%s", name, got)
		}

		batch := testBatch(ex.ID)
		batch.survey = true
		c, err = convert(ex, batch, BatchTask{}, Overlay{}, "content/examples/x.json")
		if err != nil || c.KeyNamed != 1 || !strings.Contains(strings.Join(c.Notes, "\n"), `option B names its key "Ben"`) {
			t.Errorf("%s: a survey: %v, %d, %v", name, err, c.KeyNamed, c.Notes)
		}
	}
	// The key is found as a word: Bens is no Ben.
	ex := wordExample("race-4")
	ex.Distractors["B"] = Distractor{Trap: "stopped_early", Text: "The Bens are slower still."}
	if _, err := convert(ex, testBatch(ex.ID), BatchTask{}, Overlay{}, "content/examples/x.json"); err != nil {
		t.Errorf("a longer word: error = %v", err)
	}
}

func TestConvertChoiceRefuses(t *testing.T) {
	ex := wordExample("race-5")
	ex.Options["E"] = " ann "
	if _, err := convert(ex, testBatch(ex.ID), BatchTask{}, Overlay{}, "content/examples/x.json"); err == nil ||
		!strings.Contains(err.Error(), `its option E, " ann ", is option A too`) {
		t.Errorf("options alike once normalized: error = %v", err)
	}
	ex = wordExample("race-6")
	if _, err := convert(ex, testBatch(ex.ID), BatchTask{}, Overlay{TrapTexts: map[string]string{"Eve": "x"}},
		"content/examples/x.json"); err == nil || !strings.Contains(err.Error(), "rewrites the text of Eve") {
		t.Errorf("a text of no option: error = %v", err)
	}
	ex = example("lamps-8")
	ex.Options["C"] = "14 km"
	if _, err := convert(ex, testBatch(ex.ID), metres, Overlay{}, "content/examples/x.json"); err == nil ||
		!strings.Contains(err.Error(), `unit "km" is not among`) {
		t.Errorf("a unit the batch does not give: error = %v", err)
	}
}

// A survey has no overlays, so every picture gets the same alt; a picture of a kind the port does
// not draw yet is notDrawn, which the survey counts by kind.
func TestConvertSurveyPictures(t *testing.T) {
	ex := example("lamps-9")
	ex.Picture = json.RawMessage(rowJSON)
	ex.SolutionPicture = json.RawMessage(rowJSON)
	batch := testBatch(ex.ID)
	batch.survey = true
	c, err := convert(ex, batch, metres, Overlay{}, "content/examples/x.json")
	if err != nil || strings.Count(string(c.Files["tasks/lamps-9.yaml"]), `alt: "`+surveyAlt+`"`) != 2 {
		t.Errorf("a survey: %v\n%s", err, c.Files["tasks/lamps-9.yaml"])
	}

	ex.Picture = json.RawMessage(`{"kind": "flags", "groups": []}`)
	_, err = convert(ex, batch, metres, Overlay{}, "content/examples/x.json")
	var notYet notDrawn
	if !errors.As(err, &notYet) || notYet.kind != "flags" || !strings.Contains(err.Error(), `kind "flags", which the port does not draw yet`) {
		t.Errorf("flags: error = %v", err)
	}
}

// The review of T18: yes and no stay tapped, a question that points at its options is a choice
// before its answer is read as a number, and an answer like a number that is none is named.
func TestConvertTappedAfterReview(t *testing.T) {
	ex := wordExample("race-7")
	ex.Question = "Is Ben the slowest? Answer yes or no."
	ex.Options = map[string]string{"A": "yes", "B": "maybe", "C": "no", "D": "never", "E": "always"}
	if _, c := convertOne(t, ex, BatchTask{}, Overlay{}); c.Answer != answerTapped {
		t.Errorf("no: answer = %q", c.Answer)
	}

	ex = example("lamps-10")
	ex.Question = "Which of these is the longest?"
	ex.Options["C"] = "2 km"
	if _, c := convertOne(t, ex, metres, Overlay{}); c.Answer != answerTapped {
		t.Errorf("which of these with a unit the batch does not give: answer = %q", c.Answer)
	}

	ex = example("lamps-11")
	ex.Options["C"] = "014"
	if _, c := convertOne(t, ex, BatchTask{}, Overlay{}); c.Answer != answerTapped ||
		len(c.Notes) != 1 || !strings.Contains(c.Notes[0], `"014" is no number in the notation`) {
		t.Errorf("a leading zero: %q, %v", c.Answer, c.Notes)
	}
}

// The texts of the overlay that the learner sees before the task is solved do not name the key of
// a choice, nor answer_form the key of any answer: the leak check of the format does not search a
// key that the statement shows.
func TestConvertOverlayNamesTheKey(t *testing.T) {
	for name, overlay := range map[string]Overlay{
		"givens":       {Givens: []string{"Ben runs slowest."}},
		"unknown":      {Unknown: "Whether Ben is the slowest."},
		"simpler case": {SimplerCase: "Suppose Ben runs alone."},
	} {
		ex := wordExample("race-8")
		if _, err := convert(ex, testBatch(ex.ID), BatchTask{}, overlay, "content/examples/x.json"); err == nil ||
			!strings.Contains(err.Error(), `names its key "Ben"`) {
			t.Errorf("%s: error = %v", name, err)
		}
	}
	for name, ex := range map[string]Example{"a time": clockExample("clock-2"), "a number": example("lamps-12")} {
		key := ex.Options[ex.CorrectAnswer]
		overlay := Overlay{AnswerForm: "Write it like " + key + "."}
		if _, err := convert(ex, testBatch(ex.ID), metres, overlay, "content/examples/x.json"); err == nil ||
			!strings.Contains(err.Error(), "its overlay's answer_form names its key") {
			t.Errorf("%s: error = %v", name, err)
		}
	}
}

// A time is a choice by hand even when the batch gives the task units: it is no number with them.
func TestConvertTimeWithUnits(t *testing.T) {
	if _, c := convertOne(t, clockExample("clock-3"), metres, Overlay{}); c.Answer != answerByHand {
		t.Errorf("answer = %q", c.Answer)
	}
}
