package flags

import (
	"fmt"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/spf13/cobra"
)

type DirectoryFlag struct {
	BaseFlag

	Value   *string
	Default string
}

func NewDirectoryFlag(isRequired, isPersistent bool, p *string, name, shorthand, value, usage string) Flag {
	return DirectoryFlag{
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

func (d DirectoryFlag) Validate() error {
	if d.IsRequired {
		if d.Value == nil || strings.TrimSpace(*d.Value) == "" {
			return fmt.Errorf("%s must have a non-whitespace value", d.Name)
		}
	}

	if d.Value != nil && strings.TrimSpace(*d.Value) != "" {
		return filehandler.FolderArgExists(*d.Value, d.Name)
	}

	return nil
}

func (d DirectoryFlag) AddToCmd(cmd *cobra.Command) error {
	d.flagSet(cmd).StringVarP(
		d.Value,
		d.Name,
		d.Shorthand,
		d.Default,
		d.Usage,
	)

	if d.IsRequired {
		err := markAsRequired(cmd, d.IsPersistent, d.Name)

		if err != nil {
			return err
		}
	}

	return markAsDirectory(cmd, d.IsPersistent, d.Name)
}
