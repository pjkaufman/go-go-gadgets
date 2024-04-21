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
	Use:   "song-converter",
	Short: "Some commands for converting songs from Markdown with YAML frontmatter over to html",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

func writeToFileOrStdOut(content, outputFile string) {
	if strings.TrimSpace(outputFile) != "" {
		filehandler.WriteFileContents(outputFile, content)
	} else {
		logger.WriteInfo(content)
	}
}
