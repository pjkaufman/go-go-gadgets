package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	"github.com/stretchr/testify/assert"
)

func checkFinalOutputMatches(t *testing.T, input, expectedOutput string, edits ...positions.TextEdit) {
	transformedOutput, err := positions.ApplyEdits("", input, edits)

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, transformedOutput, "The text edits after being applied to the input does not equal the expected output")
}
