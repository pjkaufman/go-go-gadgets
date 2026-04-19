package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

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
	versionDescriptor  string
	location           string
)

// CreateHtmlCmd represents the CreateSongs command
var CreateHtmlCmd = &cobra.Command{
	Use:   "html",
	Short: "Converts the cover and all Markdown files in the specified folder into html in alphabetical order generating three sections: the cover, table of contents, and songs",
	Example: heredoc.Doc(`To write the output of converting the files in the specified folder to html to a file:
	song-converter create html -d working-dir -c cover.md -o songs.html

	To write the output of converting the files in the specified folder to html to std out:
	song-converter create html -d working-dir -c cover.md
	`),
	Long: heredoc.Doc(`How it works:
	- Reads in all of the files in the specified folder
	- Sorts the files alphabetically
	- Adds the cover to the start of the content after converting it to html
	- Converts each file into html
	- Writes the content to the specified source
	`),
	PreRunE: validateCreateHtmlFile,
	Run: func(cmd *cobra.Command, args []string) {
		createHtmlFile(stagingDir, coverInputFilePath, coverOutputFile, bodyHtmlOutputFile, versionDescriptor, "", false)
	},
}

func createHtmlFile(stagingDir, coverInputFilePath, coverOutputFile, bodyHtmlOutputFile, bookType, extraCss string, isBook bool) {
	var isWritingToFile = strings.TrimSpace(coverOutputFile) == ""
	if isWritingToFile {
		logger.WriteInfo("Converting file to html cover")
	}
	coverMd, err := filehandler.ReadInFileContents(coverInputFilePath)
	if err != nil {
		logger.WriteError(err.Error())
	}

	coverHtml := converter.BuildHtmlCover(coverMd, bookType, extraCss, time.Now())

	if isWritingToFile {
		logger.WriteInfo("Finished creating html cover file")
	}

	if isWritingToFile {
		logger.WriteInfo("Converting Markdown files to html")
	}

	files, err := filehandler.GetAllFilesWithExtInASpecificFolder(stagingDir, ".md")
	if err != nil {
		logger.WriteError(err.Error())
	}

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

	var converterType = converter.Digital
	if isBook {
		mdInfo, err = converter.FilterAndSortSongs(mdInfo, location)
		if err != nil {
			logger.WriteError(err.Error())
		}

		converterType = converter.Book
	}

	songsHtml, headerIds, err := converter.BuildHtmlSongs(mdInfo, converterType)
	if err != nil {
		logger.WriteError(err.Error())
	}

	if isBook {
		writeToFileOrStdOut(fmt.Sprintf(bookFormat, coverHtml, buildBookListItems(mdInfo), songsHtml), bodyHtmlOutputFile)
	} else {
		writeToFileOrStdOut(fmt.Sprintf(fileFormat, coverHtml, buildListItems(headerIds), songsHtml), bodyHtmlOutputFile)
	}

	if isWritingToFile {
		logger.WriteInfo("Finished converting Markdown files to html")
	}
}

func validateCreateHtmlFile(cmd *cobra.Command, args []string) error {
	err := ValidateCreateHtmlAndBookFlags(stagingDir, coverInputFilePath)
	if err != nil {
		return err
	}

	err = filehandler.FolderArgExists(stagingDir, "working-dir")
	if err != nil {
		return err
	}

	err = filehandler.FileArgExists(coverInputFilePath, "cover-file")
	if err != nil {
		return err
	}

	return nil
}

func init() {
	createCmd.AddCommand(CreateHtmlCmd)

	createCommonHtmlAndBookFlags(CreateHtmlCmd)

	CreateHtmlCmd.Flags().StringVarP(&versionDescriptor, "version-type", "v", "", "the version descriptor for the type of songs to generate (generally just abridged or unabridged)")
	err := CreateHtmlCmd.MarkFlagRequired("version-type")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"version-type\" as required on create html command: %v\n", err)
	}
}

// createCommonHtmlAndBookFlags is meant to allow for de-duplicating the common flags for create html and create book
func createCommonHtmlAndBookFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&coverInputFilePath, "cover-file", "c", "", "the markdown cover file to use")
	err := cmd.MarkFlagRequired("cover-file")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"cover-file\" as required on create %s command: %v\n", cmd.Use, err)
	}

	err = cmd.MarkFlagFilename("cover-file", "md")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"cover-file\" as a looking for specific file types on create %s command: %v\n", cmd.Use, err)
	}

	cmd.Flags().StringVarP(&bodyHtmlOutputFile, "output", "o", "", "the html file to write the output to")
	err = cmd.MarkFlagFilename("output", "html")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"output\" as a looking for specific file types on create %s command: %v\n", cmd.Use, err)
	}
}

func ValidateCreateHtmlAndBookFlags(stagingDir, coverInputFilePath string) error {
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
		fmt.Fprintf(&listItems, `<li><a href="#%s"></a></li>`+"\n", headerId)
	}

	return listItems.String()
}
