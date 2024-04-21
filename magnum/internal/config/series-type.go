package config

type SeriesType string

const (
	WebNovel   SeriesType = "WN"
	Manga      SeriesType = "MN"
	LightNovel SeriesType = "LN"
)

const (
	WebNovelDisplay   string = "Web Novel"
	MangaDisplay      string = "Manga"
	LightNovelDisplay string = "Light Novel"
)

func IsSeriesType(val string) bool {
	switch val {
	case string(WebNovel):
		return true
	case string(Manga):
		return true
	case string(LightNovel):
		return true
	default:
		return false
	}
}

func SeriesTypeToDisplayText(val SeriesType) string {
	switch val {
	case WebNovel:
		return WebNovelDisplay
	case Manga:
		return MangaDisplay
	case LightNovel:
		return LightNovelDisplay
	default:
		return ""
	}
}
