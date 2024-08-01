package cmdhandler

import (
	"bytes"
	"errors"
	"strings"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	cmdtomd "github.com/pjkaufman/go-go-gadgets/pkg/cmd-to-md"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var generationDir string

const GenerationDirArgEmpty = "generation-dir must have a non-whitespace value"

type TmplData struct {
	CommandStrings string
	Todos          []string
	Description    string
	Title          string
	CustomValues   map[string]any
}

func AddGenerateCmd(rootCmd *cobra.Command, title, description string, todos []string, getCustomValues func(string) (map[string]any, error)) {
	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates the readme file for the program",
		Example: heredoc.Doc(`
			` + rootCmd.Use + ` generate -d ./` + rootCmd.Use + `
			will look for a file called README.md.tmpl and if it is found generate a readme based
			on that file and the file command info.
		`),
		Run: func(cmd *cobra.Command, args []string) {
			err := ValidateGenerateFlags(generationDir)
			if err != nil {
				logger.WriteError(err.Error())
			}

			err = filehandler.FolderArgExists(generationDir, "generation-dir")
			if err != nil {
				logger.WriteError(err.Error())
			}

			tmpl, err := template.ParseFiles(filehandler.JoinPath(generationDir, "README.md.tmpl"))
			if err != nil {
				logger.WriteError(err.Error())
			}

			var customValues = make(map[string]any)
			if getCustomValues != nil {
				customValues, err = getCustomValues(generationDir)

				if err != nil {
					logger.WriteError(err.Error())
				}
			}

			var b bytes.Buffer

			err = tmpl.Execute(&b, TmplData{
				CommandStrings: cmdtomd.RootToMd(rootCmd),
				Todos:          todos,
				Description:    description,
				Title:          title,
				CustomValues:   customValues,
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

	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&generationDir, "generation-dir", "g", "", "the path to the base folder of the "+rootCmd.Use+" program source code")
	err := generateCmd.MarkFlagRequired("generation-dir")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"generation-dir\" as required on generate command: %v\n", err)
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
