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

	// TODO: should some scenarios be removed from the following?
	for i, headerData := range headerInfo {
		if headerData.AlternateTitle == "" {
			continue
		}

		var (
			primaryPageCount   = len(headerData.PrimaryPageNumbers)
			secondaryPageCount = len(headerData.SecondaryPageNumbers)
			hasPrimaryPages    = primaryPageCount > 0
			hasSecondaryPages  = secondaryPageCount > 0
		)
		if hasPrimaryPages && !hasSecondaryPages {
			switch primaryPageCount {
			case 2: // likely wrong, but not sure I can do much about it right now...
				headerData.Header = headerData.AlternateTitle
				headerData.AlternateTitle = ""
				headerInfo = append(headerInfo, headerData)
				continue
			case 1:
				headerData.PrimaryPageNumbers = append(headerData.PrimaryPageNumbers, headerData.PrimaryPageNumbers...)
				headerData.Header = headerData.AlternateTitle
				headerData.AlternateTitle = ""
				headerInfo = append(headerInfo, headerData)
				continue
			}
		} else if hasSecondaryPages && !hasPrimaryPages {
			switch secondaryPageCount {
			case 2:
				headerData.Header = headerData.AlternateTitle
				headerData.AlternateTitle = ""
				headerInfo = append(headerInfo, headerData)
				continue
			case 1:
				headerData.SecondaryPageNumbers = append(headerData.SecondaryPageNumbers, headerData.SecondaryPageNumbers...)
				headerData.Header = headerData.AlternateTitle
				headerData.AlternateTitle = ""
				headerInfo = append(headerInfo, headerData)
				continue
			}
		} else if primaryPageCount == 1 && secondaryPageCount == 1 {
			var secondaryPages = headerData.SecondaryPageNumbers

			headerData.SecondaryPageNumbers = []int{}
			headerInfo[i] = headerData

			headerData.PrimaryPageNumbers = headerData.SecondaryPageNumbers
			headerData.SecondaryPageNumbers = secondaryPages
			headerData.Header = headerData.AlternateTitle
			headerData.AlternateTitle = ""
			headerInfo = append(headerInfo, headerData)

			continue
		}

		logger.WriteFatalf("Encountered an unhandled situation for creating book ToC entries: title %q; alternate title %q; primary page count %d; secondary page count %d\n", headerData.Header, headerData.AlternateTitle, primaryPageCount, secondaryPageCount)
	}

	c := collate.New(language.English)
	sort.Slice(headerInfo, func(i, j int) bool {
		if headerInfo[i].Header != headerInfo[j].Header {
			return c.CompareString(headerInfo[i].Header, headerInfo[j].Header) < 0
		}

		return c.CompareString(headerInfo[i].FileName, headerInfo[j].FileName) < 0
	})

	var (
		primaryPageNumberIndex   = make(map[string]int)
		secondaryPageNumberIndex = make(map[string]int)
		listItems                = strings.Builder{}
		startingLetter           = "A"
	)
	for _, headerData := range headerInfo {
		if startingLetter != headerData.Header[0:1] {
			fmt.Fprintln(&listItems, "<br>")
			startingLetter = headerData.Header[0:1]
		}

		if _, ok := primaryPageNumberIndex[headerData.FileName]; !ok {
			for range headerData.SecondaryPageNumbers {
				addToCEntry(&listItems, headerData.FileName, headerData.Header, "(%d)", headerData.SecondaryPageNumbers, secondaryPageNumberIndex)
			}
		}

		addToCEntry(&listItems, headerData.FileName, headerData.Header, "%d", headerData.PrimaryPageNumbers, primaryPageNumberIndex)
	}

	return listItems.String()
}

func addToCEntry(tocItems *strings.Builder, fileName, header, pageFormat string, pageNumbers []int, pageNumberIndex map[string]int) {
	if len(pageNumbers) == 0 {
		return
	}

	var pageNumber int
	if val, ok := pageNumberIndex[fileName]; ok {
		pageNumber = pageNumbers[val]
	} else {
		pageNumber = pageNumbers[0]
		pageNumberIndex[fileName] = 1
	}

	fmt.Fprintf(tocItems, `<li><span class="name">%s</span><span class="page">%s</span></li>`+"\n", header, fmt.Sprintf(pageFormat, pageNumber))
}
