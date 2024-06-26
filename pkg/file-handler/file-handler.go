package filehandler

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const bytesInAKiloByte float64 = 1024

func FileExists(path string) (bool, error) {
	if strings.TrimSpace(path) == "" {
		return false, nil
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf(`could not verify that %q exists: %s`, path, err)
	}

	return true, nil
}

func FileMustExist(path, name string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("%s must have a non-whitespace value", name)
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: %q must exist", name, path)
		}

		return fmt.Errorf(`could not verify that %q exists: %s`, path, err)
	}

	return nil
}

func FolderExists(path string) (bool, error) {
	if strings.TrimSpace(path) == "" {
		return false, nil
	}

	folderInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf(`could not verify that %q exists and is a directory: %s`, path, err)
	}

	if !folderInfo.IsDir() {
		return false, nil
	}

	return true, nil
}

func FolderMustExist(path, name string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("%s must have a non-whitespace value", name)
	}

	folderInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: %q must exist", name, path)
		}

		return fmt.Errorf(`could not verify that %q exists and is a directory: %w`, path, err)
	}

	if !folderInfo.IsDir() {
		return fmt.Errorf("%s: %q must be a folder", name, path)
	}

	return nil
}

func GetFoldersInCurrentFolder(path string) ([]string, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf(`could not get files/folders in %q: %w`, path, err)
	}

	var actualDirs []string
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		actualDirs = append(actualDirs, dir.Name())
	}

	return actualDirs, nil
}

func GetFileFolder(filePath string) string {
	if strings.TrimSpace(filePath) == "" {
		return ""
	}

	return path.Join(filePath, "..")
}

func JoinPath(elements ...string) string {
	return path.Join(elements...)
}

func ReadInFileContents(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", nil
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf(`could not read in file contents for %q: %w`, path, err)
	}

	return string(fileBytes), nil
}

func ReadInBinaryFileContents(path string) ([]byte, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(`could not read in file contents for %q: %w`, path, err)
	}

	return fileBytes, nil
}

func WriteFileContents(path, content string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	var fileMode fs.FileMode

	fileInfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fileMode = fs.ModePerm
		} else {
			return fmt.Errorf(`could not read in existing file info to retain existing permission for %q: %w`, path, err)
		}
	} else {
		fileMode = fileInfo.Mode()
	}

	err = os.WriteFile(path, []byte(content), fileMode)
	if err != nil {
		return fmt.Errorf(`could not write to file %q: %w`, path, err)
	}

	return nil
}

func WriteBinaryFileContents(path string, content []byte) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	var fileMode fs.FileMode

	fileInfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fileMode = fs.ModePerm
		} else {
			return fmt.Errorf(`could not read in existing file info to retain existing permission for %q: %w`, path, err)
		}
	} else {
		fileMode = fileInfo.Mode()
	}

	err = os.WriteFile(path, content, fileMode)
	if err != nil {
		return fmt.Errorf(`could not write to file %q: %w`, path, err)
	}

	return nil
}

func MustGetAllFilesWithExtInASpecificFolder(dir, ext string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf(`failed to read in folder %q: %w`, dir, err)
	}

	var fileList []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ext) {
			fileList = append(fileList, f.Name())
		}
	}

	return fileList, nil
}

// based on https://stackoverflow.com/a/67629473
func MustGetAllFilesWithExtsInASpecificFolderAndSubFolders(dir string, exts ...string) ([]string, error) {
	var a []string
	err := filepath.WalkDir(dir, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		var name = strings.ToLower(d.Name())
		if len(exts) == 1 {
			if strings.HasSuffix(name, exts[0]) {
				a = append(a, s)
			}
		} else if fileHasOneOfExts(name, exts) {
			a = append(a, s)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk dir %q: %w", dir, err)
	}

	return a, nil
}

func fileHasOneOfExts(fileName string, exts []string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}

func MustRename(src, dest string) error {
	err := os.Rename(src, dest)

	if err != nil {
		return fmt.Errorf("failed to rename %q to %q: %w", src, dest, err)
	}

	return nil
}

func MustGetFileSize(path string) (float64, error) {
	if strings.TrimSpace(path) == "" {
		return 0, fmt.Errorf("to get a file's size it must have a non-empty path")
	}

	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf(`%q does not exist so the the file size cannot be retrieved`, path)
		}

		return 0, fmt.Errorf(`could not verify that %q exists to check its size: %w`, path, err)
	}

	return float64(f.Size()) / bytesInAKiloByte, nil
}

func MustCreateFolderIfNotExists(path string) error {
	folderExists, err := FolderExists(path)
	if err != nil {
		return err
	}

	if folderExists {
		return nil
	}

	err = os.MkdirAll(path, folderPerms)
	if err != nil {
		return fmt.Errorf("failed to create folder(s) for path %q: %s", path, err)
	}

	return nil
}
