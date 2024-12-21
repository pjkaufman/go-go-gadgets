package filehandler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetFilesFromZip(src string) (*zip.ReadCloser, map[string]*zip.File, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return nil, nil, err
	}

	var zipFiles = make(map[string]*zip.File, 0)
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}

		zipFiles[f.Name] = f
	}

	return r, zipFiles, nil
}

func WriteZipCompressedString(w *zip.Writer, filename, contents string) error {
	return compressedWriteToZip(w, bytes.NewReader([]byte(contents)), filename)
}

func WriteZipUncompressedString(w *zip.Writer, filename, contents string) error {
	return uncompressedWriteToZip(w, bytes.NewReader([]byte(contents)), filename)
}

func WriteZipCompressedFile(w *zip.Writer, zipFile *zip.File) error {
	file, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	return compressedWriteToZip(w, file, zipFile.Name)
}

func WriteZipCompressedBytes(w *zip.Writer, filename string, data []byte) error {
	return compressedWriteToZip(w, bytes.NewReader(data), filename)
}

func WriteZipUncompressedFile(w *zip.Writer, zipFile *zip.File) error {
	file, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	return uncompressedWriteToZip(w, file, zipFile.Name)
}

func ReadInZipFileContents(zipFile *zip.File) (string, error) {
	file, err := zipFile.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var fileBytes = &bytes.Buffer{}
	_, err = io.Copy(fileBytes, file)
	if err != nil {
		return "", fmt.Errorf(`could not read in zip file contents for %q: %w`, zipFile.Name, err)
	}

	return fileBytes.String(), nil
}

func ReadInZipFileBytes(zipFile *zip.File) ([]byte, error) {
	file, err := zipFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileBytes = &bytes.Buffer{}
	_, err = io.Copy(fileBytes, file)
	if err != nil {
		return nil, fmt.Errorf(`could not read in zip file bytes for %q: %w`, zipFile.Name, err)
	}

	return fileBytes.Bytes(), nil
}

func compressedWriteToZip(w *zip.Writer, reader io.Reader, filename string) error {
	f, err := w.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	return err
}

func uncompressedWriteToZip(w *zip.Writer, reader io.Reader, filename string) error {
	f, err := w.CreateHeader(&zip.FileHeader{
		Name:   filename,
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	return err
}

func UnzipFile(zipFile *os.File, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer zipReader.Close()

	if err := os.MkdirAll(destDir, folderPerms); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, file := range zipReader.File {
		destPath := filepath.Join(destDir, file.Name)

		// Check for zip slip vulnerability
		if !strings.HasPrefix(destPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path in zip: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), folderPerms); err != nil {
			return fmt.Errorf("failed to create directory for file %s: %w", destPath, err)
		}

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
		}

		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			return fmt.Errorf("failed to open zip file %s: %w", file.Name, err)
		}

		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()
		if err != nil {
			return fmt.Errorf("failed to copy file contents from %s: %w", file.Name, err)
		}
	}

	return nil
}
