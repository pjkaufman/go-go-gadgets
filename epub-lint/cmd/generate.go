//go:build generate_doc

package cmd

import (
	cmdhandler "github.com/pjkaufman/go-go-gadgets/pkg/cmd-handler"
)

const (
	title       = "Epub Linter"
	description = "This is a program that helps lint and make updates to epubs."
)

func init() {
	cmdhandler.AddGenerateCmd(rootCmd, title, description, []string{
		"See about removing unused files and images when running epub linting",
	}, nil)
}
