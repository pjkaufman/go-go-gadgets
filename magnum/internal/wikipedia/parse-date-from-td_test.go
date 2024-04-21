//go:build unit

package wikipedia_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/magnum/internal/wikipedia"
	"github.com/stretchr/testify/assert"
)

type ParseDateFromTdTestCase struct {
	InputTd      string
	ExpectedDate string
}

var ParseDateFromTdTestCases = map[string]ParseDateFromTdTestCase{
	"a table data value with a print and digital date only gets the digital date string": {
		InputTd:      `<td>July 5, 2022 (print)<sup id="cite_ref-9" class="reference"><a href="#cite_note-9">[9]</a></sup><br>May 26, 2022 (digital)</td>`,
		ExpectedDate: "May 26, 2022",
	},
	"a table data value with a date with no indication of the date type gets that date string back": {
		InputTd:      `<td>July 5, 2022 <sup id="cite_ref-9" class="reference"><a href="#cite_note-9">[9]</a></sup></td>`,
		ExpectedDate: "July 5, 2022",
	},
	"a table data value with just a print date gets an empty string back": {
		InputTd:      `<td>July 9, 2024 (print)<sup id="cite_ref-23" class="reference"><a href="#cite_note-23">[23]</a></sup></td>`,
		ExpectedDate: "",
	},
	"a table data value with just a physical date gets an empty string back": {
		InputTd:      `<td>July 9, 2024 (physical)<sup id="cite_ref-23" class="reference"><a href="#cite_note-23">[23]</a></sup></td>`,
		ExpectedDate: "",
	},
	"a table data value which is just a dash gets an empty string back": {
		InputTd:      `<td>—</td>`,
		ExpectedDate: "",
	},
	"a table data value which is just 'TBA' gets an empty string back": {
		InputTd:      `<td>TBA</td>`,
		ExpectedDate: "",
	},
}

func TestParseDateFromTd(t *testing.T) {
	for name, args := range ParseDateFromTdTestCases {
		t.Run(name, func(t *testing.T) {
			actualDate := wikipedia.ParseDateFromTd(args.InputTd)

			assert.Equal(t, args.ExpectedDate, actualDate)
		})
	}
}
