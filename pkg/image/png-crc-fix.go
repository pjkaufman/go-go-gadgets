package image

// from https://github.com/scottleedavis/go-exif-remove/blob/7e059d59340538e639ab516ea037dec825d5b662/png_crc_fix.go

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
)

const chunkStartOffset = 8
const endChunk = "IEND"

type pngChunk struct {
	Offset int64
	Length uint32
	Type   [4]byte
	Data   []byte
	CRC    uint32
}

func (p pngChunk) String() string {
	return fmt.Sprintf("%s@%x - %X - Valid CRC? %v", p.Type, p.Offset, p.CRC, p.CRCIsValid())
}

func (p pngChunk) Bytes() []byte {
	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, p.Type)
	buffer.Write(p.Data)

	return buffer.Bytes()
}

func (p pngChunk) CRCIsValid() bool {
	return p.CRC == p.CalculateCRC()
}

func (p pngChunk) CalculateCRC() uint32 {
	crcTable := crc32.MakeTable(crc32.IEEE)

	return crc32.Checksum(p.Bytes(), crcTable)
}

func (p pngChunk) CRCOffset() int64 {
	return p.Offset + int64(8+p.Length)
}

func readPNGChunks(reader io.ReadSeeker) []pngChunk {
	chunks := []pngChunk{}

	_, err := reader.Seek(chunkStartOffset, io.SeekStart)
	if err != nil {
		return chunks
	}

	readChunk := func() (*pngChunk, error) {
		var chunk pngChunk
		chunk.Offset, _ = reader.Seek(0, io.SeekCurrent)

		err = binary.Read(reader, binary.BigEndian, &chunk.Length)
		if err != nil {
			goto read_error
		}

		chunk.Data = make([]byte, chunk.Length)

		err = binary.Read(reader, binary.BigEndian, &chunk.Type)
		if err != nil {
			goto read_error
		}

		if read, err := reader.Read(chunk.Data); read == 0 || err != nil {
			goto read_error
		}

		err = binary.Read(reader, binary.BigEndian, &chunk.CRC)
		if err != nil {
			goto read_error
		}

		return &chunk, nil

	read_error:
		return nil, fmt.Errorf("read error: %w", err)
	}

	chunk, err := readChunk()
	if err != nil {
		return chunks
	}

	chunks = append(chunks, *chunk)

	// Read the first chunk
	for string(chunks[len(chunks)-1].Type[:]) != endChunk {

		chunk, err := readChunk()
		if err != nil {
			break
		}

		chunks = append(chunks, *chunk)
	}

	return chunks
}
