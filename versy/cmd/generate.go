//go:build generate_doc

package cmd

import (
	cmdhandler "github.com/pjkaufman/go-go-gadgets/pkg/cmd-handler"
)

const (
	title       = "Versy"
	description = "This is a program that grabs the verse of the day or the specified verse(s) in two translations (either the default or the user specified ones)."
)

func init() {
	cmdhandler.AddGenerateCmd(rootCmd, title, description, []string{}, nil)
}
