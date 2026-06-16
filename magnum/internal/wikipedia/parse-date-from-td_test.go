//go:build unit

package wikipedia_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/wikipedia"
	"github.com/stretchr/testify/assert"
)

type parseDateFromTdTestCase struct {
	inputTd      string
	expectedDate string
}

var parseDateFromTdTestCases = map[string]parseDateFromTdTestCase{
	"a table data value with a print and digital date only gets the digital date string": {
		inputTd:      `<td>July 5, 2022 (print)<sup id="cite_ref-9" class="reference"><a href="#cite_note-9">[9]</a></sup><br>May 26, 2022 (digital)</td>`,
		expectedDate: "May 26, 2022",
	},
	"a table data value with a date with no indication of the date type gets that date string back": {
		inputTd:      `<td>July 5, 2022 <sup id="cite_ref-9" class="reference"><a href="#cite_note-9">[9]</a></sup></td>`,
		expectedDate: "July 5, 2022",
	},
	"a table data value with just a print date gets an empty string back": {
		inputTd:      `<td>July 9, 2024 (print)<sup id="cite_ref-23" class="reference"><a href="#cite_note-23">[23]</a></sup></td>`,
		expectedDate: "",
	},
	"a table data value with just a physical date gets an empty string back": {
		inputTd:      `<td>July 9, 2024 (physical)<sup id="cite_ref-23" class="reference"><a href="#cite_note-23">[23]</a></sup></td>`,
		expectedDate: "",
	},
	"a table data value which is just a dash gets an empty string back": {
		inputTd:      `<td>—</td>`,
		expectedDate: "",
	},
	"a table data value which is just 'TBA' gets an empty string back": {
		inputTd:      `<td>TBA</td>`,
		expectedDate: "",
	},
}

func TestParseDateFromTd(t *testing.T) {
	t.Parallel()

	for name, args := range parseDateFromTdTestCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			actualDate := wikipedia.ParseDateFromTd(args.inputTd)

			assert.Equal(t, args.expectedDate, actualDate)
		})
	}
}
