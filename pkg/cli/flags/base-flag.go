package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type BaseFlag struct {
	IsRequired   bool
	IsPersistent bool
	Name         string
	Shorthand    string
	Usage        string
}

func (b BaseFlag) markRequired(cmd *cobra.Command) error {
	if !b.IsRequired {
		return nil
	}

	return markAsRequired(cmd, b.IsPersistent, b.Name)
}

func (b BaseFlag) flagSet(cmd *cobra.Command) *pflag.FlagSet {
	return getPflagSet(cmd, b.IsPersistent)
}
