//go:build generate_doc

package cmd

import (
	cmdhandler "github.com/pjkaufman/go-go-gadgets/pkg/cmd-handler"
)

const (
	title       = "Cat ASCII"
	description = "This is a program that prints out cat ASCII to the terminal. Just calling `cat-ascii` will print out cat ASCII for you to enjoy."
)

func init() {
	cmdhandler.AddGenerateCmd(rootCmd, title, description, []string{}, nil)
}
