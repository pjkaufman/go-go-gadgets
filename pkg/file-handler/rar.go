package filehandler

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"

	archiver "github.com/mholt/archiver/v3"
)

func ConvertRarToCbz(src string) error {
	var (
		err     error
		rar     = archiver.NewRar()
		cbzPath = strings.Replace(src, ".cbr", ".cbz", 1)
	)
	cbzFile, err := os.Create(cbzPath)
	if err != nil {
		return fmt.Errorf("failed to create cbz file %q for %q: %w", cbzPath, src, err)
	}

	defer cbzFile.Close()

	w := zip.NewWriter(cbzFile)
	defer w.Close()

	return rar.Walk(src, func(f archiver.File) error {
		if f.IsDir() {
			return nil
		}

		defer f.Close()

		zf, err := w.Create(f.Name())
		if err != nil {
			return err
		}

		_, err = io.Copy(zf, f)
		if err != nil {
			return fmt.Errorf(`could not convert rar file to zip file for %q: %w`, f.Name(), err)
		}

		return nil
	})
}
