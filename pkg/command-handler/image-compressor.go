package commandhandler

import (
	"fmt"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

const imgComperssionProgramName = "imgp"
const minimumSizeOfImageToCompress = 150 // 150kb

var compressionParams = []string{"-x", "800x800", "-e", "-O", "-q", "40", "-m", "-w"}
var CompressableImageExts = []string{"png", "jpg", "jpeg"}

func CompressRelativeImages(baseFolder string, images map[string]struct{}) {
	for imagePath := range images {
		CompressImage(filehandler.JoinPath(baseFolder, imagePath))
	}
}

func CompressImage(imagePath string) {
	if !isCompressableImage(imagePath) || filehandler.MustGetFileSize(imagePath) <= minimumSizeOfImageToCompress {
		return
	}

	var params = append(compressionParams, imagePath)
	var errorMsg = fmt.Sprintf(`failed to compress "%s"`, imagePath)
	output := MustGetCommandOutput(imgComperssionProgramName, errorMsg, params...)
	// imgp does not have an error exist status for when an image does not exist so check for that
	// in the output
	if strings.Contains(output, "does not exist") {
		logger.WriteError(fmt.Sprintf(`%s: %s`, errorMsg, output))
	}

}

func isCompressableImage(imagePath string) bool {
	for _, ext := range CompressableImageExts {
		if strings.HasSuffix(strings.ToLower(imagePath), ext) {
			return true
		}
	}

	return false
}
