package cmd

import (
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	tableofcontents "github.com/pjkaufman/go-go-gadgets/song-converter/internal/table-of-contents"
	"github.com/spf13/cobra"
)

const (
	bookFormat = `<html>
  <body>
    <section id="cover">
      %s
    </section>
    <section id="contents">
    <h1 class="toc">Index</h1>
    <ul>
        %s
    </ul>
    </section>
    <section id="songs">
      %s
    </section>
  </body>
</html>`
)

var createBookFlags = flags.Flags{
	Flags: []flags.Flag{
		flags.NewStringFlag(true, false, &location, "location", "l", "", "the specific book to recreate by filtering songs down to just that book location"),
		// TODO: add a sort type for alphabetical and a sort type for in order for the TOC order...
	},
}

// createBookCmd represents the CreateSongs command
var createBookCmd = &cobra.Command{
	Use: "book",
	// Short: "Converts the cover and all Markdown files in the specified folder into html in alphabetical order generating three sections: the cover, table of contents, and songs",
	// Example: heredoc.Doc(`To write the output of converting the files in the specified folder to html to a file:
	// song-converter create-html -d working-dir -c cover.md -o songs.html

	// To write the output of converting the files in the specified folder to html to std out:
	// song-converter create-html -d working-dir -s cover.md
	// `),
	// Long: heredoc.Doc(`How it works:
	// - Reads in all of the files in the specified folder
	// - Sorts the files alphabetically
	// - Adds the cover to the start of the content after converting it to html
	// - Converts each file into html
	// - Writes the content to the specified source
	// `),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := validateCreateHtmlFile()
		if err != nil {
			return err
		}

		return createBookFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		createHtmlFile(stagingDir, coverInputFilePath, coverOutputFile, bodyHtmlOutputFile, "", "left", "", true, tableofcontents.BuildMultipleFileAlphabeticalListItems)
	},
}

func init() {
	createCmd.AddCommand(createBookCmd)

	err := commonBookFlags.AddToCmd(createBookCmd)
	if err != nil {
		logger.WriteFatal(err.Error())
	}

	err = createBookFlags.AddToCmd(createBookCmd)
	if err != nil {
		logger.WriteFatal(err.Error())
	}
}
