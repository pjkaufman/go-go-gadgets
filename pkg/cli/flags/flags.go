package flags

import "github.com/spf13/cobra"

type Flags struct {
	Flags []Flag
}

func (f Flags) Validate() error {
	for _, flag := range f.Flags {
		err := flag.Validate()

		if err != nil {
			return err
		}
	}

	return nil
}

func (f Flags) AddToCmd(cmd *cobra.Command) error {
	for _, flag := range f.Flags {
		err := flag.AddToCmd(cmd)

		if err != nil {
			return err
		}
	}

	return nil
}
