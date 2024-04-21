package filehandler

import (
	"fmt"
	"os"
	"strings"

	archiver "github.com/mholt/archiver/v3"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func ConvertRarToCbz(src string) {
	var err error

	var srcDir = GetFileFolder(src)
	var tempFolder = JoinPath(srcDir, "cbz")
	if FolderExists(tempFolder) {
		err = os.RemoveAll(tempFolder)

		if err != nil {
			logger.WriteError(fmt.Sprintf("failed to delete the destination directory \"%s\": %s", tempFolder, err))
		}
	}

	rar := archiver.NewRar()
	err = rar.Unarchive(src, "cbz")
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to unarchive "%s": %s`, src, err))
	}

	var dest = strings.Replace(src, ".cbr", ".cbz", 1)
	err = Rezip(tempFolder, dest)
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to zip "%s" to "%s": %s`, tempFolder, dest, err))
	}

	err = os.RemoveAll(tempFolder)
	if err != nil {
		logger.WriteError(fmt.Sprintf("failed to delete the destination directory \"%s\": %s", tempFolder, err))
	}
}
