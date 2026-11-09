package main

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// The normal forms of the format, as the port needs them to judge a choice before module-kit does:
// a choice's options must differ in the normal form of a choice (format/README.md, section 15,
// «Choices»; task.choice-options), and the key of a choice is searched for, as a phrase, in the
// explanations of its trap values (section 20). The leak check of the format does not search the
// key of a choice tapped, nor a key that its statement shows, so for a choice the port's search is
// the only one (R05; R20).

// choiceForm is a text as a choice normalizes it: NFKC, the dashes and slashes of step 3 of
// section 15, the white space of step 4, and each character to its simple lowercase mapping.
func choiceForm(text string) string { return normalized(text, true) }

// wordForm is a text as a choice with notation word normalizes it: the form of a choice, which
// then loses each ".", "!", "?" and space at its end, so that "Kim." is "kim" and "?" is blank
// (format/README.md, section 15; amendment T49 to R20).
func wordForm(text string) string { return strings.TrimRight(choiceForm(text), ".!? ") }

// normalized is a text after NFKC, the dashes and slashes of step 3 of section 15 and the white
// space of step 4, and, where lower says so, each character to its simple lowercase mapping.
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

// leakForm is a text in the normal form of the leak check: step 1 of section 15, then as a choice
// normalizes it (section 20).
func leakForm(text string) string { return normalized(stepOne(text), true) }

// casedForm is the normal form of the leak check without the lowercase mapping: a label A of a
// picture is no word a.
func casedForm(text string) string { return normalized(stepOne(text), false) }

// stepOne is step 1 of section 15: a superscript or subscript digit, and a vulgar fraction that a
// number stands right before or right after, become the replacement character.
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

// vulgar tells a vulgar fraction, such as ½, from any other character.
func vulgar(r rune) bool {
	return r >= '¼' && r <= '¾' || r >= '⅐' && r <= '⅟' || r == '↉'
}

// wordChar is a character of a word: a letter, a mark or a number (section 20).
func wordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsMark(r) || unicode.IsNumber(r)
}

// names tells whether a text holds a phrase as the leak check finds one: the phrase's normal form
// stands in the text's, and where it begins or ends with a character of a word, no character of a
// word stands right before or right after it. A phrase whose normal form is empty is found nowhere.
func names(text, phrase string) bool { return phraseIn(leakForm(text), leakForm(phrase)) }

// namesCased is names that tells capitals from small letters.
func namesCased(text, phrase string) bool { return phraseIn(casedForm(text), casedForm(phrase)) }

// phraseIn is the search of a phrase in a text, both in a normal form already.
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
