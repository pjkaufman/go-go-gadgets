//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type addCssPageBreakIfMissingTestCase struct {
	input          string
	expectedOutput string
}

var cssFileWithHrBlankSpace = `p {
height: 10px,
}
hr.blankSpace {
border:0;
height:2em;
}`

var addCssPageBreakIfMissingTestCases = map[string]addCssPageBreakIfMissingTestCase{
	"make sure that an empty input becomes the hr blank space": {
		input:          "",
		expectedOutput: linter.HrBlankSpace + "\n",
	},
	"make sure that a solely whitespace input becomes the hr blank space": {
		input: `
				   `,
		expectedOutput: linter.HrBlankSpace + "\n",
	},
	"make sure that input that already contains blank space does not get it added": {
		input:          cssFileWithHrBlankSpace,
		expectedOutput: cssFileWithHrBlankSpace,
	},
	"make sure that input that does not contain blank space does get it added": {
		input: `p {
height: 10px,
}`,
		expectedOutput: cssFileWithHrBlankSpace,
	},
}

func TestAddCssPageBreakIfMissing(t *testing.T) {
	for name, args := range addCssPageBreakIfMissingTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.AddCssPageBreakIfMissing(args.input)

			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
