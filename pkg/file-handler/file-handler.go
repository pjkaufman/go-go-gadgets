package filehandler

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

const bytesInAKiloByte float64 = 1024

func FileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		logger.WriteError(fmt.Sprintf(`could not verify that "%s" exists: %s`, path, err))
	}

	return true
}

func FileMustExist(path, name string) {
	if strings.TrimSpace(path) == "" {
		logger.WriteError(fmt.Sprintf("%s must have a non-whitespace value", name))
	}

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.WriteError(fmt.Sprintf("%s: \"%s\" must exist", name, path))
		}

		logger.WriteError(fmt.Sprintf(`could not verify that "%s" exists: %s`, path, err))
	}
}

func FolderExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}

	folderInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		logger.WriteError(fmt.Sprintf(`could not verify that "%s" exists and is a directory: %s`, path, err))
	}

	if !folderInfo.IsDir() {
		return false
	}

	return true
}

func FolderMustExist(path, name string) {
	if strings.TrimSpace(path) == "" {
		logger.WriteError(fmt.Sprintf("%s must have a non-whitespace value", name))
	}

	folderInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.WriteError(fmt.Sprintf("%s: \"%s\" must exist", name, path))
		}

		logger.WriteError(fmt.Sprintf(`could not verify that "%s" exists and is a directory: %s`, path, err))
	}

	if !folderInfo.IsDir() {
		logger.WriteError(fmt.Sprintf("%s: \"%s\" must be a folder", name, path))
	}
}

func GetFoldersInCurrentFolder(path string) []string {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	dirs, err := os.ReadDir(path)
	if err != nil {
		logger.WriteError(fmt.Sprintf(`could not get files/folders in "%s": %s`, path, err))
	}

	var actualDirs []string
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		actualDirs = append(actualDirs, dir.Name())
	}

	return actualDirs
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

func ReadInFileContents(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		logger.WriteError(fmt.Sprintf(`could not read in file contents for "%s": %s`, path, err))
	}

	return string(fileBytes)
}

func ReadInBinaryFileContents(path string) []byte {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		logger.WriteError(fmt.Sprintf(`could not read in file contents for "%s": %s`, path, err))
	}

	return fileBytes
}

func WriteFileContents(path, content string) {
	if strings.TrimSpace(path) == "" {
		return
	}

	var fileMode fs.FileMode

	fileInfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fileMode = fs.ModePerm
		} else {
			logger.WriteError(fmt.Sprintf(`could not read in existing file info to retain existing permission for "%s": %s`, path, err))
		}
	} else {
		fileMode = fileInfo.Mode()
	}

	err = os.WriteFile(path, []byte(content), fileMode)
	if err != nil {
		logger.WriteError(fmt.Sprintf(`could not write to file "%s": %s`, path, err))
	}
}

func WriteBinaryFileContents(path string, content []byte) {
	if strings.TrimSpace(path) == "" {
		return
	}

	var fileMode fs.FileMode

	fileInfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fileMode = fs.ModePerm
		} else {
			logger.WriteError(fmt.Sprintf(`could not read in existing file info to retain existing permission for "%s": %s`, path, err))
		}
	} else {
		fileMode = fileInfo.Mode()
	}

	err = os.WriteFile(path, content, fileMode)
	if err != nil {
		logger.WriteError(fmt.Sprintf(`could not write to file "%s": %s`, path, err))
	}
}

func MustGetAllFilesWithExtInASpecificFolder(dir, ext string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to read in folder "%s": %s`, dir, err))
	}

	var fileList []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ext) {
			fileList = append(fileList, f.Name())
		}
	}

	return fileList
}

// based on https://stackoverflow.com/a/67629473
func MustGetAllFilesWithExtsInASpecificFolderAndSubFolders(dir string, exts ...string) []string {
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
		logger.WriteError(fmt.Sprintf("failed to walk dir \"%s\": %s", dir, err))
	}

	return a
}

func fileHasOneOfExts(fileName string, exts []string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}

	return false
}

func MustRename(src, dest string) {
	err := os.Rename(src, dest)

	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to rename \"%s\" to \"%s\": %s", src, dest, err))
	}
}

func MustGetFileSize(path string) float64 {
	if strings.TrimSpace(path) == "" {
		logger.WriteError("to get a file's size it must have a non-empty path")
	}

	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.WriteError(fmt.Sprintf(`"%s" does not exist so the the file size cannot be retrieved`, path))
		}

		logger.WriteError(fmt.Sprintf(`could not verify that "%s" exists to check its size: %s`, path, err))
	}

	return float64(f.Size()) / bytesInAKiloByte
}

func MustCreateFolderIfNotExists(path string) {
	if FolderExists(path) {
		return
	}

	err := os.MkdirAll(path, folderPerms)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to create folder(s) for path \"%s\": %s", path, err))
	}
}
