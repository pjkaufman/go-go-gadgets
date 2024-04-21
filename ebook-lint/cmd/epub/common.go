package epub

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

const (
	EpubPathArgEmpty   = "epub-file must have a non-whitespace value"
	EpubPathArgNonEpub = "epub-file must be an Epub file"
	cliLineSeparator   = "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
)

var epubFile string

func getEpubInfo(dir, epubName string) (string, linter.EpubInfo) {
	opfFiles := filehandler.MustGetAllFilesWithExtsInASpecificFolderAndSubFolders(dir, ".opf")
	if len(opfFiles) < 1 {
		logger.WriteError(fmt.Sprintf("did not find opf file for \"%s\"", epubName))
	}

	var opfFile = opfFiles[0]
	opfText := filehandler.ReadInFileContents(opfFile)

	epubInfo, err := linter.ParseOpfFile(opfText)
	if err != nil {
		logger.WriteError(fmt.Sprintf("Failed to parse \"%s\" for \"%s\": %s", opfFile, epubName, err))
	}

	var opfFolder = filehandler.GetFileFolder(opfFile)

	return opfFolder, epubInfo
}

func validateFilesExist(opfFolder string, files map[string]struct{}) {
	for file := range files {
		var filePath = getFilePath(opfFolder, file)

		if !filehandler.FileExists(filePath) {
			logger.WriteError(fmt.Sprintf(`file from manifest not found: "%s" must exist`, filePath))
		}
	}
}

func validateCommonEpubFlags(epubPath string) error {
	if strings.TrimSpace(epubPath) == "" {
		return errors.New(EpubPathArgEmpty)
	}

	if !strings.HasSuffix(epubPath, ".epub") {
		return errors.New(EpubPathArgNonEpub)
	}

	return nil
}
