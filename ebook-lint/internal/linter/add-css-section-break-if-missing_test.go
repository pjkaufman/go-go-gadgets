//go:build unit

package linter_test

import (
	"fmt"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

const contextBreak = "------"

type addCssSectionBreakIfMissingTestCase struct {
	inputText         string
	inputContextBreak string
	expectedOutput    string
}

var cssFileWithHrCharacter = fmt.Sprintf(`p {
height: 10px,
}
hr.character {
overflow: visible;
border:0;
text-align:center;
}
hr.character:after {
content: %q;
display:inline-block;
position:relative;
font-size:1em;
padding:1em;
}`, contextBreak)

var addCssSectionBreakIfMissingTestCases = map[string]addCssSectionBreakIfMissingTestCase{
	"make sure that an empty input becomes the hr blank space": {
		inputText:         "",
		inputContextBreak: contextBreak,
		expectedOutput:    linter.HrCharacter + "\n" + fmt.Sprintf(linter.HrContentAfterTemplate, contextBreak),
	},
	"make sure that a solely whitespace input becomes the hr blank space": {
		inputText: `
				   `,
		inputContextBreak: contextBreak,
		expectedOutput:    linter.HrCharacter + "\n" + fmt.Sprintf(linter.HrContentAfterTemplate, contextBreak),
	},
	"make sure that input that already contains blank space does not get it added": {
		inputText:         cssFileWithHrCharacter,
		inputContextBreak: contextBreak,
		expectedOutput:    cssFileWithHrCharacter,
	},
	"make sure that input that does not contain blank space does get it added": {
		inputText: `p {
height: 10px,
}`,
		inputContextBreak: contextBreak,
		expectedOutput:    cssFileWithHrCharacter,
	},
}

func TestAddCssSectionBreakIfMissing(t *testing.T) {
	for name, args := range addCssSectionBreakIfMissingTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.AddCssSectionBreakIfMissing(args.inputText, args.inputContextBreak)

			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
