package jnovelclub

type JSONVolumeInfo struct {
	Series struct {
		Tags          []string `json:"tags"`
		LegacyID      string   `json:"legacyId"`
		Type          int      `json:"type"`
		Title         string   `json:"title"`
		ShortTitle    string   `json:"shortTitle"`
		OriginalTitle string   `json:"originalTitle"`
		Slug          string   `json:"slug"`
		Hidden        bool     `json:"hidden"`
		Created       struct {
			Seconds string `json:"seconds"`
			Nanos   int    `json:"nanos"`
		} `json:"created"`
		Description      string `json:"description"`
		ShortDescription string `json:"shortDescription"`
		Cover            struct {
			OriginalURL  string `json:"originalUrl"`
			CoverURL     string `json:"coverUrl"`
			ThumbnailURL string `json:"thumbnailUrl"`
		} `json:"cover"`
		Following              bool        `json:"following"`
		Catchup                bool        `json:"catchup"`
		Status                 int         `json:"status"`
		Rentals                bool        `json:"rentals"`
		ID                     string      `json:"id"`
		AgeRating              int         `json:"ageRating"`
		Banner                 interface{} `json:"banner"`
		ReaderStreamingSummary int         `json:"readerStreamingSummary"`
		TopicID                int         `json:"topicId"`
	} `json:"series"`
	Volumes []struct {
		Parts []struct {
			LegacyID string `json:"legacyId"`
			Title    string `json:"title"`
			Slug     string `json:"slug"`
			Number   int    `json:"number"`
			Preview  bool   `json:"preview"`
			Hidden   bool   `json:"hidden"`
			Created  struct {
				Seconds string `json:"seconds"`
				Nanos   int    `json:"nanos"`
			} `json:"created"`
			Expiration struct {
				Seconds string `json:"seconds"`
				Nanos   int    `json:"nanos"`
			} `json:"expiration"`
			Launch struct {
				Seconds string `json:"seconds"`
				Nanos   int    `json:"nanos"`
			} `json:"launch"`
			Cover struct {
				OriginalURL  string `json:"originalUrl"`
				CoverURL     string `json:"coverUrl"`
				ThumbnailURL string `json:"thumbnailUrl"`
			} `json:"cover"`
			Progress        int    `json:"progress"`
			OriginalTitle   string `json:"originalTitle"`
			TotalMangaPages int    `json:"totalMangaPages"`
			ID              string `json:"id"`
			ShortTitle      string `json:"shortTitle"`
			RentalCoins     int    `json:"rentalCoins"`
		} `json:"parts"`
		Volume struct {
			Creators []struct {
				Name         string `json:"name"`
				Role         int    `json:"role"`
				OriginalName string `json:"originalName"`
				ID           string `json:"id"`
			} `json:"creators"`
			LegacyID          string `json:"legacyId"`
			Title             string `json:"title"`
			Slug              string `json:"slug"`
			Number            int    `json:"number"`
			OriginalPublisher string `json:"originalPublisher"`
			Label             string `json:"label"`
			Hidden            bool   `json:"hidden"`
			ForumTopicID      int    `json:"forumTopicId"`
			Created           struct {
				Seconds string `json:"seconds"`
				Nanos   int    `json:"nanos"`
			} `json:"created"`
			Publishing struct {
				Seconds string `json:"seconds"`
				Nanos   int    `json:"nanos"`
			} `json:"publishing"`
			Description      string `json:"description"`
			ShortDescription string `json:"shortDescription"`
			Cover            struct {
				OriginalURL  string `json:"originalUrl"`
				CoverURL     string `json:"coverUrl"`
				ThumbnailURL string `json:"thumbnailUrl"`
			} `json:"cover"`
			Owned                  bool   `json:"owned"`
			OriginalTitle          string `json:"originalTitle"`
			NoEbook                bool   `json:"noEbook"`
			TotalParts             int    `json:"totalParts"`
			ID                     string `json:"id"`
			ShortTitle             string `json:"shortTitle"`
			OnSale                 bool   `json:"onSale"`
			ReaderStreamingEnabled bool   `json:"readerStreamingEnabled"`
			PremiumExtras          string `json:"premiumExtras"`
		} `json:"volume"`
	} `json:"volumes"`
}
