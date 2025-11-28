package cmd

import (
	"github.com/MakeNowJust/heredoc"
	epubcheck "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check"
	commandhandler "github.com/pjkaufman/go-go-gadgets/pkg/command-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	outputToFile string
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate an EPUB file using EPUBCheck",
	Long: heredoc.Doc(`Validates an EPUB file using W3C EPUBCheck tool.
	If EPUBCheck is not installed, it will automatically download and install the latest version.`),
	Example: heredoc.Doc(`
	epub-lint validate -f test.epub
	will run EPUBCheck against the file specified.
`),
	Run: func(cmd *cobra.Command, args []string) {
		err := validateCommonEpubFlags(epubFile)
		if err != nil {
			logger.WriteError(err.Error())
		}

		epubcheckDir, err := filehandler.GetDataDir("epubcheck")
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = epubcheck.EnsureEPUBCheckIsInstalled(epubcheckDir)
		if err != nil {
			logger.WriteError(err.Error())
		}

		jarPath := filehandler.JoinPath(epubcheckDir, "epubcheck.jar")
		extraInputs := []string{"-jar", jarPath, epubFile}

		output := commandhandler.MustGetCommandOutputEvenIfExitError("java", "failed to run EPUBCheck", extraInputs...)

		if outputToFile != "" {
			err = filehandler.WriteFileContents(outputToFile, output)

			if err != nil {
				logger.WriteError(err.Error())
			}
		} else {
			logger.WriteInfo(output)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&epubFile, "file", "f", "", "the epub file to validate")
	err := validateCmd.MarkFlagRequired("file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as required on validate command: %v\n", err)
	}

	err = validateCmd.MarkFlagFilename("file", "epub")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"file\" as looking for specific file types on validate command: %v\n", err)
	}

	validateCmd.Flags().StringVarP(&outputToFile, "out", "", "", "specifies that the validation output should be in the specified file")
}
