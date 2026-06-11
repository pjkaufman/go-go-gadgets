package flags

import (
	"github.com/spf13/cobra"
)

type BoolFlag struct {
	BaseFlag

	Value   *bool
	Default bool
}

func NewBoolFlag(isRequired, isPersistent bool, p *bool, name, shorthand string, value bool, usage string) Flag {
	return BoolFlag{
		BaseFlag: BaseFlag{
			IsRequired:   isRequired,
			IsPersistent: isPersistent,
			Name:         name,
			Shorthand:    shorthand,
			Usage:        usage,
		},
		Value:   p,
		Default: value,
	}
}

func (s BoolFlag) Validate() error {
	return nil
}

func (s BoolFlag) AddToCmd(cmd *cobra.Command) error {
	getPflagSet(cmd, s.IsPersistent).BoolVarP(s.Value, s.Name, s.Shorthand, s.Default, s.Usage)

	if s.IsRequired {
		return markAsRequired(cmd, s.IsPersistent, s.Name)
	}

	return nil
}
