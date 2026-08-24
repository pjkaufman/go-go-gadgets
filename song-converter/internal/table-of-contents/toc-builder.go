package tableofcontents

import "github.com/pjkaufman/go-go-gadgets/song-converter/internal/converter"

type TocBuilder = func(headerIds []string, headerInfo []converter.MdFileInfo) (string, error)
