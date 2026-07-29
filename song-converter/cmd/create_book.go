package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/spf13/cobra"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
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
		createHtmlFile(stagingDir, coverInputFilePath, coverOutputFile, bodyHtmlOutputFile, "", "font-size: 52pt;", true)
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

// TODO: this needs to take in option for putting the songs in page order or putting them in alphabetical order...
// Note: that alphabetical order will not be perfect given the discrepancy between some of the names in the digital vs. book versions
func buildBookListItems(headerInfo []converter.MdFileInfo) string {
	if len(headerInfo) == 0 {
		return ""
	}

	c := collate.New(language.English)
	sort.Slice(headerInfo, func(i, j int) bool {
		return c.CompareString(headerInfo[i].FileName, headerInfo[j].FileName) < 0
	})

	var (
		pageNumberIndex = make(map[string]int)
		listItems       = strings.Builder{}
		pageNumber      int
		pageInfo        []int
		pageFormat      string
		startingLetter  = "A"
	)
	for _, headerData := range headerInfo {
		pageInfo = headerData.PrimaryPageNumbers
		pageFormat = "%d"
		if len(pageInfo) == 0 {
			pageInfo = headerData.SecondaryPageNumbers
			pageFormat = "(%d)"
		}

		if val, ok := pageNumberIndex[headerData.FileName]; ok {
			pageNumber = pageInfo[val]
		} else {
			pageNumber = pageInfo[0]
			pageNumberIndex[headerData.FileName] = 1
		}

		if startingLetter != headerData.Header[0:1] {
			fmt.Fprintln(&listItems, "<br>")
			startingLetter = headerData.Header[0:1]
		}

		fmt.Fprintf(&listItems, `<li><span class="name">%s</span><span class="page">%s</span></li>`+"\n", headerData.Header, fmt.Sprintf(pageFormat, pageNumber))
	}

	return listItems.String()
}
