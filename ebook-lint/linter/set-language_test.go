//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	"github.com/stretchr/testify/assert"
)

type SetLanguageTestCase struct {
	InputText    string
	InputLang    string
	ExpectedText string
}

var SetLanguageTestCases = map[string]SetLanguageTestCase{
	"when the html tag is missing, no change is made": {
		InputText:    "",
		InputLang:    "en",
		ExpectedText: "",
	},
	"when the html tag is present, but does not have lang, lang should be added": {
		InputText:    `<html xml:lang="en"></html>`,
		InputLang:    "en",
		ExpectedText: `<html xml:lang="en" lang="en"></html>`,
	},
	"when the html tag is present, but does not have xml:lang, xml:lang should be added": {
		InputText:    `<html lang="en"></html>`,
		InputLang:    "en",
		ExpectedText: `<html lang="en" xml:lang="en"></html>`,
	},
	"when the html tag is present with an empty lang attribute, the lang should be added": {
		InputText:    `<html lang="" xml:lang="es"></html>`,
		InputLang:    "en",
		ExpectedText: `<html lang="en" xml:lang="es"></html>`,
	},
	"when the html tag is present with a whitespace lang attribute, the lang should be added": {
		InputText:    `<html lang="   " xml:lang="es"></html>`,
		InputLang:    "en",
		ExpectedText: `<html lang="en" xml:lang="es"></html>`,
	},
	"when the html tag is present with an empty xml:lang attribute, the xml:lang should be added": {
		InputText:    `<html xml:lang="" lang="es"></html>`,
		InputLang:    "en",
		ExpectedText: `<html xml:lang="en" lang="es"></html>`,
	},
	"when the html tag is present with a whitespace xml:lang attribute, the xml:lang should be added": {
		InputText:    `<html xml:lang="   " lang="es"></html>`,
		InputLang:    "en",
		ExpectedText: `<html xml:lang="en" lang="es"></html>`,
	},
	"when the html tag is present with a value for lang, no change should be made even if the lang differs from the provided lang": {
		InputText:    `<html xml:lang="en" lang="es"></html>`,
		InputLang:    "en",
		ExpectedText: `<html xml:lang="en" lang="es"></html>`,
	},
	"when the html tag is present with a value for xml:lang, no change should be made even if the lang differs from the provided xml:lang": {
		InputText:    `<html xml:lang="es" lang="en"></html>`,
		InputLang:    "en",
		ExpectedText: `<html xml:lang="es" lang="en"></html>`,
	},
}

func TestSetLanguage(t *testing.T) {
	for name, args := range SetLanguageTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.EnsureLanguageIsSet(args.InputText, args.InputLang)
			assert.Equal(t, args.ExpectedText, actual)
		})
	}
}
