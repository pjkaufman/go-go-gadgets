package images

import (
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
)

func CompressRelativeImages(baseFolder string, images map[string]struct{}) {
	for imagePath := range images {
		CompressImage(filehandler.JoinPath(baseFolder, imagePath))
	}
}
