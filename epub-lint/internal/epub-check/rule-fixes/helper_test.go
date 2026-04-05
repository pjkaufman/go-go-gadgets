package rulefixes_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check/positions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkFinalOutputMatches(t *testing.T, input, expectedOutput string, edits ...positions.TextEdit) {
	transformedOutput, err := positions.ApplyEdits("", input, edits)

	require.NoError(t, err)
	assert.Equal(t, expectedOutput, transformedOutput, "The text edits after being applied to the input does not equal the expected output")
}
