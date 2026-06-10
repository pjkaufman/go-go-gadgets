//go:build unit

package potentiallyfixableissue_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type suggesterTestCase struct {
	inputText           string
	expectedSuggestions map[string]string
}

func testSuggesterNoError(t *testing.T, tests map[string]suggesterTestCase, suggester func(string) (map[string]string, error)) {
	t.Helper()
	t.Parallel()

	for name, args := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actual, err := suggester(args.inputText)

			require.NoError(t, err)
			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
