package markdown

import (
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	TableHeader = `| Short Name | Long Name | Description | Value Type | Default Value | Is Required | Other Notes |`
	Separator   = `| ---------- | --------- | ----------- | ---------- | ------------- | ----------- | ----------- |`
)

func FlagsToMd(cliFlags *pflag.FlagSet, builder *strings.Builder) {
	if cliFlags == nil {
		return
	}

	builder.WriteString(TableHeader + "\n")
	builder.WriteString(Separator)

	cliFlags.VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		builder.WriteString("\n| ")
		builder.WriteString(flag.Shorthand)
		builder.WriteString(" | ")
		builder.WriteString(flag.Name)
		builder.WriteString(" | ")

		flagType, usage := pflag.UnquoteUsage(flag)

		builder.WriteString(usage)
		builder.WriteString(" | ")
		builder.WriteString(flagType)
		builder.WriteString(" | ")
		builder.WriteString(flag.DefValue)
		builder.WriteString(" | ")

		if val, ok := flag.Annotations[cobra.BashCompOneRequiredFlag]; ok && len(val) > 0 {
			builder.WriteString(val[0])
		} else {
			builder.WriteString("false")
		}

		builder.WriteString(" | ")

		if val, ok := flag.Annotations[cobra.BashCompFilenameExt]; ok && len(val) > 0 {
			builder.WriteString("Should be a file with one of the following extensions: ")
			builder.WriteString(strings.Join(val, ", "))
		} else if _, ok := flag.Annotations[cobra.BashCompSubdirsInDir]; ok {
			builder.WriteString("Should be a directory")
		} else if val, ok := flag.Annotations[flags.CustomEnumValuesFlagAnnotation]; ok && len(val) > 0 {
			builder.WriteString("Should be a one of the following: ")
			builder.WriteString(strings.Join(val, ", "))
		}

		builder.WriteString(" |")
	})
}
