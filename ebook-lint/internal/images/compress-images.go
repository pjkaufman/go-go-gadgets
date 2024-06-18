package images

import (
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
)

func CompressRelativeImages(baseFolder string, images map[string]struct{}) error {
	for imagePath := range images {
		err := CompressImage(filehandler.JoinPath(baseFolder, imagePath))

		if err != nil {
			return err
		}
	}

	return nil
}
