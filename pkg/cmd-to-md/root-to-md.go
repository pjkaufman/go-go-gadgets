package cmdtomd

import (
	"strings"

	"github.com/spf13/cobra"
)

func RootToMd(rootCmd *cobra.Command) string {
	var builder strings.Builder

	handleSubCommands(rootCmd, &builder, 3)

	return builder.String()
}

func handleSubCommands(cmd *cobra.Command, builder *strings.Builder, level int) {
	for _, subCmd := range cmd.Commands() {
		if subCmd.Hidden || subCmd.Name() == "completion" {
			continue
		}

		CommandToMd(subCmd, builder, level)

		if len(subCmd.Commands()) != 0 {
			handleSubCommands(subCmd, builder, level+1)
		}
	}
}
