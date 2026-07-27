package hanashimedia

import "time"

type JSONVolumeInfo struct {
	Data struct {
		At          string `json:"at"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
		Entries     []struct {
			Order int `json:"order"`
			Ebook struct {
				At           string    `json:"at"`
				Title        string    `json:"title"`
				Subtitle     string    `json:"subtitle"`
				Description  string    `json:"description"`
				Visible      bool      `json:"visible"`
				Categories   []string  `json:"categories"`
				Release      time.Time `json:"release"`
				Authors      []string  `json:"authors"`
				Nsfw         bool      `json:"nsfw"`
				AnimeAdapted bool      `json:"animeAdapted"`
				Status       string    `json:"status"`
			} `json:"ebook"`
		} `json:"entries"`
		SeriesStatus string `json:"seriesStatus"`
	} `json:"data"`
}
