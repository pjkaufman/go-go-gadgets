//go:build generate_doc

package cmd

import "github.com/pjkaufman/go-go-gadgets/pkg/cli"

const (
	title       = "Epub Linter"
	description = "This is a program that helps lint and make updates to epubs. The logic is designed to be broken into different functionality for each command, so a user does not have to use all functionality if they do not want to."
)

func init() {
	cli.AddGenerateCmd(rootCmd, title, description, []string{
		"See about removing unused files and images when running epub linting",
	}, nil)
}
