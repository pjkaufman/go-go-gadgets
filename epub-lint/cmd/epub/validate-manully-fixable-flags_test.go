//go:build unit

package epub_test

import (
	"errors"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/epub-lint/cmd/epub"
	"github.com/stretchr/testify/assert"
)

type validateManuallyFixableFlagsTestCase struct {
	inputEpubFile          string
	inputRunAll            bool
	inputRunBrokenLines    bool
	inputRunSectionBreak   bool
	inputRunPageBreak      bool
	inputRunOxfordCommas   bool
	inputRunAlthoughBut    bool
	inputRunThoughts       bool
	inputRunConversation   bool
	inputRunNecessaryWords bool
	inputRunSingleQuotes   bool
	expectedErr            error
}

var validateManuallyFixableFlagsTestCases = map[string]validateManuallyFixableFlagsTestCase{
	"make sure that all bool flags being false causes a validation error": {
		inputEpubFile: "test.epub",
		expectedErr:   epub.ErrOneRunBoolArgMustBeEnabled,
	},
	"make sure that an empty epub file path causes a validation error": {
		inputEpubFile: "	  ",
		inputRunAll:   true,
		expectedErr:   epub.ErrEpubPathArgEmpty,
	},
	"make sure that a non-epub file path causes a validation error": {
		inputEpubFile: "file.txt",
		inputRunAll:   true,
		expectedErr:   epub.ErrEpubPathArgNonEpub,
	},
	"make sure that an epub file with at least 1 boolean passes validation": {
		inputEpubFile:        "test.epub",
		inputRunSectionBreak: true,
	},
}

func TestValidateManuallyFixableFlags(t *testing.T) {
	for name, args := range validateManuallyFixableFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := epub.ValidateManuallyFixableFlags(args.inputEpubFile, args.inputRunAll, args.inputRunBrokenLines, args.inputRunSectionBreak, args.inputRunPageBreak, args.inputRunOxfordCommas, args.inputRunAlthoughBut, args.inputRunThoughts, args.inputRunConversation, args.inputRunNecessaryWords, args.inputRunSingleQuotes)

			if err != nil {
				assert.True(t, errors.Is(err, args.expectedErr))
			} else {
				assert.Nil(t, args.expectedErr)
			}
		})
	}
}
