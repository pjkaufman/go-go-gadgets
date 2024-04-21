package epub

import (
	"github.com/spf13/cobra"
)

// EpubCmd represents the epub command
var EpubCmd = &cobra.Command{
	Use:   "epub",
	Short: "Deals with epub related commands",
	Long:  `Handles operations on epub files in particular`,
}

func init() {
}
