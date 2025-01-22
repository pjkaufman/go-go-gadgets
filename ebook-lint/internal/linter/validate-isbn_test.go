//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type isValidIsbnTestCase struct {
	inputVal       string
	expectedOutput bool
}

var isValidIsbnTestCases = map[string]isValidIsbnTestCase{
	"Valid ISBN-10 with ISBN prefix": {
		inputVal:       "ISBN-10: 0-306-40615-2",
		expectedOutput: true,
	},
	"Valid ISBN-10 with hyphens": {
		inputVal:       "0-306-40615-2",
		expectedOutput: true,
	},
	"Valid ISBN-10 without hyphens": {
		inputVal:       "0306406152",
		expectedOutput: true,
	},
	"Valid ISBN-13 with ISBN prefix": {
		inputVal:       "ISBN-13: 978-3-16-148410-0",
		expectedOutput: true,
	},
	"Valid ISBN-13 with hyphens": {
		inputVal:       "978-3-16-148410-0",
		expectedOutput: true,
	},
	"Valid ISBN-13 without hyphens": {
		inputVal:       "9783161484100",
		expectedOutput: true,
	},
	"Valid ISBN-10 with X": {
		inputVal:       "123456789X",
		expectedOutput: true,
	},
	"Invalid ISBN-10": {
		inputVal:       "1234567890",
		expectedOutput: false,
	},
	"Invalid ISBN-13 because check sum does not add up": {
		inputVal:       "9781234567896",
		expectedOutput: false,
	},
	"DOI is not ISBN": {
		inputVal:       "10.1000/182",
		expectedOutput: false,
	},
	"UUID is not ISBN": {
		inputVal:       "123e4567-e89b-12d3-a456-426614174000",
		expectedOutput: false,
	},
	"URN ISBN is not ISBN": {
		inputVal:       "urn:isbn:978-3-16-148410-0",
		expectedOutput: false,
	},
}

func TestIsValidISBN(t *testing.T) {
	for name, args := range isValidIsbnTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.IsValidISBN(args.inputVal)
			assert.Equal(t, args.expectedOutput, actual)
		})
	}
}
