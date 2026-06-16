package flags

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/spf13/cobra"
)

type FileFlag struct {
	BaseFlag

	Value         *string
	Default       string
	Extensions    []string
	FileMustExist bool
}

func NewFileFlag(isRequired, isPersistent bool, p *string, name, shorthand, value, usage string, extensions []string, fileMustExist bool) Flag {
	return FileFlag{
		BaseFlag: BaseFlag{
			IsRequired:   isRequired,
			IsPersistent: isPersistent,
			Name:         name,
			Shorthand:    shorthand,
			Usage:        usage,
		},
		Default:       value,
		Value:         p,
		Extensions:    extensions,
		FileMustExist: fileMustExist,
	}
}

func (s FileFlag) Validate() error {
	if s.IsRequired {
		if s.Value == nil || strings.TrimSpace(*s.Value) == "" {
			return fmt.Errorf("%s must have a non-whitespace value", s.Name)
		}
	}

	if len(s.Extensions) != 0 && s.Value != nil {
		var ext = strings.TrimPrefix(filepath.Ext(strings.TrimSpace(*s.Value)), ".")
		if !slices.Contains(s.Extensions, ext) {
			return fmt.Errorf("%s has extension %q, must have one of the following extensions: %s", s.Name, ext, strings.Join(s.Extensions, ", "))
		}
	}

	if s.FileMustExist && s.Value != nil && strings.TrimSpace(*s.Value) != "" {
		return filehandler.FileArgExists(*s.Value, s.Name)
	}

	return nil
}

func (f FileFlag) AddToCmd(cmd *cobra.Command) error {
	f.flagSet(cmd).StringVarP(
		f.Value,
		f.Name,
		f.Shorthand,
		f.Default,
		f.Usage,
	)

	if f.IsRequired {
		err := f.markRequired(cmd)

		if err != nil {
			return err
		}
	}

	return markAsFileTypes(cmd, f.IsPersistent, f.Name, f.Extensions)
}
