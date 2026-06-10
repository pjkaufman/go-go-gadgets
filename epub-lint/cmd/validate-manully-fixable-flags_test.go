//go:build unit

package cmd_test

import (
	"testing"

	epub "github.com/pjkaufman/go-go-gadgets/epub-lint/cmd"
	"github.com/stretchr/testify/require"
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
	t.Parallel()

	for name, args := range validateManuallyFixableFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := epub.ValidateManuallyFixableFlags(args.inputEpubFile, args.inputRunAll, args.inputRunBrokenLines, args.inputRunSectionBreak, args.inputRunPageBreak, args.inputRunOxfordCommas, args.inputRunAlthoughBut, args.inputRunThoughts, args.inputRunConversation, args.inputRunNecessaryWords, args.inputRunSingleQuotes)

			if err != nil {
				require.ErrorIs(t, err, args.expectedErr)
			} else {
				require.NoError(t, args.expectedErr)
			}
		})
	}
}
