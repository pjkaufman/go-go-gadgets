//go:build unit

package cbz_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/cbz"
	"github.com/stretchr/testify/assert"
)

type ValidateCompressFlagsTestCase struct {
	InputDir      string
	ExpectedError string
}

var ValidateCompressFlagsTestCases = map[string]ValidateCompressFlagsTestCase{
	"make sure that an empty directory causes a validation error": {
		InputDir:      "	",
		ExpectedError: cbz.DirArgEmpty,
	},
	"make sure that a non-whitespace directory passes validation": {
		InputDir:      "folder",
		ExpectedError: "",
	},
}

func TestValidateCompressFlags(t *testing.T) {
	for name, args := range ValidateCompressFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := cbz.ValidateCompressFlags(args.InputDir)

			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
