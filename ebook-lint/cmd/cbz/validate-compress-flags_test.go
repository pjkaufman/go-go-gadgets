//go:build unit

package cbz_test

import (
	"errors"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/cmd/cbz"
	"github.com/stretchr/testify/assert"
)

type validateCompressFlagsTestCase struct {
	inputDir    string
	expectedErr error
}

var validateCompressFlagsTestCases = map[string]validateCompressFlagsTestCase{
	"make sure that an empty directory causes a validation error": {
		inputDir:    "	",
		expectedErr: cbz.ErrDirArgEmpty,
	},
	"make sure that a non-whitespace directory passes validation": {
		inputDir: "folder",
	},
}

func TestValidateCompressFlags(t *testing.T) {
	for name, args := range validateCompressFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := cbz.ValidateCompressFlags(args.inputDir)

			if err != nil {
				assert.True(t, errors.Is(err, args.expectedErr))
			} else {
				assert.Nil(t, args.expectedErr)
			}
		})
	}
}
