package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/cli/flags"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var seriesRemoveFlags = flags.Flags{
	Flags: []flags.Flag{
		flags.NewStringFlag(true, false, &seriesName, "name", "n", "", "the name of the series"),
	},
}

// RemoveCmd represents the remove book info command
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the provided series from the list of series to keep track of",
	Example: heredoc.Doc(`To remove a series use the following command:
	magnum series remove -n "Lady and the Tramp"
	`),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return seriesRemoveFlags.Validate()
	},
	Run: func(cmd *cobra.Command, args []string) {
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

	err := seriesRemoveFlags.AddToCmd(RemoveCmd)
	if err != nil {
		logger.WriteFatal(err.Error())
	}
}
