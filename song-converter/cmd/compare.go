package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	commandhandler "github.com/pjkaufman/go-go-gadgets/pkg/command-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/compare"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/spf13/cobra"
)

var (
	pdfFile, htmlFile string
	numJoinLines      int
	ignoreToCLineNums bool

	compareFlags = flags.Flags{
		Flags: []flags.Flag{
			flags.NewFileFlag(true, false, &htmlFile, "source", "s", "", "the html file that was used to generate the pdf file", []string{"html"}, true),
			flags.NewFileFlag(true, false, &pdfFile, "file", "f", "", "the pdf file to compare with the html file", []string{"pdf"}, true),
			flags.NewIntFlag(false, false, &numJoinLines, "join-lines", "", 0, "the number of lines at the start of the pdf to join together to help make the html and pdf content as similar as possible"),
			flags.NewBoolFlag(false, false, &ignoreToCLineNums, "ignore-page-numbers", "", false, "whether to ignore table of contents page numbers (this is for when the HTML or PDF will not have line numbers in the table of contents, but the other will)"),
		},
	}
)

// compareCmd represents the Compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compares the provided html and pdf file to see if there are any potentially meaningful difference like linebreaks and whitespace differences",
	Example: heredoc.Doc(`To compare a pdf and its html source:
	song-converter compare -s songs.html -f songs.pdf

	To compare a pdf and its html source where the first several lines of text are meant to be the heading on a single line:
	song-converter compare -s songs.html -f songs.pdf --join-lines 4
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return compareFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		pdfText := commandhandler.MustGetCommandOutput("pdftotext", "PDF extraction error", "-layout", pdfFile, "-")
		pdfLines := converter.PdfTextCleanup(pdfText, numJoinLines, ignoreToCLineNums)

		htmlContent, err := filehandler.ReadInFileContents(htmlFile)
		if err != nil {
			logger.WriteFatal(err.Error())
		}

		htmlLines := converter.HtmlToText(htmlContent)

		logger.WriteInfo("-- Alignment of PDF vs HTML lines --")
		diffs := compare.CompareLines(pdfLines, htmlLines)
		for _, diff := range diffs {
			logger.WriteInfo(diff.ToDisplayText())
		}
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)

	err := compareFlags.AddToCmd(compareCmd)
	if err != nil {
		logger.WriteFatal(err.Error())
	}
}
