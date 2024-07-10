//go:build unit

package cmdtomd_test

import (
	"strings"
	"testing"

	cmdtomd "github.com/pjkaufman/go-go-gadgets/pkg/cmd-to-md"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type ParseOpfContentsTestCase struct {
	Command        *cobra.Command
	ExpectedOutput string
}

var ParseOpfContentsTestCases = map[string]ParseOpfContentsTestCase{
	"make sure that a nil command results in an empty string being generated": {
		Command:        nil,
		ExpectedOutput: "",
	},
	"make sure that a command with no long description uses the short description": {
		Command: &cobra.Command{
			Use:   "test",
			Short: "Short description",
			// Long:  "Long description",
			// Example: `To run the test command:
			// cmd tomd test --flag`,
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
	// TODO: add simple and complex test with actual flags
}

func TestCommandToMd(t *testing.T) {
	for name, args := range ParseOpfContentsTestCases {
		t.Run(name, func(t *testing.T) {
			var actual strings.Builder
			cmdtomd.CommandToMd(args.Command, &actual, 3)

			assert.Equal(t, args.ExpectedOutput, actual.String())
		})
	}
}
