package flags_test

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/pjkaufman/go-go-gadgets/pkg/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type validationFlagTestCase struct {
	flag                      flags.Flag
	setup                     func(t *testing.T, flag flags.Flag)
	expectedErrorStringSubset *string
}

// error indicators
var (
	whitespaceOrEmptyIndicator        = createStringPointer("non-whitespace value")
	requiredExtensionsIndicator       = createStringPointer("following extensions")
	fileOrDirectoryExistenceIndicator = createStringPointer("must exist")
	enumOptionIndicator               = createStringPointer("must be one of")
)

var stringFlagValidationTestCases = map[string]validationFlagTestCase{
	"A string flag that is required and nil should return an error": {
		flag:                      flags.NewStringFlag(true, false, nil, "test flag", "f", "", ""),
		expectedErrorStringSubset: whitespaceOrEmptyIndicator,
	},
	"A string flag that is required and is just whitespace should return an error": {
		flag:                      flags.NewStringFlag(true, false, createStringPointer("    "), "test flag", "f", "", ""),
		expectedErrorStringSubset: whitespaceOrEmptyIndicator,
	},
	"A string flag that is required and is not just whitespace should not return an error": {
		flag: flags.NewStringFlag(true, false, createStringPointer(" test value   "), "test flag", "f", "", ""),
	},
	"A string flag that is not required and is just whitespace should not return an error": {
		flag: flags.NewStringFlag(false, false, createStringPointer("    "), "test flag", "f", "", ""),
	},
	"A string flag that is not required and is nil should not return an error": {
		flag: flags.NewStringFlag(false, false, nil, "test flag", "f", "", ""),
	},
	"A string flag that is not required and is not nil or whitespace should not return an error": {
		flag: flags.NewStringFlag(false, false, createStringPointer("test"), "test flag", "f", "", ""),
	},
}

var fileFlagValidationTestCases = map[string]validationFlagTestCase{
	"A file flag that is required and nil should return an error": {
		flag:                      flags.NewFileFlag(true, false, nil, "test flag", "f", "", "", nil, false),
		expectedErrorStringSubset: whitespaceOrEmptyIndicator,
	},
	"A file flag that is required and is just whitespace should return an error": {
		flag:                      flags.NewFileFlag(true, false, createStringPointer("    "), "test flag", "f", "", "", nil, false),
		expectedErrorStringSubset: whitespaceOrEmptyIndicator,
	},
	"A file flag that is required, is not whitespace, has a value, and does not have a required extension should not return an error": {
		flag: flags.NewFileFlag(true, false, createStringPointer("file.txt"), "test flag", "f", "", "", nil, false),
	},
	"A file flag that is required, is not whitespace, has a value, and does not have a required extension should return an error": {
		flag:                      flags.NewFileFlag(true, false, createStringPointer("file.json"), "test flag", "f", "", "", []string{"txt"}, false),
		expectedErrorStringSubset: requiredExtensionsIndicator,
	},
	"A file flag that is required, is not whitespace, has a value, has a required extension, and does not exist but must should return an error": {
		flag:                      flags.NewFileFlag(true, false, createStringPointer("file.txt"), "test flag", "f", "", "", []string{"txt"}, true),
		expectedErrorStringSubset: fileOrDirectoryExistenceIndicator,
	},
	"A file flag that is required, is not whitespace, has a value, has a required extension, and exists and it must should not return an error": {
		flag:  flags.NewFileFlag(true, false, createStringPointer("file.txt"), "test flag", "f", "", "", []string{"txt"}, true),
		setup: setupFiletTest,
	},
	"A file flag that is not required and is nil should not return an error": {
		flag: flags.NewFileFlag(false, false, nil, "test flag", "f", "", "", nil, false),
	},
	"A file flag that is not required and is whitespace should not return an error": {
		flag: flags.NewFileFlag(false, false, createStringPointer("  "), "test flag", "f", "", "", nil, false),
	},
	"A file flag that is not required, is not whitespace, and has a value should not return an error": {
		flag: flags.NewFileFlag(false, false, createStringPointer("file.txt"), "test flag", "f", "", "", nil, false),
	},
	"A file flag that is not required, is not whitespace, has a value, and has an extension in the required list should not return an error": {
		flag: flags.NewFileFlag(false, false, createStringPointer("file.txt"), "test flag", "f", "", "", []string{"txt"}, false),
	},
	"A file flag that is not required, is not whitespace, has a value, has an extension in the required list, and exists and it must should not return an error": {
		flag:  flags.NewFileFlag(false, false, createStringPointer("new.txt"), "test flag", "f", "", "", []string{"txt"}, true),
		setup: setupFiletTest,
	},
	"A file flag that is not required, is not whitespace, has a value, has an extension in the required list, and does not exist and it must should return an error": {
		flag:                      flags.NewFileFlag(false, false, createStringPointer("new.txt"), "test flag", "f", "", "", []string{"txt"}, true),
		expectedErrorStringSubset: fileOrDirectoryExistenceIndicator,
	},
}

