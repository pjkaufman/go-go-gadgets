//go:build generate_doc

package cmd

import "github.com/pjkaufman/go-go-gadgets/pkg/cli"

const (
	title       = "Song Converter"
	description = "This is a program that helps converter some Markdown files with YAML frontmatter into html or csv to help with creating a song book."
)

func init() {
	cli.AddGenerateCmd(rootCmd, title, description, []string{}, nil)
}
