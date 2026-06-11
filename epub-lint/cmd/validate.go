package cmd

import (
	"github.com/MakeNowJust/heredoc"
	epubcheck "github.com/pjkaufman/go-go-gadgets/epub-lint/internal/epub-check"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	commandhandler "github.com/pjkaufman/go-go-gadgets/pkg/command-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	outputToFile string
	validateFlag = flags.Flags{
		Flags: []flags.Flag{
			flags.NewFileFlag(true, false, &epubFile, "file", "f", "", "the epub file to validate", []string{"epub"}, true),
			flags.NewFileFlag(false, false, &outputToFile, "out", "", "", "specifies that the validation output should be in the specified file", nil, false),
		},
	}
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFlag.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
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

	err := validateFlag.AddToCmd(validateCmd)
	if err != nil {
		logger.WriteError(err.Error())
	}
}
