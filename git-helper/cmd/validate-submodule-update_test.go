//go:build unit

package cmd_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/git-helper/cmd"
	"github.com/stretchr/testify/assert"
)

type ValidateSubmoduleUpdateTestCase struct {
	InputBranchName     string
	InputRepoParentPath string
	InputSubmoduleName  string
	ExpectedError       string
}

var ValidateSubmoduleUpdateTestCases = map[string]ValidateSubmoduleUpdateTestCase{
	"make sure that an empty branch name causes a validation error": {
		InputBranchName:     "",
		InputRepoParentPath: "/users/username/home/",
		InputSubmoduleName:  "submodule",
		ExpectedError:       cmd.BranchNameArgEmpty,
	},
	"make sure that an empty repo parent directory causes a validation error": {
		InputBranchName:     "name",
		InputRepoParentPath: "",
		InputSubmoduleName:  "submodule",
		ExpectedError:       cmd.RepoParentPathArgEmpty,
	},
	"make sure that an empty submodule name causes a validation error": {
		InputBranchName:     "name",
		InputRepoParentPath: "/users/username/home/",
		InputSubmoduleName:  "",
		ExpectedError:       cmd.SubmoduleNameArgEmpty,
	},
	"make sure that ticket, branch name, and repo parent directory having values passes validation": {
		InputBranchName:     "name",
		InputRepoParentPath: "/users/username/home/",
		InputSubmoduleName:  "submodule",
		ExpectedError:       "",
	},
}

func TestValidateSubmoduleUpdate(t *testing.T) {
	for name, args := range ValidateSubmoduleUpdateTestCases {
		t.Run(name, func(t *testing.T) {
			err := cmd.ValidateSubmoduleUpdate(args.InputBranchName, args.InputRepoParentPath, args.InputSubmoduleName)

			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
