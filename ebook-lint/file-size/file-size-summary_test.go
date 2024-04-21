//go:build unit

package filesize_test

import (
	"fmt"
	"testing"

	filesize "github.com/pjkaufman/go-go-gadgets/ebook-lint/file-size"
	"github.com/stretchr/testify/assert"
)

type FileSizeSummaryTestCase struct {
	InputFile           string
	InputBeforeKbSize   float64
	InputAfterKbSize    float64
	ExpectedBeforeSize  string
	ExpectedAfterString string
}

var FileSizeSummaryTestCases = map[string]FileSizeSummaryTestCase{
	"make sure that kilobytes are left as is when they do not exceed 1,024": {
		InputFile:           "test.cbz",
		InputBeforeKbSize:   100,
		InputAfterKbSize:    50,
		ExpectedBeforeSize:  "100.00 KB",
		ExpectedAfterString: "50.00 KB",
	},
	"make sure that kilobytes are truncated when they have more than 2 decimal places": {
		InputFile:           "test.cbz",
		InputBeforeKbSize:   100.5678,
		InputAfterKbSize:    50.878567,
		ExpectedBeforeSize:  "100.57 KB",
		ExpectedAfterString: "50.88 KB",
	},
	"make sure that kilobytes are converted to megabytes when there are more than 1024 of them": {
		InputFile:           "test.cbz",
		InputBeforeKbSize:   1025,
		InputAfterKbSize:    50.878567,
		ExpectedBeforeSize:  "1.00 MB",
		ExpectedAfterString: "50.88 KB",
	},
	"make sure that kilobytes are converted to gigabytes when there are more than 1000000 of them": {
		InputFile:           "test.cbz",
		InputBeforeKbSize:   2000000,
		InputAfterKbSize:    50.878567,
		ExpectedBeforeSize:  "2.00 GB",
		ExpectedAfterString: "50.88 KB",
	},
}

func TestFileSizeSummary(t *testing.T) {
	for name, args := range FileSizeSummaryTestCases {
		t.Run(name, func(t *testing.T) {
			var originalFile = args.InputFile + ".original"
			var expected = fmt.Sprintf(filesize.FileSummaryTemplate, filesize.CliLineSeparator, originalFile, args.ExpectedBeforeSize, args.InputFile, args.ExpectedAfterString)
			actual := filesize.FileSizeSummary(originalFile, args.InputFile, args.InputBeforeKbSize, args.InputAfterKbSize)

			assert.Equal(t, expected, actual)
		})
	}
}
