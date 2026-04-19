package cmd

import (
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Deals with creating files from the song Markdown files",
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&stagingDir, "working-dir", "d", "", "the directory where the Markdown files are located")
	err := createCmd.MarkPersistentFlagRequired("working-dir")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"working-dir\" as required on create command: %v\n", err)
	}

	err = createCmd.MarkPersistentFlagDirname("working-dir")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"working-dir\" as a directory on create command: %v\n", err)
	}
}
