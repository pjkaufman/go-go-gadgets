//go:build generate

package cmd

import (
	cmdhandler "github.com/pjkaufman/go-go-gadgets/pkg/cmd-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
)

const (
	title       = "Ebook Linter"
	description = "This is a program that helps lint and make updates to ebooks."
)

func getListOfAvailableTypes(generationDir string) (map[string]any, error) {
	folders, err := filehandler.GetFoldersInCurrentFolder(filehandler.JoinPath(generationDir, "cmd"))
	if err != nil {
		return nil, err
	}

	customValues := make(map[string]any)
	customValues["supportedFileTypes"] = folders

	return customValues, nil
}

func init() {
	cmdhandler.AddGenerateCmd(rootCmd, title, description, []string{
		"See about removing unused files and images when running epub linting",
	}, getListOfAvailableTypes)
}
