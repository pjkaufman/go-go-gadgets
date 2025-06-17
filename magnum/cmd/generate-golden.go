/*go:build generate_test*/

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var goldenFilePath string

// GenerateTestCmd represents the generate golden test file command
var GenerateTestCmd = &cobra.Command{
	Use:   "golden",
	Short: "Gets a snapshot of the current html for several series to use for golden testing",
	Run: func(cmd *cobra.Command, args []string) {
		if goldenFilePath == "" {
			logger.WriteErrorf("output path is required\n")
			return
		}

		type goldenTestInfo struct {
			url      string
			filename string
		}

		var goldenList = []goldenTestInfo{
			// JNovel-Club
			{
				url:      "https://j-novel.club/series/arifureta-zero",
				filename: "jnovel-club/test/arifureta-zero.golden",
			},
			{
				url:      "https://j-novel.club/series/how-a-realist-hero-rebuilt-the-kingdom",
				filename: "jnovel-club/test/how-a-realist-hero-rebuilt-the-kingdom.golden",
			},
			// Seven Seas Entertainment
			// TODO: fix the forbidden errors
			{
				url:      "https://sevenseasentertainment.com/series/mushoku-tensei-jobless-reincarnation-light-novel/",
				filename: "sevenseasentertainment/test/mushoku-tensei-jobless-reincarnation-light-novel.golden",
			},
			{
				url:      "https://sevenseasentertainment.com/series/berserk-of-gluttony-light-novel/",
				filename: "sevenseasentertainment/test/berserk-of-gluttony-light-novel.golden",
			},
		}

		for _, test := range goldenList {
			err := createGoldenFile(test.url, filepath.Join(goldenFilePath, test.filename))

			if err != nil {
				logger.WriteErrorf("failed to create golden file for %s: %v", test.url, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(GenerateTestCmd)

	GenerateTestCmd.Flags().StringVarP(&goldenFilePath, "out", "o", "", "the output path for where to store the resulting golden files for the tests which should point to magnum's internal folder")
}

func createGoldenFile(url string, out string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "text/html")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request for url %q: %w", url, err)
	}
	defer resp.Body.Close()

	file, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("error creating file %q: %w", out, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file %q: %w", out, err)
	}

	logger.WriteInfof("Response saved to %q\n", out)

	return nil
}
