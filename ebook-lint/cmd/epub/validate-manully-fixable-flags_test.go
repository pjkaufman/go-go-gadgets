//go:build unit

package epub_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/epub"
	"github.com/stretchr/testify/assert"
)

type ValidateManuallyFixableFlagsTestCase struct {
	InputEpubFile        string
	InputRunAll          bool
	InputRunBrokenLines  bool
	InputRunSectionBreak bool
	InputRunPageBreak    bool
	InputRunOxfordCommas bool
	InputRunAlthoughBut  bool
	ExpectedError        string
}

var ValidateManuallyFixableFlagsTestCases = map[string]ValidateManuallyFixableFlagsTestCase{
	"make sure that all bool flags being false causes a validation error": {
		InputEpubFile: "test.epub",
		ExpectedError: epub.OneRunBoolArgMustBeEnabled,
	},
	"make sure that an empty epub file path causes a validation error": {
		InputEpubFile: "	  ",
		InputRunAll:   true,
		ExpectedError: epub.EpubPathArgEmpty,
	},
	"make sure that a non-epub file path causes a validation error": {
		InputEpubFile: "file.txt",
		InputRunAll:   true,
		ExpectedError: epub.EpubPathArgNonEpub,
	},
	"make sure that an epub file with at least 1 boolean passes validation": {
		InputEpubFile:        "test.epub",
		InputRunSectionBreak: true,
		ExpectedError:        "",
	},
}

func TestValidateManuallyFixableFlags(t *testing.T) {
	for name, args := range ValidateManuallyFixableFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := epub.ValidateManuallyFixableFlags(args.InputEpubFile, args.InputRunAll, args.InputRunBrokenLines, args.InputRunSectionBreak, args.InputRunPageBreak, args.InputRunOxfordCommas, args.InputRunAlthoughBut)

			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
