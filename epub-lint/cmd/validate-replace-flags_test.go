//go:build unit

package cmd_test

import (
	"errors"
	"testing"

	epub "github.com/pjkaufman/go-go-gadgets/epub-lint/cmd"
	"github.com/stretchr/testify/assert"
)

type validateReplaceFlagsTestCase struct {
	inputEpubFile                string
	inputExtraReplaceStringsPath string
	expectedErr                  error
}

var validateReplaceFlagsTestCases = map[string]validateReplaceFlagsTestCase{
	"make sure that an empty epub file paths causes a validation error": {
		inputEpubFile:                "	",
		inputExtraReplaceStringsPath: "file.md",
		expectedErr:                  epub.ErrEpubPathArgEmpty,
	},
	"make sure that a non-epub file for epub file causes a validation error": {
		inputEpubFile:                "file.txt",
		inputExtraReplaceStringsPath: "file.md",
		expectedErr:                  epub.ErrEpubPathArgNonEpub,
	},
	"make sure that an empty extra string replace path causes a validation error": {
		inputEpubFile:                "file.epub",
		inputExtraReplaceStringsPath: "",
		expectedErr:                  epub.ErrExtraStringReplaceArgEmpty,
	},
	"make sure that a non-md extra string replace path causes a validation error": {
		inputEpubFile:                "file.epub",
		inputExtraReplaceStringsPath: "file.txt",
		expectedErr:                  epub.ErrExtraStringReplaceArgNonMd,
	},
	"make sure that an extra string replace path as an md file and a an epub file for epub file passes validation": {
		inputEpubFile:                "file.epub",
		inputExtraReplaceStringsPath: "file.md",
	},
}

func TestValidateReplaceFlags(t *testing.T) {
	for name, args := range validateReplaceFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := epub.ValidateReplaceFlags(args.inputEpubFile, args.inputExtraReplaceStringsPath)

			if err != nil {
				assert.True(t, errors.Is(err, args.expectedErr))
			} else {
				assert.Nil(t, args.expectedErr)
			}
		})
	}
}
