package main

import "testing"

func TestChoiceForm(t *testing.T) {
	for text, want := range map[string]string{
		"  ODD ":    "odd",
		"Ann \tKim": "ann kim",
		"＜":         "<",
		"7–7":       "7-7",
		"1⁄2":       "1/2",
		"Kim.":      "kim.",
		"x−y∕z":     "x-y/z",
	} {
		if got := choiceForm(text); got != want {
			t.Errorf("choiceForm(%q) = %q, want %q", text, got, want)
		}
	}
}

// The search of a phrase, by the examples of format/README.md, section 20.
func TestNames(t *testing.T) {
	for _, test := range []struct {
		text, phrase string
		want         bool
	}{
		{"Odd.", "odd", true},
		{"even/odd", "odd", true},
		{"odds", "odd", false},
		{"56<65", "<", true},
		{"It ends at 7:30.", "7:30", true},
		{"It ends at 17:30.", "7:30", false},
		{"Ask BEN first", "Ben", true},
		{"Benny", "Ben", false},
		{"x", " ", false},
		{"7²", "72", false},
		{"1½ apples", "11", false},
		{"½ an apple", "1/2", true},
		{"Ann and  Ben", "ann and ben", true},
	} {
		if got := names(test.text, test.phrase); got != test.want {
			t.Errorf("names(%q, %q) = %v, want %v", test.text, test.phrase, got, test.want)
		}
	}
}
