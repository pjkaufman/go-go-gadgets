package epubhandler

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
)

const defaultMimetypeContents = "application/epub+zip"

func UpdateEpub(src string, operation func(map[string]*zip.File, *zip.Writer, EpubInfo, string) ([]string, error)) error {
	r, zipFiles, err := filehandler.GetFilesFromZip(src)
	if err != nil {
		return fmt.Errorf("failed to get zip contents for %q: %w", src, err)
	}

	defer r.Close()

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
		return fmt.Errorf("failed to parse %q for %q: %s", opfFilename, src, err)
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
				return fmt.Errorf("failed to copy mimetype to zip file: %w", err)
			}
		} else {
			err = filehandler.WriteZipUncompressedString(w, "mimetype", defaultMimetypeContents)

			if err != nil {
				return fmt.Errorf("failed to add default mimetype to zip file: %w", err)
			}
		}

		filesHandled, err := operation(zipFiles, w, epubInfo, opfFolder)
		if err != nil {
			return err
		}

		filesHandled = append(filesHandled, "mimetype")

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

	// we are closing this here to make sure that the operations run correctly, but it needed a defer for possible errors, so we ignore the error in the defer
	err = r.Close()
	if err != nil {
		return fmt.Errorf("failed to close zip reader: %w", err)
	}

	err = filehandler.Rename(src, src+".original")
	if err != nil {
		return err
	}

	err = filehandler.Rename(tempEpub, src)
	if err != nil {
		return err
	}

	return nil
}
