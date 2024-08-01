package epub

import (
	"archive/zip"
	"errors"
	"fmt"
	"strings"
)

const (
	EpubPathArgEmpty   = "epub-file must have a non-whitespace value"
	EpubPathArgNonEpub = "epub-file must be an Epub file"
	cliLineSeparator   = "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
)

var epubFile string

func validateFilesExist(opfFolder string, files map[string]struct{}, zipFiles map[string]*zip.File) error {
	for file := range files {
		var filePath = getFilePath(opfFolder, file)

		if _, ok := zipFiles[filePath]; !ok {
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
