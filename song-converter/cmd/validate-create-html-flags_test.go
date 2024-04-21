//go:build unit

package cmd_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/song-converter/cmd"
	"github.com/stretchr/testify/assert"
)

type ValidateCreateHtmlFlagsTestCase struct {
	InputCoverPath  string
	InputStagingDir string
	ExpectedError   string
}

var ValidateCreateHtmlFlagsTestCases = map[string]ValidateCreateHtmlFlagsTestCase{
	"make sure that an empty staging dir causes a validation error": {
		InputStagingDir: "",
		ExpectedError:   cmd.StagingDirArgEmpty,
	},

	"make sure that an empty cover path causes a validation error": {
		InputStagingDir: "value",
		InputCoverPath:  "",
		ExpectedError:   cmd.CoverPathArgEmpty,
	},
	"make sure that an non-md styles path causes a validation error": {
		InputStagingDir: "dir",
		InputCoverPath:  "cover.txt",
		ExpectedError:   cmd.CoverPathNotMdFile,
	},
	"make sure that the cover path that is an md file and a non-empty staging dir passes validation": {
		InputStagingDir: "dir",
		InputCoverPath:  "cover.md",
		ExpectedError:   "",
	},
}

func TestValidateCreateHtmlFlags(t *testing.T) {
	for name, args := range ValidateCreateHtmlFlagsTestCases {
		t.Run(name, func(t *testing.T) {
			err := cmd.ValidateCreateHtmlFlags(args.InputStagingDir, args.InputCoverPath)
			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
