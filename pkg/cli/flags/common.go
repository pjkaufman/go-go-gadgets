package flags

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func getPflagSet(cmd *cobra.Command, isPersistent bool) *pflag.FlagSet {
	if isPersistent {
		return cmd.PersistentFlags()
	}

	return cmd.Flags()
}

func markAsRequired(cmd *cobra.Command, isPersistent bool, name string) error {
	var err error
	if isPersistent {
		err = cmd.MarkPersistentFlagRequired(name)
	} else {
		err = cmd.MarkFlagRequired(name)
	}

	if err != nil {
		return fmt.Errorf("failed to mark flag %q as required on %s command: %w", name, cmd.Name(), err)
	}

	return nil
}

func markAsDirectory(cmd *cobra.Command, isPersistent bool, name string) error {
	var err error
	if isPersistent {
		err = cmd.MarkPersistentFlagDirname(name)
	} else {
		err = cmd.MarkFlagDirname(name)
	}

	if err != nil {
		return fmt.Errorf("failed to mark flag %q as a directory on %s command: %w", name, cmd.Name(), err)
	}

	return nil
}

func markAsFileTypes(cmd *cobra.Command, isPersistent bool, name string, extensions []string) error {
	if len(extensions) == 0 {
		return fmt.Errorf("failed to mark flag %q as looking for specific file types on %s command: needs at least one file extension", name, cmd.Name())
	}

	var err error
	if isPersistent {
		err = cmd.MarkPersistentFlagFilename(name, extensions...)
	} else {
		err = cmd.MarkFlagFilename(name, extensions...)
	}

	if err != nil {
		return fmt.Errorf("failed to mark flag %q as looking for specific file types on %s command: %w", name, cmd.Name(), err)
	}

	return nil
}
