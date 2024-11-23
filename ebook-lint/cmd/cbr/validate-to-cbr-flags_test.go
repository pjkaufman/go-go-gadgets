//go:build unit

package cbr_test

import (
	"errors"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/cbr"
	"github.com/stretchr/testify/assert"
)

type validateToCbrFlagsTestCase struct {
	inputDir    string
	expectedErr error
}

var ValidateToCbrFlagsTestCases = map[string]validateToCbrFlagsTestCase{
	"make sure that an empty directory causes a validation error": {
		inputDir:    "	",
		expectedErr: cbr.ErrDirArgEmpty,
	},
	"make sure that a non-whitespace directory passes validation": {
		inputDir: "folder",
	},
}

func TestValidateToCbrFlags(t *testing.T) {
	for name, args := range ValidateToCbrFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := cbr.ValidateToCbrFlags(args.inputDir)

			if err != nil {
				assert.True(t, errors.Is(err, args.expectedErr))
			} else {
				assert.Nil(t, args.expectedErr)
			}
		})
	}
}
