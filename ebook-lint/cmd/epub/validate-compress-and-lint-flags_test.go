//go:build unit

package epub_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub"
	"github.com/stretchr/testify/assert"
)

type ValidateCompressAndLintFlagsTestCase struct {
	InputLintDir  string
	InputLang     string
	ExpectedError string
}

var ValidateCompressAndLintFlagsTestCases = map[string]ValidateCompressAndLintFlagsTestCase{
	"make sure that an empty lint dir causes a validation error": {
		InputLintDir:  "",
		InputLang:     "en",
		ExpectedError: epub.LintDirArgEmpty,
	},
	"make sure that an empty lang causes a validation error": {
		InputLintDir:  "package.opf",
		InputLang:     "",
		ExpectedError: epub.LangArgEmpty,
	},
	"make sure that a non-whitespace lint dir and a non-whitespace lang value passes validation": {
		InputLintDir:  "folder",
		InputLang:     "en",
		ExpectedError: "",
	},
}

func TestValidateCompressAndLintFlags(t *testing.T) {
	for name, args := range ValidateCompressAndLintFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := epub.ValidateCompressAndLintFlags(args.InputLintDir, args.InputLang)

			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
