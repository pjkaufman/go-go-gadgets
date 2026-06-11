//go:build unit

package markdown_test

import (
	"strings"
	"testing"

	markdown "github.com/pjkaufman/go-go-gadgets/pkg/cli/markdown"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type CommandToMdTestCase struct {
	Command        *cobra.Command
	ExpectedOutput string
}

var CommandToMdTestCases = map[string]CommandToMdTestCase{
	"make sure that a nil command results in an empty string being generated": {
		Command:        nil,
		ExpectedOutput: "",
	},
	"make sure that a command with no long description uses the short description": {
		Command: &cobra.Command{
			Use:   "test",
			Short: "Short description",
		},
		ExpectedOutput: "### test\n\nShort description\n\n",
	},
	"make sure that a command with a long description uses the long description": {
		Command: &cobra.Command{
			Use:   "test",
			Short: "Short description",
			Long:  "Long description",
		},
		ExpectedOutput: "### test\n\nLong description\n\n",
	},
	"make sure that a command with an example has it included in the result": {
		Command: &cobra.Command{
			Use:   "test",
			Short: "Short description",
			Long:  "Long description",
			Example: `To run the test command:
cmd to md test --flag`,
		},
		ExpectedOutput: "### test\n\nLong description\n\n#### Usage\n\n``` bash\n# To run the test command:\ncmd to md test --flag\n```\n\n",
	},
}

func TestCommandToMd(t *testing.T) {
	parent := cobra.Command{}
	for name, args := range CommandToMdTestCases {
		t.Run(name, func(t *testing.T) {
			if args.Command != nil {
				parent.AddCommand(args.Command) // makes sure the command is not registered as a root command

				args.Command.Run = func(cmd *cobra.Command, args []string) {} // allows the command to be considered runnable
			}

			var actual strings.Builder
			markdown.CommandToMd(args.Command, &actual, 3)

			assert.Equal(t, args.ExpectedOutput, actual.String())
		})
	}
}