var directoryFlagValidationTestCases = map[string]validationFlagTestCase{
	"A directory flag that is required and nil should return an error": {
		flag:                      flags.NewDirectoryFlag(true, false, nil, "test flag", "f", "", ""),
		expectedErrorStringSubset: whitespaceOrEmptyIndicator,
	},
	"A directory flag that is required and is just whitespace should return an error": {
		flag:                      flags.NewDirectoryFlag(true, false, createStringPointer("    "), "test flag", "f", "", ""),
		expectedErrorStringSubset: whitespaceOrEmptyIndicator,
	},
	"A directory flag that is required, is not whitespace, has a value, and does not exist should return an error": {
		flag:                      flags.NewDirectoryFlag(true, false, createStringPointer("directory"), "test flag", "f", "", ""),
		expectedErrorStringSubset: fileOrDirectoryExistenceIndicator,
	},
	"A directory flag that is required, is not whitespace, has a value, and does exist should not return an error": {
		flag:  flags.NewDirectoryFlag(true, false, createStringPointer("directory"), "test flag", "f", "", ""),
		setup: setupDirectoryTest,
	},
	"A directory flag that is not required and is nil should not return an error": {
		flag: flags.NewDirectoryFlag(false, false, nil, "test flag", "f", "", ""),
	},
	"A directory flag that is not required and is whitespace should not return an error": {
		flag: flags.NewDirectoryFlag(false, false, createStringPointer("  "), "test flag", "f", "", ""),
	},
	"A directory flag that is not required, is not whitespace, has a value, and exists should not return an error": {
		flag:  flags.NewDirectoryFlag(false, false, createStringPointer("directory"), "test flag", "f", "", ""),
		setup: setupDirectoryTest,
	},
	"A directory flag that is not required, is not whitespace, has a value, and does not exist should return an error": {
		flag:                      flags.NewDirectoryFlag(false, false, createStringPointer("directory"), "test flag", "f", "", ""),
		expectedErrorStringSubset: fileOrDirectoryExistenceIndicator,
	},
}

