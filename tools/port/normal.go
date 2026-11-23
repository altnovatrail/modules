package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// The normal forms of the format, so that the port judges a choice before module-kit does: the
// options must differ in the normal form of a choice (task.choice-options), and the key is searched
// for, as a phrase, in the explanations of the trap values. The leak check of the format searches
// neither the key of a tapped choice nor a key that the statement shows, so for a choice the
// port's search is the only one.

// choiceForm is a text in the normal form of a choice of the format.
func choiceForm(text string) string { return normalized(text, true) }

// wordForm is a text as a choice with notation word normalizes it: the form of a choice without
// the ".", "!", "?" and spaces at its end, so that "Kim." is "kim" and "?" is blank.
func wordForm(text string) string { return strings.TrimRight(choiceForm(text), ".!? ") }

// timeValue is a value in the notation of a time read. It admits minutes of three digits, so that
// 5:60 and 10:100 are times a trap value may have.
var timeValue = regexp.MustCompile(`^([0-9]{1,2}):([0-9]{2}|[1-9][0-9]{2})$`)

// timeForm is a time as the kind time compares it, its hours and minutes as numbers: 07:30 is
// 7:30, and 5:60 stays 5:60, since the kind never takes the sum of minutes. It is blank for a
// value that is no time read.
func timeForm(text string) string {
	parts := timeValue.FindStringSubmatch(text)
	if parts == nil {
		return ""
	}
	hours, _ := strconv.Atoi(parts[1])
	minutes, _ := strconv.Atoi(parts[2])
	return fmt.Sprintf("%d:%02d", hours, minutes)
}

// twelveHoursOff is a time as a clock face shows it where that is not the time itself, 9:10 for
// 21:10 and 12:15 for 0:15: the slip that the kind time, which counts hours to 24, leaves to a trap
// value. It is blank from 1:00 to 12:59, which a clock face shows as it is.
func twelveHoursOff(key string) string {
	var hours, minutes int
	if _, err := fmt.Sscanf(key, "%d:%d", &hours, &minutes); err != nil || hours > 23 {
		return ""
	}
	switch {
	case hours == 0:
		return fmt.Sprintf("12:%02d", minutes)
	case hours >= 13:
		return fmt.Sprintf("%d:%02d", hours-12, minutes)
	}
	return ""
}

// normalized is a text after NFKC and the format's mapping of dashes, slashes and white space,
// and, where lower says so, the simple lowercase mapping.
func normalized(text string, lower bool) string {
	text = norm.NFKC.String(text)
	text = strings.Map(func(r rune) rune {
		switch {
		case r == '‐', r >= '‒' && r <= '―', r == '−':
			return '-'
		case r == '⁄', r == '∕':
			return '/'
		case lower:
			return unicode.ToLower(r)
		}
		return r
	}, text)
	return strings.Join(strings.FieldsFunc(text, func(r rune) bool { return unicode.Is(unicode.White_Space, r) }), " ")
}

// leakForm is a text in the normal form of the leak check of the format.
func leakForm(text string) string { return normalized(stepOne(text), true) }

// casedForm is the normal form of the leak check without the lowercase mapping: a label A of a
// picture is no word a.
func casedForm(text string) string { return normalized(stepOne(text), false) }

// stepOne is step 1 of the format's normal form: a superscript or subscript digit, and a vulgar
// fraction that a number stands right before or right after, become the replacement character.
func stepOne(text string) string {
	runes := []rune(text)
	out := make([]rune, len(runes))
	for i, r := range runes {
		out[i] = r
		switch {
		case r == '²', r == '³', r == '¹', r >= '⁰' && r <= '₟':
			out[i] = unicode.ReplacementChar
		case vulgar(r) && (i > 0 && unicode.IsNumber(runes[i-1]) || i+1 < len(runes) && unicode.IsNumber(runes[i+1])):
			out[i] = unicode.ReplacementChar
		}
	}
	return string(out)
}

func vulgar(r rune) bool {
	return r >= '¼' && r <= '¾' || r >= '⅐' && r <= '⅟' || r == '↉'
}

// wordChar is a character of a word as the leak check of the format counts one.
func wordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsNumber(r)
}

// names tells whether a text holds a phrase as the leak check finds one: the phrase's normal form
// stands in the text's, and where it begins or ends with a character of a word, no character of a
// word stands right before or right after it. A phrase whose normal form is empty is found nowhere.
func names(text, phrase string) bool { return phraseIn(leakForm(text), leakForm(phrase)) }

// namesCased is names that tells capitals from small letters.
func namesCased(text, phrase string) bool { return phraseIn(casedForm(text), casedForm(phrase)) }

// phraseIn takes a text and a phrase both in a normal form already.
func phraseIn(form, phraseForm string) bool {
	in, what := []rune(form), []rune(phraseForm)
	if len(what) == 0 {
		return false
	}
	for start := 0; start+len(what) <= len(in); start++ {
		if string(in[start:start+len(what)]) != string(what) {
			continue
		}
		end := start + len(what)
		if wordChar(what[0]) && start > 0 && wordChar(in[start-1]) {
			continue
		}
		if wordChar(what[len(what)-1]) && end < len(in) && wordChar(in[end]) {
			continue
		}
		return true
	}
	return false
}
