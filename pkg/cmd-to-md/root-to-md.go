package cmdtomd

import (
	"strings"

	"github.com/spf13/cobra"
)

const baseLevel = 3

func RootToMd(rootCmd *cobra.Command) string {
	var commandMd, commandToc strings.Builder

	handleSubCommands(rootCmd, &commandMd, &commandToc, baseLevel)

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
