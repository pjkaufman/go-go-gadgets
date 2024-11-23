//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type setLanguageTestCase struct {
	inputText    string
	inputLang    string
	expectedText string
}

var setLanguageTestCases = map[string]setLanguageTestCase{
	"when the html tag is missing, no change is made": {
		inputText:    "",
		inputLang:    "en",
		expectedText: "",
	},
	"when the html tag is present, but does not have lang, lang should be added": {
		inputText:    `<html xml:lang="en"></html>`,
		inputLang:    "en",
		expectedText: `<html xml:lang="en" lang="en"></html>`,
	},
	"when the html tag is present, but does not have xml:lang, xml:lang should be added": {
		inputText:    `<html lang="en"></html>`,
		inputLang:    "en",
		expectedText: `<html lang="en" xml:lang="en"></html>`,
	},
	"when the html tag is present with an empty lang attribute, the lang should be added": {
		inputText:    `<html lang="" xml:lang="es"></html>`,
		inputLang:    "en",
		expectedText: `<html lang="en" xml:lang="es"></html>`,
	},
	"when the html tag is present with a whitespace lang attribute, the lang should be added": {
		inputText:    `<html lang="   " xml:lang="es"></html>`,
		inputLang:    "en",
		expectedText: `<html lang="en" xml:lang="es"></html>`,
	},
	"when the html tag is present with an empty xml:lang attribute, the xml:lang should be added": {
		inputText:    `<html xml:lang="" lang="es"></html>`,
		inputLang:    "en",
		expectedText: `<html xml:lang="en" lang="es"></html>`,
	},
	"when the html tag is present with a whitespace xml:lang attribute, the xml:lang should be added": {
		inputText:    `<html xml:lang="   " lang="es"></html>`,
		inputLang:    "en",
		expectedText: `<html xml:lang="en" lang="es"></html>`,
	},
	"when the html tag is present with a value for lang, no change should be made even if the lang differs from the provided lang": {
		inputText:    `<html xml:lang="en" lang="es"></html>`,
		inputLang:    "en",
		expectedText: `<html xml:lang="en" lang="es"></html>`,
	},
	"when the html tag is present with a value for xml:lang, no change should be made even if the lang differs from the provided xml:lang": {
		inputText:    `<html xml:lang="es" lang="en"></html>`,
		inputLang:    "en",
		expectedText: `<html xml:lang="es" lang="en"></html>`,
	},
	"make sure we preserve other values in the html element": {
		inputText: `<?xml version='1.0' encoding='utf-8'?>
		<html lang="en" xml:lang="en" xmlns="http://www.w3.org/1999/xhtml"></body></html>`,
		inputLang: "en",
		expectedText: `<?xml version='1.0' encoding='utf-8'?>
		<html lang="en" xml:lang="en" xmlns="http://www.w3.org/1999/xhtml"></body></html>`,
	},
	"make sure we preserve other values in the html element when they are between the lang attributes": {
		inputText: `<?xml version='1.0' encoding='utf-8'?>
		<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xml:lang="en"></body></html>`,
		inputLang: "en",
		expectedText: `<?xml version='1.0' encoding='utf-8'?>
		<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xml:lang="en"></body></html>`,
	},
	"when no language currently exists and the html element does not start the file, make sure that the language attributes are added correctly": {
		inputText: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml">
	<head></head>
	<body>
		<div id="body" xml:lang="en-US">
			<div class="story">
			</div>
		</div>
	</body>
</html>`,
		inputLang: "en",
		expectedText: `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml" lang="en" xml:lang="en">
	<head></head>
	<body>
		<div id="body" xml:lang="en-US">
			<div class="story">
			</div>
		</div>
	</body>
</html>`,
	},
}

func TestSetLanguage(t *testing.T) {
	for name, args := range setLanguageTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.EnsureLanguageIsSet(args.inputText, args.inputLang)
			assert.Equal(t, args.expectedText, actual)
		})
	}
}
