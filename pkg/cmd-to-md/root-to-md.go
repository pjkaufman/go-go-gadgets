package cmdtomd

import (
	"strings"

	"github.com/spf13/cobra"
)

func RootToMd(rootCmd *cobra.Command) string {
	var builder strings.Builder

	handleSubCommands(rootCmd, &builder)

	return builder.String()
}

func handleSubCommands(cmd *cobra.Command, builder *strings.Builder) {
	for _, subCmd := range cmd.Commands() {
		CommandToMarkdown(subCmd, builder)

		if len(subCmd.Commands()) != 0 {
			handleSubCommands(subCmd, builder)
		}
	}
}
