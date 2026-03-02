package cmd

import (
	"errors"
	"strings"

	"github.com/MakeNowJust/heredoc"
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
	// wsCollapse        = regexp.MustCompile(`\s+`)
	// tocCollapse       = regexp.MustCompile(`(.+?)  +(\d+)$`) // finds toc page numbers
)

const (
	PdfPathArgEmpty     = "file must have a non-whitespace value"
	HtmlPathArgEmpty    = "source must have a non-whitespace value"
	PdfPathNotPdfFile   = "file must be a pdf file"
	HtmlPathNotHtmlFile = "source must be an html file"
)

// CompareCmd represents the Compare command
var CompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compares the provided html and pdf file to see if there are any potentially meaningful difference like linebreaks and whitespace differences",
	Example: heredoc.Doc(`To compare a pdf and its html source:
	song-converter compare -s songs.html -f songs.pdf

	To compare a pdf and its html source where the first several lines of text are meant to be the heading on a single line:
	song-converter compare -s songs.html -f songs.pdf --join-lines 4
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateCompareHtmlFlags(htmlFile, pdfFile)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(htmlFile, "source")
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileArgExists(pdfFile, "file")
		if err != nil {
			logger.WriteError(err.Error())
		}

		pdfText := commandhandler.MustGetCommandOutput("pdftotext", "PDF extraction error", "-layout", pdfFile, "-")
		pdfLines := converter.PdfTextCleanup(pdfText, numJoinLines, ignoreToCLineNums)

		htmlContent, err := filehandler.ReadInFileContents(htmlFile)
		if err != nil {
			logger.WriteError(err.Error())
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
	rootCmd.AddCommand(CompareCmd)

	CompareCmd.Flags().StringVarP(&htmlFile, "source", "s", "", "the html file that was used to generate the pdf file")
	err := CompareCmd.MarkFlagRequired("source")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"source\" as required on compare command: %v\n", err)
	}

	err = CompareCmd.MarkFlagFilename("source", "html")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"source\" as a looking for specific file types on compare command: %v\n", err)
	}

	CompareCmd.Flags().StringVarP(&pdfFile, "file", "f", "", "the pdf file to compare with the html file")
	err = CompareCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as required on compare command: %v\n", err)
	}

	err = CompareCmd.MarkFlagFilename("file", "pdf")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as a looking for specific file types on compare command: %v\n", err)
	}

	CompareCmd.Flags().IntVarP(&numJoinLines, "join-lines", "", 0, "the number of lines at the start of the pdf to join together to help make the html and pdf content as similar as possible")
	CompareCmd.Flags().BoolVarP(&ignoreToCLineNums, "ignore-page-numbers", "", false, "whether to ignore table of contents page numbers (this is for when the HTML or PDF will not have line numbers in the table of contents, but the other will)")
}

func ValidateCompareHtmlFlags(htmlFilePath, pdfFilePath string) error {
	if strings.TrimSpace(htmlFilePath) == "" {
		return errors.New(HtmlPathArgEmpty)
	}

	if !strings.HasSuffix(htmlFilePath, ".html") {
		return errors.New(HtmlPathNotHtmlFile)
	}

	if strings.TrimSpace(pdfFilePath) == "" {
		return errors.New(PdfPathArgEmpty)
	}

	if !strings.HasSuffix(pdfFilePath, ".pdf") {
		return errors.New(PdfPathNotPdfFile)
	}

	return nil
}
