package epub

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// createTranslatorsNotesCmd represents the createTranslatorsNotes command
var createTranslatorsNotesCmd = &cobra.Command{
	Use:     "create-notes",
	Short:   "Attempts to create translator's notes for fan translated works.",
	Example: heredoc.Doc(`ebook-lint epub create-notes -f test.epub`),
	Long: heredoc.Doc(`Tries to create translator's notes (TN) for fan translated works.
	How it works (not implemented yet):
	- Look for instances of "["
	- For any hits, check with user if the contents of [] should be considered translator's notes
	- If so, check the last 5 lines for "*" as the reference point
	- If not found, then assume it is the end of the last line
	- Lastly check for instances of the word note and repeat the prior process
	- Once all values have been found, make sure they are in order based on the opf
	- Create the TN page
	- Add TN page to opf
	- Replace references
	- Finish everything up
`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateCreateTranslatorsNotesFlags(epubFile)
		if err != nil {
			logger.WriteError(err.Error())
		}

		logger.WriteError("TODO: implement logic described above")
	},
}

func init() {
	EpubCmd.AddCommand(createTranslatorsNotesCmd)

	createTranslatorsNotesCmd.Flags().StringVarP(&epubFile, "epub-file", "f", "", "the epub file to replace strings in in")
	createTranslatorsNotesCmd.MarkFlagRequired("epub-file")
}

func ValidateCreateTranslatorsNotesFlags(epubPath string) error {
	err := validateCommonEpubFlags(epubPath)
	if err != nil {
		return err
	}

	return nil
}
