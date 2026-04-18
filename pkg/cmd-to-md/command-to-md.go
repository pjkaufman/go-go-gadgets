package cmdtomd

import (
	"strings"

	"github.com/spf13/cobra"
)

func CommandToMd(cmd *cobra.Command, builder *strings.Builder, level int) {
	if cmd == nil {
		return
	}

	var (
		name            = cmd.Name()
		longDescription = cmd.Long
		example         = cmd.Example
	)

	builder.WriteString(strings.Repeat("#", level) + " " + name + "\n\n")
	if longDescription != "" {
		builder.WriteString(longDescription)
	} else {
		builder.WriteString(cmd.Short)
	}

	builder.WriteString("\n\n")

	// This may look useless, but it pulls in persistent flags from the parent, so it is called here
	// and then flags gets populated with both the inherited and regular flags.
	// It claims it does not modify the current flags, but that is not what I am seeing.
	cmd.InheritedFlags()

	var (
		flags = cmd.Flags()
	)
	if flags != nil && flags.HasFlags() {
		builder.WriteString(strings.Repeat("#", level+1) + " Flags\n\n")
		FlagsToMd(flags, builder)
		builder.WriteString("\n\n")
	}

	if example != "" {
		builder.WriteString(strings.Repeat("#", level+1) + " Usage\n\n``` bash\n")
		builder.WriteString(strings.TrimRight(strings.ReplaceAll(example, "To ", "# To "), "\n"))
		builder.WriteString("\n```\n\n")
	}
}
