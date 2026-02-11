package cmd

import (
	"errors"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// RemoveCmd represents the remove book info command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the provided series from the list of series to keep track of",
	Example: heredoc.Doc(`To remove a series use the following command:
	magnum series remove -n "Lady and the Tramp"
	`),
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateRemoveSeriesFlags(seriesName)
		if err != nil {
			logger.WriteError(err.Error())
		}

		seriesInfo := config.GetConfig()
		if !seriesInfo.RemoveSeriesIfExists(seriesName) {
			logger.WriteInfo("The series does not exists in the list.")

			return
		}

		config.WriteConfig(seriesInfo)

		logger.WriteInfof("The %q was removed from the series list.\n", seriesName)
	},
}

func init() {
	seriesCmd.AddCommand(RemoveCmd)

	RemoveCmd.Flags().StringVarP(&seriesName, "name", "n", "", "the name of the series")
	err := RemoveCmd.MarkFlagRequired("name")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"name\" as required on remove command: %v\n", err)
	}
}

func ValidateRemoveSeriesFlags(seriesName string) error {
	if strings.TrimSpace(seriesName) == "" {
		return errors.New(NameArgEmpty)
	}

	return nil
}
