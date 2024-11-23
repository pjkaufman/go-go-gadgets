//go:build unit

package epub_test

import (
	"errors"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub"
	"github.com/stretchr/testify/assert"
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
	for name, args := range ValidateCompressAndLintFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := epub.ValidateCompressAndLintFlags(args.inputLintDir, args.inputLang)

			if err != nil {
				assert.True(t, errors.Is(err, args.expectedErr))
			} else {
				assert.Nil(t, args.expectedErr)
			}
		})
	}
}
