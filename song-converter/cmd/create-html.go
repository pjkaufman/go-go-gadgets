package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"
	"github.com/spf13/cobra"
)

const (
	CoverPathArgEmpty  = "cover-file must have a non-whitespace value"
	CoverPathNotMdFile = "cover-file must be an md file"
	StagingDirArgEmpty = "working-dir must have a non-whitespace value"
	fileFormat         = `<html>
  <body>
    <section id="cover">
      %s
    </section>
    <section id="contents">
    <h1 class="toc">Table of Contents</h1>
    <ul>
        %s
    </ul>
    </section>
    <section id="songs">
      %s
    </section>
  </body>
</html>`
)

var (
	stagingDir         string
	bodyHtmlOutputFile string
	coverOutputFile    string
	coverInputFilePath string
)

// CreateHtmlCmd represents the CreateSongs command
var CreateHtmlCmd = &cobra.Command{
	Use:   "create-html",
	Short: "Converts the cover and all Markdown files in the specified folder into html in alphabetical order generating three sections: the cover, table of contents, and songs",
	Example: heredoc.Doc(`To write the output of converting the files in the specified folder to html to a file:
	song-converter create-html -d working-dir -c cover.md -o songs.html

	To write the output of converting the files in the specified folder to html to std out:
	song-converter create-html -d working-dir -s cover.md
	`),
	Long: heredoc.Doc(`How it works:
	- Reads in all of the files in the specified folder
	- Sorts the files alphabetically
	- Adds the cover to the start of the content after converting it to html
	- Converts each file into html
	- Writes the content to the specified source
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateCreateHtmlFlags(stagingDir, coverInputFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FolderMustExist(stagingDir, "working-dir")
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FileMustExist(coverInputFilePath, "cover-file")
		if err != nil {
			logger.WriteError(err.Error())
		}

		var isWritingToFile = strings.TrimSpace(coverOutputFile) == ""
		if isWritingToFile {
			logger.WriteInfo("Converting file to html cover")
		}
		coverMd, err := filehandler.ReadInFileContents(coverInputFilePath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		coverHtml := converter.BuildHtmlCover(coverMd)

		if isWritingToFile {
			logger.WriteInfo("Finished creating html cover file")
		}

		if isWritingToFile {
			logger.WriteInfo("Converting Markdown files to html")
		}

		files := filehandler.MustGetAllFilesWithExtInASpecificFolder(stagingDir, ".md")
		sort.Strings(files)

		var mdInfo = make([]converter.MdFileInfo, len(files))

		for i, fileName := range files {
			var filePath = filehandler.JoinPath(stagingDir, fileName)
			fileContents, err := filehandler.ReadInFileContents(filePath)
			if err != nil {
				logger.WriteError(err.Error())
			}

			mdInfo[i] = converter.MdFileInfo{
				FilePath:     filePath,
				FileName:     fileName,
				FileContents: fileContents,
			}
		}

		songsHtml, headerIds, err := converter.BuildHtmlSongs(mdInfo)
		if err != nil {
			logger.WriteError(err.Error())
		}

		writeToFileOrStdOut(fmt.Sprintf(fileFormat, coverHtml, buildListItems(headerIds), songsHtml), bodyHtmlOutputFile)

		if isWritingToFile {
			logger.WriteInfo("Finished converting Markdown files to html")
		}
	},
}

func init() {
	rootCmd.AddCommand(CreateHtmlCmd)

	CreateHtmlCmd.Flags().StringVarP(&stagingDir, "working-dir", "d", "", "the directory where the Markdown files are located")
	err := CreateHtmlCmd.MarkFlagRequired("working-dir")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "working-dir" as required on create html command: %v`, err))
	}

	err = CreateHtmlCmd.MarkFlagDirname("working-dir")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "working-dir" as a directory on create html command: %v`, err))
	}

	CreateHtmlCmd.Flags().StringVarP(&coverInputFilePath, "cover-file", "c", "", "the markdown cover file to use")
	err = CreateHtmlCmd.MarkFlagRequired("cover-file")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "cover-file" as required on create html command: %v`, err))
	}

	err = CreateHtmlCmd.MarkFlagFilename("cover-file", "md")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "cover-file" as a looking for specific file types on create html command: %v`, err))
	}

	CreateHtmlCmd.Flags().StringVarP(&bodyHtmlOutputFile, "output", "o", "", "the html file to write the output to")
	err = CreateHtmlCmd.MarkFlagFilename("output", "html")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to mark flag "output" as a looking for specific file types on create html command: %v`, err))
	}
}

func ValidateCreateHtmlFlags(stagingDir, coverInputFilePath string) error {
	if strings.TrimSpace(stagingDir) == "" {
		return errors.New(StagingDirArgEmpty)
	}

	if strings.TrimSpace(coverInputFilePath) == "" {
		return errors.New(CoverPathArgEmpty)
	}

	if !strings.HasSuffix(coverInputFilePath, ".md") {
		return errors.New(CoverPathNotMdFile)
	}

	return nil
}

func buildListItems(headerIds []string) string {
	if len(headerIds) == 0 {
		return ""
	}

	var listItems = strings.Builder{}
	for _, headerId := range headerIds {
		listItems.WriteString(fmt.Sprintf(`<li><a href="#%s"></a></li>`+"\n", headerId))
	}

	return listItems.String()
}
