package filehandler

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

const (
	tempZip = "compress.zip"
	// have to use these or similar permissions to avoid permission denied errors in some cases
	folderPerms fs.FileMode = 0755
	numWorkers  int         = 5
)

// UnzipRunOperationAndRezip starts by deleting the destination directory if it exists,
// then it goes ahead an unzips the contents into the destination directory
// once that is done it runs the operation func on the destination folder
// lastly it rezips the folder back to compress.zip
func UnzipRunOperationAndRezip(src, dest string, operation func() error) error {
	dest = filepath.Clean(dest)

	err := os.RemoveAll(dest)
	if err != nil {
		return fmt.Errorf("failed to delete the destination directory %q: %w", dest, err)
	}

	err = Unzip(src, dest)
	if err != nil {
		return fmt.Errorf("failed to unzip %q: %w", src, err)
	}

	err = operation()
	if err != nil {
		return err
	}

	err = Rezip(dest, tempZip)
	if err != nil {
		return fmt.Errorf("failed to rezip content for source %q: %w", src, err)
	}

	err = os.RemoveAll(dest)
	if err != nil {
		return fmt.Errorf("failed to cleanup the destination directory %q: %w", dest, err)
	}

	err = MustRename(src, src+".original")
	if err != nil {
		return err
	}

	err = MustRename(tempZip, src)
	if err != nil {
		return err
	}

	return nil
}

// Unzip is based on https://stackoverflow.com/a/24792688
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		r.Close()
	}()

	err = os.MkdirAll(dest, folderPerms)
	if err != nil {
		return err
	}

	files := make(chan *zip.File, len(r.File))
	g, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for {
				select {
				case file, ok := <-files:
					if ok {
						wErr := extractAndWriteFile(dest, file)

						if wErr != nil {
							return wErr
						}
					} else {
						return nil
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	}

	for _, f := range r.File {
		files <- f
	}

	close(files)

	return g.Wait()
}

func extractAndWriteFile(dest string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		rc.Close()
	}()

	path := filepath.Join(dest, f.Name)

	// Check for ZipSlip (Directory traversal)
	if !strings.HasPrefix(path, dest+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if f.FileInfo().IsDir() {
		err = os.MkdirAll(path, folderPerms)

		if err != nil {
			return err
		}
	} else {
		err = os.MkdirAll(filepath.Dir(path), folderPerms)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			f.Close()
		}()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

// Rezip is based on https://stackoverflow.com/a/63233911
func Rezip(src, dest string) error {
	file, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	var mimetypePath = src + string(os.PathSeparator) + "mimetype"
	err = CopyMimetypeToZip(w, src, mimetypePath)
	if err != nil {
		return err
	}

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip empty directories
		if info.IsDir() {
			return nil
		}

		if mimetypePath == path {
			return nil
		}

		err = writeToZip(w, src, path)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(src, walker)
	if err != nil {
		return err
	}

	return nil
}

func writeToZip(w *zip.Writer, src, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// need a zip relative path to avoid creating extra directories inside of the zip
	var zipRelativePath = strings.Replace(path, src+string(os.PathSeparator), "", 1)
	f, err := w.Create(zipRelativePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func CopyMimetypeToZip(w *zip.Writer, src, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// need a zip relative path to avoid creating extra directories inside of the zip
	var zipRelativePath = strings.Replace(path, src+string(os.PathSeparator), "", 1)
	f, err := w.CreateHeader(&zip.FileHeader{
		Name:   strings.ReplaceAll(zipRelativePath, string(os.PathSeparator), "/"),
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

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
	var reader = bytes.NewReader([]byte(contents))
	f, err := w.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	return nil
}

func WriteZipCompressedFile(w *zip.Writer, zipFile *zip.File) error {
	file, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := w.Create(zipFile.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
}

func WriteZipCompressedBytes(w *zip.Writer, filename string, data []byte) error {
	var reader = bytes.NewReader(data)
	f, err := w.CreateHeader(&zip.FileHeader{
		Name:   filename,
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(f, reader)
	if err != nil {
		return err
	}

	return nil
}

func WriteZipUncompressedFile(w *zip.Writer, zipFile *zip.File) error {
	file, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := w.CreateHeader(&zip.FileHeader{
		Name:   zipFile.Name,
		Method: zip.Store,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	return nil
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
