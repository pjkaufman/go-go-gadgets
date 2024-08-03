package ziphandler

import (
	"archive/zip"
	"fmt"
	"os"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
)

func UpdateZip(src string, operation func(map[string]*zip.File, *zip.Writer) ([]string, error)) error {
	r, zipFiles, err := filehandler.GetFilesFromZip(src)
	if err != nil {
		return fmt.Errorf("failed to get zip contents for %q: %w", src, err)
	}

	defer r.Close()

	var tempZip = src + ".temp"
	var runOperation = func() error {
		tempZipFile, err := os.Create(tempZip)
		if err != nil {
			return fmt.Errorf("failed to create temporary zip file %q for %q: %w", tempZip, src, err)
		}

		defer tempZipFile.Close()

		w := zip.NewWriter(tempZipFile)
		defer w.Close()

		if mimetypeFile, ok := zipFiles["mimetype"]; ok {
			err = filehandler.WriteZipUncompressedFile(w, mimetypeFile)

			if err != nil {
				return fmt.Errorf("failed to copy mimetype to zip file")
			}
		}

		filesHandled, err := operation(zipFiles, w)
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

	err = filehandler.Rename(tempZip, src)
	if err != nil {
		return err
	}

	return nil
}
