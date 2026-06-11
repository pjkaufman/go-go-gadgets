package flags

import "github.com/spf13/cobra"

type Flag interface {
	Validate() error
	AddToCmd(cmd *cobra.Command) error
}
