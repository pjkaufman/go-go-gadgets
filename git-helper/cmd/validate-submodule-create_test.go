//go:build unit

package cmd_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/git-helper/cmd"
	"github.com/stretchr/testify/assert"
)

type ValidateCreateSubmoduleTestCase struct {
	InputTicket           string
	InputBranchName       string
	InputRepoParentPath   string
	InputSubmoduleName    string
	InputBranchPrefixName string
	ExpectedError         string
}

var ValidateCreateSubmoduleTestCases = map[string]ValidateCreateSubmoduleTestCase{
	"make sure that an empty ticket name causes a validation error": {
		InputTicket:           "",
		InputBranchName:       "name",
		InputRepoParentPath:   "/users/username/home/",
		InputSubmoduleName:    "submodule",
		InputBranchPrefixName: "prefix",
		ExpectedError:         cmd.TicketArgEmpty,
	},
	"make sure that an empty branch name causes a validation error": {
		InputTicket:           "ticket",
		InputBranchName:       "",
		InputRepoParentPath:   "/users/username/home/",
		InputSubmoduleName:    "submodule",
		InputBranchPrefixName: "prefix",
		ExpectedError:         cmd.BranchNameArgEmpty,
	},
	"make sure that an empty repo parent directory causes a validation error": {
		InputTicket:           "ticket",
		InputBranchName:       "name",
		InputRepoParentPath:   "",
		InputSubmoduleName:    "submodule",
		InputBranchPrefixName: "prefix",
		ExpectedError:         cmd.RepoParentPathArgEmpty,
	},
	"make sure that an empty submodule name causes a validation error": {
		InputTicket:           "ticket",
		InputBranchName:       "name",
		InputRepoParentPath:   "/users/username/home/",
		InputSubmoduleName:    "",
		InputBranchPrefixName: "prefix",
		ExpectedError:         cmd.SubmoduleNameArgEmpty,
	},
	"make sure that an empty branch prefix name causes a validation error": {
		InputTicket:           "ticket",
		InputBranchName:       "name",
		InputRepoParentPath:   "/users/username/home/",
		InputSubmoduleName:    "submodule",
		InputBranchPrefixName: "",
		ExpectedError:         cmd.BranchPrefixArgEmpty,
	},
	"make sure that ticket, branch name, branch prefix, and repo parent directory having values passes validation": {
		InputTicket:           "ticket",
		InputBranchName:       "name",
		InputRepoParentPath:   "/users/username/home/",
		InputSubmoduleName:    "submodule",
		InputBranchPrefixName: "prefix",
		ExpectedError:         "",
	},
}

func TestValidateCreateSubmodule(t *testing.T) {
	for name, args := range ValidateCreateSubmoduleTestCases {
		t.Run(name, func(t *testing.T) {
			err := cmd.ValidateSubmoduleCreate(args.InputTicket, args.InputBranchName, args.InputRepoParentPath, args.InputSubmoduleName, args.InputBranchPrefixName)

			if err != nil {
				assert.Equal(t, args.ExpectedError, err.Error())
			} else {
				assert.Equal(t, args.ExpectedError, "")
			}
		})
	}
}
