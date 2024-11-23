//go:build unit

package filesize_test

import (
	"fmt"
	"testing"

	filesize "github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/file-size"
	"github.com/stretchr/testify/assert"
)

type fileSizeSummaryTestCase struct {
	inputFile           string
	inputBeforeKbSize   float64
	inputAfterKbSize    float64
	expectedBeforeSize  string
	expectedAfterString string
}

var fileSizeSummaryTestCases = map[string]fileSizeSummaryTestCase{
	"make sure that kilobytes are left as is when they do not exceed 1,024": {
		inputFile:           "test.cbz",
		inputBeforeKbSize:   100,
		inputAfterKbSize:    50,
		expectedBeforeSize:  "100.00 KB",
		expectedAfterString: "50.00 KB",
	},
	"make sure that kilobytes are truncated when they have more than 2 decimal places": {
		inputFile:           "test.cbz",
		inputBeforeKbSize:   100.5678,
		inputAfterKbSize:    50.878567,
		expectedBeforeSize:  "100.57 KB",
		expectedAfterString: "50.88 KB",
	},
	"make sure that kilobytes are converted to megabytes when there are more than 1024 of them": {
		inputFile:           "test.cbz",
		inputBeforeKbSize:   1025,
		inputAfterKbSize:    50.878567,
		expectedBeforeSize:  "1.00 MB",
		expectedAfterString: "50.88 KB",
	},
	"make sure that kilobytes are converted to gigabytes when there are more than 1000000 of them": {
		inputFile:           "test.cbz",
		inputBeforeKbSize:   2000000,
		inputAfterKbSize:    50.878567,
		expectedBeforeSize:  "2.00 GB",
		expectedAfterString: "50.88 KB",
	},
}

func TestFileSizeSummary(t *testing.T) {
	for name, args := range fileSizeSummaryTestCases {
		t.Run(name, func(t *testing.T) {
			var originalFile = args.inputFile + ".original"
			var expected = fmt.Sprintf(filesize.FileSummaryTemplate, filesize.CliLineSeparator, originalFile, args.expectedBeforeSize, args.inputFile, args.expectedAfterString)
			actual := filesize.FileSizeSummary(originalFile, args.inputFile, args.inputBeforeKbSize, args.inputAfterKbSize)

			assert.Equal(t, expected, actual)
		})
	}
}
