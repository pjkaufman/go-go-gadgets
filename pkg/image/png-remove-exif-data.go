package image

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image/png"

	pngstructure "github.com/dsoprea/go-png-image-structure"
)

var ErrNotPng = errors.New("the provided bytes do not look to pertain to a png file")

// Based on https://github.com/scottleedavis/go-exif-remove/blob/7e059d59340538e639ab516ea037dec825d5b662/exif_remove.go
func PngRemoveExifData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	pmp := pngstructure.NewPngMediaParser()

	if !pmp.LooksLikeFormat(data) {
		return nil, ErrNotPng
	}

	segmentInfo, err := pmp.ParseBytes(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse out png bytes: %w", err)
	}

	_, exifData, err := segmentInfo.Exif()
	if err != nil {
		// for some reasons errors.Is(err, pngstructure.ErrNoExif) did not work here so we will compare the error text instead
		if err.Error() == pngstructure.ErrNoExif.Error() {
			return data, nil
		}

		return nil, fmt.Errorf("failed to get png exif data: %w", err)
	}

	if len(exifData) == 0 {
		return data, nil
	}

	// we replace the exif data with empty bytes because otherwise it corrupts the file
	newData := bytes.Replace(data, exifData, make([]byte, len(exifData)), 1)

	chunks := readPNGChunks(bytes.NewReader(newData))

	for _, chunk := range chunks {
		if !chunk.CRCIsValid() {
			offset := int(chunk.Offset) + 8 + int(chunk.Length)
			crc := chunk.CalculateCRC()

			buf := new(bytes.Buffer)
			binary.Write(buf, binary.BigEndian, crc)
			crcBytes := buf.Bytes()

			copy(newData[offset:], crcBytes)
		}
	}

	chunks = readPNGChunks(bytes.NewReader(newData))
	for _, chunk := range chunks {
		if !chunk.CRCIsValid() {
			return nil, errors.New("EXIF removal failed CRC")
		}
	}

	_, err = png.Decode(bytes.NewReader(newData))
	if err != nil {
		return nil, fmt.Errorf("EXIF removal corrupted the data: %w", err)
	}

	return newData, nil
}
