//go:build unit

package cmd_test

import (
	"testing"

	epub "github.com/pjkaufman/go-go-gadgets/epub-lint/cmd"
	"github.com/stretchr/testify/require"
)

type validateCompressAndLintFlagsTestCase struct {
	inputLintDir string
	inputLang    string
	expectedErr  error
}

var ValidateCompressAndLintFlagsTestCases = map[string]validateCompressAndLintFlagsTestCase{
	"make sure that an empty lint dir causes a validation error": {
		inputLintDir: "",
		inputLang:    "en",
		expectedErr:  epub.ErrLintDirArgEmpty,
	},
	"make sure that an empty lang causes a validation error": {
		inputLintDir: "package.opf",
		inputLang:    "",
		expectedErr:  epub.ErrLangArgEmpty,
	},
	"make sure that a non-whitespace lint dir and a non-whitespace lang value passes validation": {
		inputLintDir: "folder",
		inputLang:    "en",
	},
}

func TestValidateCompressAndLintFlags(t *testing.T) {
	t.Parallel()

	for name, args := range ValidateCompressAndLintFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := epub.ValidateOptimizeFlags(args.inputLintDir, args.inputLang)

			if err != nil {
				require.ErrorIs(t, err, args.expectedErr)
			} else {
				require.NoError(t, args.expectedErr)
			}
		})
	}
}
