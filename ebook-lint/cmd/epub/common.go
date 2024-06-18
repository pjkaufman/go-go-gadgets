package epub

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
)

const (
	EpubPathArgEmpty   = "epub-file must have a non-whitespace value"
	EpubPathArgNonEpub = "epub-file must be an Epub file"
	cliLineSeparator   = "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
)

var epubFile string

func getEpubInfo(dir, epubName string) (string, linter.EpubInfo, error) {
	opfFiles, err := filehandler.MustGetAllFilesWithExtsInASpecificFolderAndSubFolders(dir, ".opf")
	if err != nil {
		return "", linter.EpubInfo{}, err
	}

	if len(opfFiles) < 1 {
		return "", linter.EpubInfo{}, fmt.Errorf("did not find opf file for %q", epubName)
	}

	var opfFile = opfFiles[0]
	opfText, err := filehandler.ReadInFileContents(opfFile)
	if err != nil {
		return "", linter.EpubInfo{}, err
	}

	epubInfo, err := linter.ParseOpfFile(opfText)
	if err != nil {
		return "", linter.EpubInfo{}, fmt.Errorf("failed to parse %q for %q: %w", opfFile, epubName, err)
	}

	var opfFolder = filehandler.GetFileFolder(opfFile)

	return opfFolder, epubInfo, nil
}

func validateFilesExist(opfFolder string, files map[string]struct{}) error {
	for file := range files {
		var filePath = getFilePath(opfFolder, file)

		fileExists, err := filehandler.FileExists(filePath)
		if err != nil {
			return err
		}

		if !fileExists {
			return fmt.Errorf(`file from manifest not found: %q must exist`, filePath)
		}
	}

	return nil
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
