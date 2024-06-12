package epubhandler

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func UpdateEpub(src string, operation func(map[string]*zip.File, *zip.Writer, EpubInfo, string) []string) error {
	zipFiles, err := filehandler.GetFilesFromZip(src)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to get zip contents for %q: %s", src, err))
	}

	var (
		opfFilename string
		opfFile     *zip.File
	)
	for filename, file := range zipFiles {
		if strings.HasSuffix(filename, "opf") {
			opfFilename = filename
			opfFile = file
			break
		}
	}

	if opfFile == nil {
		return fmt.Errorf("failed to find the opf file for %q", src)
	}

	fileContents, err := filehandler.ReadInZipFileContents(opfFile)
	if err != nil {
		return err
	}

	epubInfo, err := ParseOpfFile(fileContents)
	if err != nil {
		logger.WriteError(fmt.Sprintf("Failed to parse %q for %q: %s", opfFilename, src, err))
	}
	var opfFolder = filehandler.GetFileFolder(opfFilename)

	var tempEpub = src + ".temp"
	var runOperation = func() error {
		tempEpubFile, err := os.Create(tempEpub)
		if err != nil {
			return fmt.Errorf("failed to create temporary epub file %q for %q: %w", tempEpub, src, err)
		}

		defer tempEpubFile.Close()

		w := zip.NewWriter(tempEpubFile)
		defer w.Close()

		if mimetypeFile, ok := zipFiles["mimetype"]; ok {
			err = filehandler.WriteZipUncompressedFile(w, mimetypeFile)
			if err != nil {
				return fmt.Errorf("failed to copy mimetype to zip file")
			}
		} else {
			return fmt.Errorf("no mimetype exists for %q", src)
		}

		filesHandled := operation(zipFiles, w, epubInfo, opfFolder)

		var handled bool
		for filename, zipFile := range zipFiles {
			handled = false
			for _, handledFile := range filesHandled {
				if filename == handledFile {
					handled = true
					break
				}
			}

			if handled {
				continue
			}

			err = filehandler.WriteZipCompressedFile(w, zipFile)
			if err != nil {
				return fmt.Errorf("failed to write file %q to zip for %q", zipFile.Name, src)
			}
		}

		return nil
	}

	err = runOperation()
	if err != nil {
		return err
	}

	filehandler.MustRename(src, src+".original")
	filehandler.MustRename(tempEpub, src)

	return nil
}
