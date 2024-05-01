package jpeg

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"

	"github.com/dsoprea/go-exif/v2"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
)

var ErrNotJpeg = errors.New("the provided bytes do not look to pertain to a jpeg file")

// Based on https://github.com/scottleedavis/go-exif-remove/blob/7e059d59340538e639ab516ea037dec825d5b662/exif_remove.go
func JpegRemoveExifData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	jmp := jpegstructure.NewJpegMediaParser()

	if !jmp.LooksLikeFormat(data) {
		return nil, ErrNotJpeg
	}

	segmentInfo, err := jmp.ParseBytes(data)
	if err != nil && segmentInfo == nil {
		return nil, fmt.Errorf("failed to parse out jpeg bytes: %w", err)
	}

	_, exifData, err := segmentInfo.Exif()
	if err != nil {
		// for some reasons errors.Is(err, exif.ErrNoExif) did not work here so we will compare the error text instead
		if err.Error() == exif.ErrNoExif.Error() {
			return data, nil
		}

		return nil, fmt.Errorf("failed to get jpeg exif data: %w", err)
	}

	if len(exifData) == 0 {
		return data, nil
	}

	// we replace the exif data with empty bytes because otherwise it corrupts the file
	newData := bytes.Replace(data, exifData, make([]byte, len(exifData)), 1)

	_, err = jpeg.Decode(bytes.NewReader(newData))
	if err != nil {
		return nil, fmt.Errorf("EXIF removal corrupted the data: %w", err)
	}

	return newData, nil
}
