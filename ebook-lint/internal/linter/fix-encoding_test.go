//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

type fixEncodingTestCase struct {
	input    string
	expected string
}

var fixEncodingTestCases = map[string]fixEncodingTestCase{
	"when the xml tag is missing, an xml tag is added": {
		input: "<html></html>",
		expected: `<?xml version="1.0" encoding="utf-8"?>
<html></html>`,
	},
	"when the xml tag is present, but does not have encoding, encoding should be added": {
		input: `<?xml version="1.0"?>
<html></html>`,
		expected: `<?xml version="1.0" encoding="utf-8"?>
<html></html>`,
	},
	"when the xml tag is present, and does have encoding, encoding should be left as is": {
		input: `<?xml version="1.0" encoding="text"?>
<html></html>`,
		expected: `<?xml version="1.0" encoding="text"?>
<html></html>`,
	},
	"when there are multiple xml tags present, only the 1st one will be modified": {
		input: `<?xml version="1.0"?><?xml version="1.0"?>
<html></html>`,
		expected: `<?xml version="1.0" encoding="utf-8"?><?xml version="1.0"?>
<html></html>`,
	},
}

func TestFixEncoding(t *testing.T) {
	for name, args := range fixEncodingTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.EnsureEncodingIsPresent(args.input)

			assert.Equal(t, args.expected, actual)
		})
	}
}
