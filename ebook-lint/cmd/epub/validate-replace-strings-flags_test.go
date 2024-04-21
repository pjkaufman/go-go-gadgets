//go:build unit

package epub_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub"
	"github.com/stretchr/testify/assert"
)

type ValidateReplaceStringsFlagsTestCase struct {
	InputEpubFile                string
	InputExtraReplaceStringsPath string
	ExpectedError                string
}

var ValidateReplaceStringsFlagsTestCases = map[string]ValidateReplaceStringsFlagsTestCase{
	"make sure that an empty epub file paths causes a validation error": {
		InputEpubFile:                "	",
		InputExtraReplaceStringsPath: "file.md",
		ExpectedError:                epub.EpubPathArgEmpty,
	},
	"make sure that a non-epub file for epub file causes a validation error": {
		InputEpubFile:                "file.txt",
		InputExtraReplaceStringsPath: "file.md",
		ExpectedError:                epub.EpubPathArgNonEpub,
	},
	"make sure that an empty extra string replace path causes a validation error": {
		InputEpubFile:                "file.epub",
		InputExtraReplaceStringsPath: "",
		ExpectedError:                epub.ExtraStringReplaceArgEmpty,
	},
	"make sure that a non-md extra string replace path causes a validation error": {
		InputEpubFile:                "file.epub",
		InputExtraReplaceStringsPath: "file.txt",
		ExpectedError:                epub.ExtraStringReplaceArgNonMd,
	},
	"make sure that an extra string replace path as an md file and a an epub file for epub file passes validation": {
		InputEpubFile:                "file.epub",
		InputExtraReplaceStringsPath: "file.md",
		ExpectedError:                "",
	},
}

func TestValidateReplaceStringsFlags(t *testing.T) {
	for name, args := range ValidateReplaceStringsFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := epub.ValidateReplaceStringsFlags(args.InputEpubFile, args.InputExtraReplaceStringsPath)

			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
