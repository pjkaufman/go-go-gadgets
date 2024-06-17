package filehandler

import (
	"fmt"
	"os"
	"strings"

	archiver "github.com/mholt/archiver/v3"
)

func ConvertRarToCbz(src string) error {
	var err error

	var srcDir = GetFileFolder(src)
	var tempFolder = JoinPath(srcDir, "cbz")
	fileExists, err := FolderExists(tempFolder)
	if err != nil {
		return err
	}

	if fileExists {
		err = os.RemoveAll(tempFolder)

		if err != nil {
			return fmt.Errorf("failed to delete the destination directory %q: %w", tempFolder, err)
		}
	}

	rar := archiver.NewRar()
	err = rar.Unarchive(src, "cbz")
	if err != nil {
		return fmt.Errorf(`failed to unarchive %q: %w`, src, err)
	}

	var dest = strings.Replace(src, ".cbr", ".cbz", 1)
	err = Rezip(tempFolder, dest)
	if err != nil {
		return fmt.Errorf(`failed to zip %q to %q: %w`, tempFolder, dest, err)
	}

	err = os.RemoveAll(tempFolder)
	if err != nil {
		return fmt.Errorf("failed to delete the destination directory %q: %w", tempFolder, err)
	}

	return nil
}
