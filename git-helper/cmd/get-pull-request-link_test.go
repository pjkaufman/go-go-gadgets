//go:build unit

package cmd_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/git-helper/cmd"
	"github.com/stretchr/testify/assert"
)

type GetPullRequestLinkTestCase struct {
	Input  string
	Output string
}

var GetPullRequestLinkTestCases = map[string]GetPullRequestLinkTestCase{
	"make sure that an empty input results in a blank string": {
		Input:  "",
		Output: "",
	},
	"make sure that an empty branch name causes a validation error": {
		Input: `Enumerating objects: 11, done.
Counting objects: 100% (11/11), done.
Delta compression using up to 12 threads
Compressing objects: 100% (6/6), done.
Writing objects: 100% (6/6), 536 bytes | 0 bytes/s, done.
Total 6 (delta 5), reused 0 (delta 0), pack-reused 0
remote: Resolving deltas: 100% (5/5), completed with 5 local objects.
remote:
remote: Create a pull request for 'branch-name' on GitHub by visiting:
remote:      https://github.com/pjkaufman/dotfiles/pull/new/branch-name
remote:
To github.com:pjkaufman/dotfiles.git
	* [new branch]        branch-name -> branch-name
branch 'branch-name' set up to track 'origin/branch-name'.`,
		Output: "https://github.com/pjkaufman/dotfiles/pull/new/branch-name",
	},
}

func TestGetPullRequestLink(t *testing.T) {
	for name, args := range GetPullRequestLinkTestCases {
		t.Run(name, func(t *testing.T) {
			prLink := cmd.GetPullRequestLink(args.Input)

			assert.Equal(t, args.Output, prLink)
		})
	}
}
