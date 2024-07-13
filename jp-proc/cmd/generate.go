//go:build generate

package cmd

import (
	cmdhandler "github.com/pjkaufman/go-go-gadgets/pkg/cmd-handler"
)

const (
	title       = "Jpeg and Png Processor"
	description = `This is meant to be a replacement for my usage of imgp.

Currently I use imgp for the following things:
- image resizing
- exif data removal
- image quality setting

Given how this works, I find it easier to just go ahead and do a simple program in Go to see how things stack up and not be so reliant on Python. This also helps me learn some more about imaging processing as well. So a win-win in my book.`
)

func init() {
	cmdhandler.AddGenerateCmd(rootCmd, title, description, []string{
		"Resize png test",
	}, nil)
}
