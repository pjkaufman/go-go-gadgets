package cmd

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/pjkaufman/go-go-gadgets/cat-ascii/internal/ascii"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cat-ascii",
	Short: "A cat ascii art generator that displays a random cat ascii art on each invocation",
	Run: func(cmd *cobra.Command, args []string) {
		displayRandomCatAscii()
	},
}

func displayRandomCatAscii() {
	width, height := consolesize.GetConsoleSize()

	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	indices := generator.Perm(len(ascii.CAT_ASCII))

	for _, idx := range indices {
		catAscii := ascii.CAT_ASCII[idx]
		if fitsInTerminal(catAscii.Ascii, width, height) {
			logger.WriteInfo(catAscii.Ascii)
			return
		}
	}

	logger.WriteWarn("No ASCII art fits in the current terminal size.")
}

func fitsInTerminal(art string, termWidth, termHeight int) bool {
	lines := strings.Split(art, "\n")
	artHeight := len(lines)
	artWidth := 0
	for _, line := range lines {
		if len(line) > artWidth {
			artWidth = len(line)
		}
	}
	return artWidth <= termWidth && artHeight <= termHeight
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.WriteErrorf("Error: %v\n", err)
		os.Exit(1)
	}
}
