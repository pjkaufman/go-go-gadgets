package flags

import (
	"github.com/spf13/cobra"
)

type IntFlag struct {
	BaseFlag

	Value   *int
	Default int
}

func NewIntFlag(isRequired, isPersistent bool, p *int, name, shorthand string, value int, usage string) Flag {
	return IntFlag{
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

func (f IntFlag) AddToCmd(cmd *cobra.Command) error {
	f.flagSet(cmd).IntVarP(
		f.Value,
		f.Name,
		f.Shorthand,
		f.Default,
		f.Usage,
	)

	return f.markRequired(cmd)
}

func (s IntFlag) Validate() error {
	return nil
}
