//go:build unit

package cmdtomd_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	cmdtomd "github.com/pjkaufman/go-go-gadgets/pkg/cmd-to-md"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type FlagsToMdTestCase struct {
	CreateFlags    func() *pflag.FlagSet
	ExpectedOutput string
}

var FlagsToMdTestCases = map[string]FlagsToMdTestCase{
	"make sure that a nil command results in an empty string being generated": {
		CreateFlags: func() *pflag.FlagSet {
			return nil
		},
		ExpectedOutput: "",
	},
	"make sure that flags are properly displayed in the generated table": {
		CreateFlags: func() *pflag.FlagSet {

			flagSet := pflag.NewFlagSet("flags", pflag.ContinueOnError)

			flagSet.BoolP("boolean", "b", false, "A boolean flag")
			flagSet.IntP("integer", "i", 0, "An integer flag")
			flagSet.StringP("string", "s", "default", "A string flag")

			return flagSet
		},
		ExpectedOutput: fmt.Sprintf("%s\n%s\n%s", cmdtomd.TableHeader, cmdtomd.Separator, "| b | boolean | A boolean flag |  | false | false |  |\n| i | integer | An integer flag | int | 0 | false |  |\n| s | string | A string flag | string | default | false |  |"),
	},
	"make sure that required flags are properly listed as required in the generated table": {
		CreateFlags: func() *pflag.FlagSet {

			flagSet := pflag.NewFlagSet("flags", pflag.ContinueOnError)

			flagSet.BoolP("required boolean", "b", false, "A required boolean flag")
			err := flagSet.SetAnnotation("required boolean", cobra.BashCompOneRequiredFlag, []string{"true"})
			if err != nil {
				log.Fatalf("failed to add required flag %q: %s", "required boolean", err)
			}

			return flagSet
		},
		ExpectedOutput: fmt.Sprintf("%s\n%s\n%s", cmdtomd.TableHeader, cmdtomd.Separator, "| b | required boolean | A required boolean flag |  | false | true |  |"),
	},
	"make sure that specific file extension flags have their other notes properly list their possible extensions in the generated table": {
		CreateFlags: func() *pflag.FlagSet {

			flagSet := pflag.NewFlagSet("flags", pflag.ContinueOnError)

			flagSet.StringP("cover-file", "c", "", "The cover file to use")
			err := flagSet.SetAnnotation("cover-file", cobra.BashCompFilenameExt, []string{"md"})
			if err != nil {
				log.Fatalf("failed to add filename extension flag %q: %s", "cover-file", err)
			}

			flagSet.StringP("image", "i", "", "The image to use")
			err = flagSet.SetAnnotation("image", cobra.BashCompFilenameExt, []string{"jpg", "png", "jpeg"})
			if err != nil {
				log.Fatalf("failed to add filename extension flag %q: %s", "image", err)
			}

			return flagSet
		},
		ExpectedOutput: fmt.Sprintf("%s\n%s\n%s", cmdtomd.TableHeader, cmdtomd.Separator, "| c | cover-file | The cover file to use | string |  | false | Should be a file with one of the following extensions: md |\n| i | image | The image to use | string |  | false | Should be a file with one of the following extensions: jpg, png, jpeg |"),
	},
	"make sure that directory flags have their other notes properly list that they are for directories in the generated table": {
		CreateFlags: func() *pflag.FlagSet {

			flagSet := pflag.NewFlagSet("flags", pflag.ContinueOnError)

			flagSet.StringP("working-dir", "w", ".", "The directory to do operations in")
			err := flagSet.SetAnnotation("working-dir", cobra.BashCompSubdirsInDir, []string{})
			if err != nil {
				log.Fatalf("failed to add directory flag %q: %s", "working-dir", err)
			}

			return flagSet
		},
		ExpectedOutput: fmt.Sprintf("%s\n%s\n%s", cmdtomd.TableHeader, cmdtomd.Separator, "| w | working-dir | The directory to do operations in | string | . | false | Should be a directory |"),
	},
	"make sure that hidden flags are not included in the generated table": {
		CreateFlags: func() *pflag.FlagSet {

			flagSet := pflag.NewFlagSet("flags", pflag.ContinueOnError)

			flagSet.BoolP("required boolean", "b", false, "A required boolean flag")
			err := flagSet.SetAnnotation("required boolean", cobra.BashCompOneRequiredFlag, []string{"true"})
			if err != nil {
				log.Fatalf("failed to add required flag %q: %s", "required boolean", err)
			}

			flagSet.BoolP("help", "h", false, "Show the help")
			err = flagSet.MarkHidden("help")
			if err != nil {
				log.Fatalf("failed to add hidden flag %q: %s", "help", err)
			}

			return flagSet
		},
		ExpectedOutput: fmt.Sprintf("%s\n%s\n%s", cmdtomd.TableHeader, cmdtomd.Separator, "| b | required boolean | A required boolean flag |  | false | true |  |"),
	},
}

func TestFlagsToMd(t *testing.T) {
	for name, args := range FlagsToMdTestCases {
		t.Run(name, func(t *testing.T) {
			var actual strings.Builder
			cmdtomd.FlagsToMd(args.CreateFlags(), &actual)

			assert.Equal(t, args.ExpectedOutput, actual.String())
		})
	}
}
