//go:build generate

package cmd

import (
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	cmdhandler "github.com/pjkaufman/go-go-gadgets/pkg/cmd-handler"
)

const (
	title       = "Magnum"
	description = "Magnum is a program that checks if the list of specified light novels has any updates and notes the release dates of any new entries."
)

var (
	novelsToAdd = []string{
		"Daily Life of the Immortal King - Novel Updates?",
		"Eighth Son - Novel Updates",
	}
	publisherTypes = []config.PublisherType{config.YenPress, config.JNovelClub, config.SevenSeasEntertainment, config.OnePeaceBooks, config.VizMedia, config.HanashiMedia}
)

func getCustomValues(generationDir string) (map[string]any, error) {

	customValues := make(map[string]any)
	customValues["novelsToAdd"] = novelsToAdd

	var supportedPublishers = make([]string, len(publisherTypes))
	for i, publisherType := range publisherTypes {
		supportedPublishers[i] = config.PublisherToDisplayString(publisherType)

		if publisherType == config.SevenSeasEntertainment {
			supportedPublishers[i] += " (uses Google Cache)"
		} else if publisherType == config.OnePeaceBooks || publisherType == config.HanashiMedia {
			supportedPublishers[i] += " (uses Wikipedia)"
		}
	}
	customValues["supportedPublishers"] = supportedPublishers

	return customValues, nil
}

func init() {
	cmdhandler.AddGenerateCmd(rootCmd, title, description, []string{
		"Add more unit tests and validation for commands and parsing logic to make sure it works as intended and is easier to refactor down the road since breaking changes should be easier to catch",
	}, getCustomValues)
}
