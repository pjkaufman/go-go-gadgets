package cbr

import (
	"github.com/spf13/cobra"
)

// CbrCmd represents the cbr command
var CbrCmd = &cobra.Command{
	Use:   "cbr",
	Short: "Deals with cbr related commands",
	Long:  `Handles operations on cbr files in particular`,
}

func init() {
}
