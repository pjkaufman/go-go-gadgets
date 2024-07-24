//go:build generate

package cmd

import (
	cmdhandler "github.com/pjkaufman/go-go-gadgets/pkg/cmd-handler"
)

const (
	title       = "Song Converter"
	description = "This is a program that helps converter some Markdown files with YAML frontmatter into html or csv to help with creating a song book."
)

func init() {
	cmdhandler.AddGenerateCmd(rootCmd, title, description, []string{}, nil)
}
