package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pjkaufman/go-go-gadgets/magnum/internal/config"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func selectBookName(series []config.SeriesInfo, includeCompleted bool) string {
	var seriesNames = []string{}
	for _, series := range series {
		if series.Status != config.Completed || includeCompleted {
			seriesNames = append(seriesNames, series.Name)
		}
	}

	prompt := promptui.Select{
		Label: "Select Book Name",
		Items: seriesNames,
		Searcher: func(input string, index int) bool {
			seriesName := seriesNames[index]
			name := strings.Replace(strings.ToLower(seriesName), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, input)
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		logger.WriteErrorf("Book name prompt failed %v\n", err)
	}

	return result
}

func selectBookStatus() config.SeriesStatus {
	var statuses = []config.SeriesStatus{
		config.Ongoing,
		config.Hiatus,
		config.Completed,
	}
	var seriesStatuses = make([]string, len(statuses))
	for i, status := range statuses {
		seriesStatuses[i] = fmt.Sprintf("%s - %s", status, config.SeriesStatusToDisplayText(status))
	}

	prompt := promptui.Select{
		Label: "Select Book Status",
		Items: seriesStatuses,
	}

	i, _, err := prompt.Run()
	if err != nil {
		logger.WriteErrorf("Book status prompt failed %v\n", err)
	}

	return statuses[i]
}

func selectPublisher() config.PublisherType {
	var publishers = []config.PublisherType{
		config.YenPress,
		config.JNovelClub,
		config.SevenSeasEntertainment,
		config.OnePeaceBooks,
		config.VizMedia,
		config.HanashiMedia,
	}
	var publisherTypes = make([]string, len(publishers))
	for i, publisherType := range publishers {
		publisherTypes[i] = fmt.Sprintf("%[1]s - %[1]s", publisherType)
	}

	prompt := promptui.Select{
		Label: "Select Book Publisher",
		Items: publisherTypes,
	}

	i, _, err := prompt.Run()
	if err != nil {
		logger.WriteErrorf("Book publisher prompt failed %v\n", err)
	}

	return publishers[i]
}

func selectSeriesType() config.SeriesType {
	var types = []config.SeriesType{
		config.WebNovel,
		config.Manga,
		config.LightNovel,
	}
	var seriesTypes = make([]string, len(types))
	for i, seriesType := range types {
		seriesTypes[i] = fmt.Sprintf("%s - %s", seriesType, config.SeriesTypeToDisplayText(seriesType))
	}

	prompt := promptui.Select{
		Label: "Select Series Type",
		Items: seriesTypes,
	}

	i, _, err := prompt.Run()
	if err != nil {
		logger.WriteErrorf("Book series type prompt failed %v\n", err)
	}

	return types[i]
}
