package cmd

import (
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Deals with creating files from the song Markdown files",
}

var createFlags = flags.Flags{
	Flags: []flags.Flag{
		flags.NewDirectoryFlag(true, true, &stagingDir, "working-dir", "d", "", "the directory where the Markdown files are located"),
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	err := createFlags.AddToCmd(createCmd)
	if err != nil {
		logger.WriteFatal(err.Error())
	}
}
