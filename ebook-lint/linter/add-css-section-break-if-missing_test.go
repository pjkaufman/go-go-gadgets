//go:build unit

package linter_test

import (
	"fmt"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	"github.com/stretchr/testify/assert"
)

var contextBreak = "------"

type AddCssSectionBreakIfMissingTestCase struct {
	InputText         string
	InputContextBreak string
	ExpectedOutput    string
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
content: "%s";
display:inline-block;
position:relative;
font-size:1em;
padding:1em;
}`, contextBreak)

var AddCssSectionBreakIfMissingTestCases = map[string]AddCssSectionBreakIfMissingTestCase{
	"make sure that an empty input becomes the hr blank space": {
		InputText:         "",
		InputContextBreak: contextBreak,
		ExpectedOutput:    linter.HrCharacter + "\n" + fmt.Sprintf(linter.HrContentAfterTemplate, contextBreak),
	},
	"make sure that a solely whitespace input becomes the hr blank space": {
		InputText: `
				   `,
		InputContextBreak: contextBreak,
		ExpectedOutput:    linter.HrCharacter + "\n" + fmt.Sprintf(linter.HrContentAfterTemplate, contextBreak),
	},
	"make sure that input that already contains blank space does not get it added": {
		InputText:         cssFileWithHrCharacter,
		InputContextBreak: contextBreak,
		ExpectedOutput:    cssFileWithHrCharacter,
	},
	"make sure that input that does not contain blank space does get it added": {
		InputText: `p {
height: 10px,
}`,
		InputContextBreak: contextBreak,
		ExpectedOutput:    cssFileWithHrCharacter,
	},
}

func TestAddCssSectionBreakIfMissing(t *testing.T) {
	for name, args := range AddCssSectionBreakIfMissingTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.AddCssSectionBreakIfMissing(args.InputText, args.InputContextBreak)

			assert.Equal(t, args.ExpectedOutput, actual)
		})
	}
}
