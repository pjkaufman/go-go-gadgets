package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type StringFlag struct {
	BaseFlag

	Value   *string
	Default string
}

func NewStringFlag(isRequired, isPersistent bool, p *string, name, shorthand, value, usage string) Flag {
	return StringFlag{
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

func (s StringFlag) Validate() error {
	if s.IsRequired {
		if s.Value == nil || strings.TrimSpace(*s.Value) == "" {
			return fmt.Errorf("%s must have a non-whitespace value", s.Name)
		}
	}

	return nil
}

func (f StringFlag) AddToCmd(cmd *cobra.Command) error {
	f.flagSet(cmd).StringVarP(
		f.Value,
		f.Name,
		f.Shorthand,
		f.Default,
		f.Usage,
	)

	return f.markRequired(cmd)
}
