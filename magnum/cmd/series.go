package cmd

import (
	"github.com/spf13/cobra"
)

// seriesCmd represents the series command
var seriesCmd = &cobra.Command{
	Use:   "series",
	Short: "Deals with series related information that magnum tracks",
}

func init() {
	rootCmd.AddCommand(seriesCmd)
}
