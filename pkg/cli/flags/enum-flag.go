package flags

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

type EnumFlag struct {
	StringFlag

	Allowed []string
}

const CustomEnumValuesFlagAnnotation = "custom_enum_values_flag_annotation"

func NewEnumFlag(isRequired, isPersistent bool, p *string, name, shorthand, value, usage string, allowed []string) Flag {
	return EnumFlag{
		StringFlag: StringFlag{
			BaseFlag: BaseFlag{
				IsRequired:   isRequired,
				IsPersistent: isPersistent,
				Name:         name,
				Shorthand:    shorthand,
				Usage:        usage,
			},
			Value:   p,
			Default: value,
		},
		Allowed: allowed,
	}
}

func (f EnumFlag) AddToCmd(cmd *cobra.Command) error {
	if len(f.Allowed) == 0 {
		return fmt.Errorf("failed to add flag %q on command %s because no allowed values were provided.", f.Name, cmd.Name())
	}

	err := f.StringFlag.AddToCmd(cmd)
	if err != nil {
		return err
	}

	flag := f.flagSet(cmd).Lookup(f.Name)
	if flag != nil {
		if flag.Annotations == nil {
			flag.Annotations = make(map[string][]string)
		}

		flag.Annotations[CustomEnumValuesFlagAnnotation] = f.Allowed
	}

	err = cmd.RegisterFlagCompletionFunc(f.Name, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return f.Allowed, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return fmt.Errorf("failed to add custom shell completion registration for flag %q on %s command: %w", f.Name, cmd.Name(), err)
	}

	return nil
}

func (f EnumFlag) Validate() error {
	if err := f.StringFlag.Validate(); err != nil {
		return err
	}

	if f.Value == nil || strings.TrimSpace(*f.Value) == "" || slices.Contains(f.Allowed, *f.Value) {
		return nil
	}

	return fmt.Errorf("%s must be one of %s", f.Name, strings.Join(f.Allowed, ", "))
}
