package cmd

import (
	"os"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "song-converter",
	Short:         "Some commands for converting songs from Markdown with YAML frontmatter over to html",
	SilenceErrors: true, // avoids double printing of errors when thrown
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.WriteFatalf("Error: %v\n", err)
	}
}

func init() {
	rootCmd.SetOut(os.Stdout)
}

func writeToFileOrStdOut(content, outputFile string) {
	if strings.TrimSpace(outputFile) != "" {
		err := filehandler.WriteFileContents(outputFile, content)

		if err != nil {
			logger.WriteFatal(err.Error())
		}
	} else {
		logger.WriteInfo(content)
	}
}
