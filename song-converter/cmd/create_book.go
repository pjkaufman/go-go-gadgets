package cmd

import (
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
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

// CreateBookCmd represents the CreateSongs command
var CreateBookCmd = &cobra.Command{
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
	PreRunE: validateCreateHtmlFile,
	Run: func(cmd *cobra.Command, args []string) {
		createHtmlFile(stagingDir, coverInputFilePath, coverOutputFile, bodyHtmlOutputFile, "", "font-size: 52pt;", true)
	},
}

func init() {
	createCmd.AddCommand(CreateBookCmd)

	createCommonHtmlAndBookFlags(CreateBookCmd)

	CreateBookCmd.Flags().StringVarP(&location, "location", "l", "", "the specific book to recreate by filtering songs down to just that book location")
	err := CreateBookCmd.MarkFlagRequired("location")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"location\" as required on create book command: %v\n", err)
	}
	// TODO: add a sort type for alphabetical and a sort type for in order for the TOC order...
}

// TODO: this needs to take in option for putting the songs in page order or putting them in alphabetical order...
// Note: that alphabetical order will not be perfect given the discrepancy between some of the names in the digital vs. book versions
func buildBookListItems(headerInfo []converter.MdFileInfo) string {
	if len(headerInfo) == 0 {
		return ""
	}

	var (
		pageNumberIndex = make(map[string]int)
		listItems       = strings.Builder{}
		pageNumber      int
	)
	for _, headerData := range headerInfo {
		if val, ok := pageNumberIndex[headerData.FileName]; ok {
			pageNumber = headerData.PageNumbers[val]
		} else {
			pageNumber = headerData.PageNumbers[0]
			pageNumberIndex[headerData.FileName] = 1
		}

		fmt.Fprintf(&listItems, `<li><span class="name">%s</span><span class="page">%d</span></li>`+"\n", headerData.Header, pageNumber)
	}

	return listItems.String()
}
