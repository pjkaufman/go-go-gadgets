//go:build unit

package cmd_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/cmd"
	"github.com/stretchr/testify/assert"
)

type ValidateCreateCsvFlagsTestCase struct {
	InputStagingDir string
	ExpectedError   string
}

// errors that get handled as errors are represented as panics
var ValidateCreateCsvFlagsTestCases = map[string]ValidateCreateCsvFlagsTestCase{
	"make sure that an empty working dir causes a validation error": {
		InputStagingDir: "",
		ExpectedError:   cmd.StagingDirArgEmpty,
	},
	"make sure that working dir with a value passes validation": {
		InputStagingDir: "folder",
	},
}

func TestValidateCreateCsvFlags(t *testing.T) {
	for name, args := range ValidateCreateCsvFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := cmd.ValidateCreateCsvFlags(args.InputStagingDir)
			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
