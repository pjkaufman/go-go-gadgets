//go:build generate

package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	cmdtomd "github.com/pjkaufman/go-go-gadgets/pkg/cmd-to-md"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var generationDir string

const (
	GenerationDirArgEmpty = "generation-dir must have a non-whitespace value"
	title                 = "Ebook Linter"
	description           = "This is a program that helps lint and make updates to ebooks."
)

type TmplData struct {
	CommandStrings string
	Todos          []string
	Description    string
	Title          string
	CustomValues   map[string]any
}

// generateCmd generates the readme for the command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates the readme file for the program",
	Example: heredoc.Doc(`
		ebook-lint generate -d ./ebook-lint
		will look for a file called README.md.tmpl and if it is found generate a readme based
		on that file and the file command info.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateGenerateFlags(generationDir)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FolderMustExist(generationDir, "generation-dir")
		if err != nil {
			logger.WriteError(err.Error())
		}

		folders, err := filehandler.GetFoldersInCurrentFolder(filehandler.JoinPath(generationDir, "cmd"))
		if err != nil {
			logger.WriteError(err.Error())
		}

		tmpl, err := template.ParseFiles(filehandler.JoinPath(generationDir, "README.md.tmpl"))
		if err != nil {
			logger.WriteError(err.Error())
		}

		customValues := make(map[string]any)
		customValues["supportedFileTypes"] = folders

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		err = tmpl.Execute(writer, TmplData{
			CommandStrings: cmdtomd.RootToMd(rootCmd),
			Todos: []string{
				"See about removing unused files and images when running epub linting",
			},
			Title:        title,
			Description:  description,
			CustomValues: customValues,
		})
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.WriteFileContents(filehandler.JoinPath(generationDir, "README.md"), b.String())
		if err != nil {
			logger.WriteError(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&generationDir, "generation-dir", "g", "", "the path to the base folder of the ebook-lint program source code")
	err := generateCmd.MarkFlagRequired("generation-dir")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "generation-dir" as required on generate command: %v`, err))
	}

	// keep from showing up in the output of the command generation
	generateCmd.Hidden = true
}

func ValidateGenerateFlags(generationDir string) error {
	if strings.TrimSpace(generationDir) == "" {
		return errors.New(GenerationDirArgEmpty)
	}

	return nil
}
