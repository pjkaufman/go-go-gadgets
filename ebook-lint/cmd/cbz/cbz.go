package cbz

import (
	"github.com/spf13/cobra"
)

// CbzCmd represents the cbz command
var CbzCmd = &cobra.Command{
	Use:   "cbz",
	Short: "Deals with cbz related commands",
	Long:  `Handles operations on cbz files in particular`,
}

func init() {
}
