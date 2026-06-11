//go:build generate_doc

package cmd

import "github.com/pjkaufman/go-go-gadgets/pkg/cli"

const (
	title       = "Versy"
	description = "This is a program that grabs the verse of the day or the specified verse(s) in two translations (either the default or the user specified ones)."
)

func init() {
	cli.AddGenerateCmd(rootCmd, title, description, []string{}, nil)
}
