package jnovelclub

type JSONVolumeInfo struct {
	Props struct {
		PageProps struct {
			Aggregate struct {
				Volumes []struct {
					Volume struct {
						Title             string `json:"title,omitempty"`
						Slug              string `json:"slug,omitempty"`
						Number            int    `json:"number,omitempty"`
						OriginalPublisher string `json:"originalPublisher,omitempty"`
						Label             string `json:"label,omitempty"`
						Hidden            bool   `json:"hidden,omitempty"`
						ForumTopicID      int    `json:"forumTopicId,omitempty"`
						Created           struct {
							Seconds string `json:"seconds,omitempty"`
							Nanos   int    `json:"nanos,omitempty"`
						} `json:"created,omitempty"`
						Publishing struct {
							Seconds string `json:"seconds,omitempty"`
							Nanos   int    `json:"nanos,omitempty"`
						} `json:"publishing,omitempty"`
						TotalParts int `json:"totalParts,omitempty"`
					} `json:"volume,omitempty"`
				} `json:"volumes,omitempty"`
			} `json:"aggregate,omitempty"`
			ID string `json:"id,omitempty"`
		} `json:"pageProps,omitempty"`
	} `json:"props,omitempty"`
}
