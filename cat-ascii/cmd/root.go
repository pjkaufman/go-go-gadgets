package cmd

import (
	"math/rand"
	"os"
	"time"

	"github.com/pjkaufman/go-go-gadgets/cat-ascii/internal/ascii"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cat-acii",
	Short: "A cat ascii art generator that displays a random cat ascii art on each invocation",
	Run: func(cmd *cobra.Command, args []string) {
		generator := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := int(generator.Int63()) % len(ascii.CAT_ASCII)

		logger.WriteInfo(ascii.CAT_ASCII[n].Ascii)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
