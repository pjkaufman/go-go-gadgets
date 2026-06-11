package markdown

import (
	"strings"

	"github.com/spf13/cobra"
)

const baseLevel = 3

func RootToMd(rootCmd *cobra.Command) string {
	var commandMd, commandToc strings.Builder

	var level = baseLevel
	if rootCmd.Runnable() {
		commandToc.WriteString("- [" + rootCmd.Name() + "](#" + rootCmd.Name() + "-base-command)\n")
		CommandToMd(rootCmd, &commandMd, level)

		level++
	}

	handleSubCommands(rootCmd, &commandMd, &commandToc, level)

	return commandToc.String() + "\n" + commandMd.String()
}

func handleSubCommands(cmd *cobra.Command, commandMd, commandToc *strings.Builder, level int) {
	for _, subCmd := range cmd.Commands() {
		if subCmd.Hidden || subCmd.Name() == "completion" || subCmd.Name() == "help" {
			continue
		}

		commandToc.WriteString(strings.Repeat(" ", (level-baseLevel)*2) + "- [" + subCmd.Name() + "](#" + subCmd.Name() + ")\n")

		CommandToMd(subCmd, commandMd, level)

		if len(subCmd.Commands()) != 0 {
			handleSubCommands(subCmd, commandMd, commandToc, level+1)
		}
	}
}
