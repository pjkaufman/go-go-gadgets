//go:build unit

package rulefixes_test

import (
	"testing"

	rulefixes "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/rule-fixes"
)

type fixMismatchedLangAttributesTestCase struct {
	input        string
	line, column int
	language     string
	expected     string
}

var fixMismatchedLangAttributesTestCases = map[string]fixMismatchedLangAttributesTestCase{
	"When there are mismatched lang attributes, they are changed to the provided language": {
		input:    `<html lang="en" xml:lang="fr"><body></body></html>`,
		line:     1,
		column:   31,
		language: "en",
		expected: `<html lang="en" xml:lang="en"><body></body></html>`,
	},
	"When there are mismatched lang attributes and none match the provided language, then both are changed to the provided language": {
		input:    `<html lang="en" xml:lang="fr"><body></body></html>`,
		line:     1,
		column:   31,
		language: "de",
		expected: `<html lang="de" xml:lang="de"><body></body></html>`,
	},
	"When the lang attributes already match the provided language,then no change is made": {
		input:    `<html lang="en" xml:lang="en"><body></body></html>`,
		line:     1,
		column:   31,
		language: "en",
		expected: `<html lang="en" xml:lang="en"><body></body></html>`,
	},
	"When there are mismatched languages in single quotes, the single quotes are preserved and the values are updated": {
		input:    `<html lang='en' xml:lang='fr'><body></body></html>`,
		line:     1,
		column:   31,
		language: "de",
		expected: `<html lang='de' xml:lang='de'><body></body></html>`,
	},
	"When there is a mix of single quoted and double quoted lang attributes, the quote usage is preserved when the values are updated": {
		input:    `<html lang='en' xml:lang="fr"><body></body></html>`,
		line:     1,
		column:   31,
		language: "de",
		expected: `<html lang='de' xml:lang="de"><body></body></html>`,
	},
	"When there are a mismatched languages, the whitespace is preserved when updating the values": {
		input:    `<html   lang = 'en'   xml:lang = "fr" ><body></body></html>`,
		line:     1,
		column:   40,
		language: "de",
		expected: `<html   lang = 'de'   xml:lang = "de" ><body></body></html>`,
	},
	"When there are mismatched languages on a non-html element, they will get changed": {
		input:    `<body lang="en" xml:lang="fr"><p>text</p></body>`,
		line:     1,
		column:   31,
		language: "de",
		expected: `<body lang="de" xml:lang="de"><p>text</p></body>`,
	},
	"When the lang lang attribute is not present, then no changes are made": {
		input:    `<html xml:lang="fr"><body></body></html>`,
		line:     1,
		column:   21,
		language: "en",
		expected: `<html xml:lang="fr"><body></body></html>`,
	},
	"When the xml lang attribute is not present, then no changes are made": {
		input:    `<html lang="fr"><body></body></html>`,
		line:     1,
		column:   17,
		language: "en",
		expected: `<html lang="fr"><body></body></html>`,
	},
	"When an empty language is provided, then no changes are made": {
		input:    `<html lang="en" xml:lang="fr"><body></body></html>`,
		line:     1,
		column:   35,
		language: "",
		expected: `<html lang="en" xml:lang="fr"><body></body></html>`,
	},
	"When an invalid line is provided, then no changes are made": {
		input:    `<html lang="en" xml:lang="fr"><body></body></html>`,
		line:     0,
		column:   35,
		language: "de",
		expected: `<html lang="en" xml:lang="fr"><body></body></html>`,
	},
	"When an invalid column is provided, then no changes are made": {
		input:    `<html lang="en" xml:lang="fr"><body></body></html>`,
		line:     1,
		column:   100,
		language: "de",
		expected: `<html lang="en" xml:lang="fr"><body></body></html>`,
	},
}

func TestFixMismatchedLangAttributes(t *testing.T) {
	t.Parallel()

	for name, args := range fixMismatchedLangAttributesTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			edits := rulefixes.FixMismatchedLangAttributes(
				args.line,
				args.column,
				args.input,
				args.language,
			)

			checkFinalOutputMatches(t, args.input, args.expected, edits...)
		})
	}
}
