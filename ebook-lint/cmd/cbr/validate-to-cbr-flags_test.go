//go:build unit

package cbr_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/cbr"
	"github.com/stretchr/testify/assert"
)

type ValidateToCbrFlagsTestCase struct {
	InputDir      string
	ExpectedError string
}

var ValidateToCbrFlagsTestCases = map[string]ValidateToCbrFlagsTestCase{
	"make sure that an empty directory causes a validation error": {
		InputDir:      "	",
		ExpectedError: cbr.DirArgEmpty,
	},
	"make sure that a non-whitespace directory passes validation": {
		InputDir:      "folder",
		ExpectedError: "",
	},
}

func TestValidateToCbrFlags(t *testing.T) {
	for name, args := range ValidateToCbrFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := cbr.ValidateToCbrFlags(args.InputDir)

			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
