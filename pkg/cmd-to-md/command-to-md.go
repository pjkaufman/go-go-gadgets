package cmdtomd

import (
	"strings"

	"github.com/spf13/cobra"
)

func CommandToMd(cmd *cobra.Command, builder *strings.Builder) {
	if cmd == nil {
		return
	}

	var (
		name            = cmd.Name()
		longDescription = cmd.Long
		example         = cmd.Example
	)

	builder.WriteString("### " + name + "\n\n")
	if longDescription != "" {
		builder.WriteString(longDescription)
	} else {
		builder.WriteString(cmd.Short)
	}

	builder.WriteString("\n\n")

	var flags = cmd.Flags()
	if flags != nil && flags.HasFlags() {
		builder.WriteString("#### Flags\n\n")
		FlagsToMd(flags, builder)
		builder.WriteString("\n\n")
	}

	if example != "" {
		builder.WriteString("#### Usage\n\n``` bash\n")
		builder.WriteString(strings.ReplaceAll(example, "To ", "# To "))
		builder.WriteString("\n```\n\n")
	}
}
