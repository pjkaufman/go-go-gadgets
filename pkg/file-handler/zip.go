package filehandler

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

const (
	tempZip                 = "compress.zip"
	folderPerms fs.FileMode = 0755
)

// UnzipRunOperationAndRezip starts by deleting the destination directory if it exists,
// then it goes ahead an unzips the contents into the destination directory
// once that is done it runs the operation func on the destination folder
// lastly it rezips the folder back to compress.zip
func UnzipRunOperationAndRezip(src, dest string, operation func()) {
	var err error
	if FolderExists(dest) {
		err = os.RemoveAll(dest)

		if err != nil {
			logger.WriteError(fmt.Sprintf("failed to delete the destination directory \"%s\": %s", dest, err))
		}
	}

	err = Unzip(src, dest)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to unzip \"%s\": %s", src, err))
	}

	operation()

	err = Rezip(dest, tempZip)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to rezip content for source \"%s\": %s", src, err))
	}

	err = os.RemoveAll(dest)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to cleanup the destination directory \"%s\": %s", dest, err))
	}

	MustRename(src, src+".original")
	MustRename(tempZip, src)
}

// Unzip is based on https://stackoverflow.com/a/24792688
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	// have to use these or similar permissions to avoid permission denied errors in some cases
	var folderPerms fs.FileMode = 0755
	err = os.MkdirAll(dest, folderPerms)
	if err != nil {
		return err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
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
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)

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

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip empty directories
		if info.IsDir() {
			return nil
		}

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
	err = filepath.Walk(src, walker)
	if err != nil {
		return err
	}

	return nil
}

func UnzipSource(source, destination string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
