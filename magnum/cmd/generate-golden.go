/*go:build generate_test*/

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
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
			// this value indicates that the html pulled at the time of generation should be the source of truth
			// if the page changes with regards to how the parsing needs to happen, then it needs to be updated
			// manually to reflect the change in html structure
			frozen bool
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
			{
				url:      "https://sevenseasentertainment.com/series/mushoku-tensei-jobless-reincarnation-light-novel/",
				filename: "sevenseasentertainment/test/mushoku-tensei-jobless-reincarnation-light-novel.golden",
			},
			{
				url:      "https://sevenseasentertainment.com/series/berserk-of-gluttony-light-novel/",
				filename: "sevenseasentertainment/test/berserk-of-gluttony-light-novel.golden",
			},
			// Viz Media
			{
				url:      "https://www.viz.com/nausicaa-of-the-valley-of-the-wind",
				filename: "vizmedia/test/nausicaa-of-the-valley-of-the-wind.golden",
			},
			{
				url:      "https://www.viz.com/manga-books/nausicaa-of-the-valley-of-the-wind/section/115444/more",
				filename: "vizmedia/test/manga-books-nausicaa-of-the-valley-of-the-wind-section-115444-more.golden",
			},
			// Yen Press
			{
				url:      "https://yenpress.com/series/the-asterisk-war",
				filename: "yenpress/test/the-asterisk-war.golden",
			},
			{
				url:      "https://yenpress.com/titles/9781975369095-the-asterisk-war-vol-17-light-novel",
				filename: "yenpress/test/titles-9781975369095-the-asterisk-war-vol-17-light-novel.golden",
			},
			{
				url:      "https://yenpress.com/series/a-certain-magical-index-light-novel",
				filename: "yenpress/test/a-certain-magical-index-light-novel.golden",
				frozen:   true, // meant to test the omnibus ignore logic
			},
			{
				url:      "https://yenpress.com/titles/9781975317997-a-certain-magical-index-ss-vol-2-light-novel",
				filename: "yenpress/test/titles-9781975317997-a-certain-magical-index-ss-vol-2-light-novel.golden",
			},
		}

		for _, test := range goldenList {
			err := createGoldenFile(test.url, filepath.Join(goldenFilePath, test.filename), test.frozen)

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

func createGoldenFile(url, out string, frozen bool) error {
	if frozen {
		exists, err := filehandler.FileExists(out)
		if err != nil {
			return err
		}

		if exists {
			return nil
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "text/html")
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request for url %q: %w", url, err)
	}
	defer resp.Body.Close()

	// TODO: make sure to create the folder if it does not already exist
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