var enumFlagValidationTestCases = map[string]validationFlagTestCase{
	"An enum flag that is required and nil should return an error": {
		flag:                      flags.NewEnumFlag(true, false, nil, "test flag", "f", "", "", []string{"option1", "option2"}),
		expectedErrorStringSubset: whitespaceOrEmptyIndicator,
	},
	"An enum flag that is required and is just whitespace should return an error": {
		flag:                      flags.NewEnumFlag(true, false, createStringPointer("    "), "test flag", "f", "", "", []string{"option1", "option2"}),
		expectedErrorStringSubset: whitespaceOrEmptyIndicator,
	},
	"An enum flag that is required, is not just whitespace, and is one of the allowed options should not return an error": {
		flag: flags.NewEnumFlag(true, false, createStringPointer("option1"), "test flag", "f", "", "", []string{"option1", "option2"}),
	},
	"An enum flag that is required, is not just whitespace, and is not one of the allowed options should return an error": {
		flag:                      flags.NewEnumFlag(true, false, createStringPointer("option"), "test flag", "f", "", "", []string{"option1", "option2"}),
		expectedErrorStringSubset: enumOptionIndicator,
	},
	"An enum flag that is not required and is just whitespace should not return an error": {
		flag: flags.NewEnumFlag(false, false, createStringPointer("    "), "test flag", "f", "", "", []string{"option1", "option2"}),
	},
	"An enum flag that is not required and is nil should not return an error": {
		flag: flags.NewEnumFlag(false, false, nil, "test flag", "f", "", "", []string{"option1", "option2"}),
	},
	"An enum flag that is not required, is not nil or whitespace, and is one of the allowed options should not return an error": {
		flag: flags.NewEnumFlag(false, false, createStringPointer("option2"), "test flag", "f", "", "", []string{"option1", "option2"}),
	},
	"An enum flag that is not required, is not nil or whitespace, and is not one of the allowed options should return an error": {
		flag:                      flags.NewEnumFlag(false, false, createStringPointer("option3"), "test flag", "f", "", "", []string{"option1", "option2"}),
		expectedErrorStringSubset: enumOptionIndicator,
	},
}

func TestFlagValidation(t *testing.T) {
	t.Parallel()

	t.Run("String Flags", func(t *testing.T) {
		t.Parallel()
		for name, args := range stringFlagValidationTestCases {
			t.Run(name, func(t *testing.T) {
				runValidationTest(t, args)
			})
		}
	})

	t.Run("File Flags", func(t *testing.T) {
		t.Parallel()
		for name, args := range fileFlagValidationTestCases {
			t.Run(name, func(t *testing.T) {
				runValidationTest(t, args)
			})
		}
	})

	t.Run("Directory Flags", func(t *testing.T) {
		t.Parallel()
		for name, args := range directoryFlagValidationTestCases {
			t.Run(name, func(t *testing.T) {
				runValidationTest(t, args)
			})
		}
	})

	t.Run("Enum Flags", func(t *testing.T) {
		t.Parallel()
		for name, args := range enumFlagValidationTestCases {
			t.Run(name, func(t *testing.T) {
				runValidationTest(t, args)
			})
		}
	})
}

func runValidationTest(t *testing.T, args validationFlagTestCase) {
	t.Helper()

	if args.setup != nil {
		args.setup(t, args.flag)
	}

	err := args.flag.Validate()
	if args.expectedErrorStringSubset == nil {
		require.NoError(t, err)
	} else {
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), *args.expectedErrorStringSubset)
	}
}

func createStringPointer(value string) *string {
	return &value
}

func setupFiletTest(t *testing.T, flag flags.Flag) {
	t.Helper()

	fileFlag, ok := flag.(flags.FileFlag)
	if !ok {
		return
	}

	dir := t.TempDir()

	if fileFlag.Value != nil && strings.TrimSpace(*fileFlag.Value) != "" {
		var filePath = filepath.Join(dir, *fileFlag.Value)

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("failed to create test file %q: %s", filePath, err)
		}

		defer tests.MustClose(t, file)

		*fileFlag.Value = filePath
	}
}

func setupDirectoryTest(t *testing.T, flag flags.Flag) {
	t.Helper()

	directoryFlag, ok := flag.(flags.DirectoryFlag)
	if !ok {
		return
	}

	dir := t.TempDir()

	if directoryFlag.Value != nil && strings.TrimSpace(*directoryFlag.Value) != "" {
		var filePath = filepath.Join(dir, *directoryFlag.Value)

		err := os.Mkdir(filePath, 0755)
		if err != nil {
			log.Fatalf("failed to create test directory %q: %s", filePath, err)
		}

		*directoryFlag.Value = filePath
	}
}
