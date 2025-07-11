package cmd

import (
	"archive/zip"
	"errors"
	"fmt"
	"strings"
)

const (
	cliLineSeparator = "-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
)

var (
	ErrEpubPathArgEmpty   = errors.New("epub-file must have a non-whitespace value")
	ErrEpubPathArgNonEpub = errors.New("epub-file must be an Epub file")
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
		return ErrEpubPathArgEmpty
	}

	if !strings.HasSuffix(epubPath, ".epub") {
		return ErrEpubPathArgNonEpub
	}

	return nil
}
