package config

type PublisherType string

const (
	YenPress               PublisherType = "YenPress"
	JNovelClub             PublisherType = "JNovelClub"
	SevenSeasEntertainment PublisherType = "SevenSeasEntertainment"
	OnePeaceBooks          PublisherType = "OnePeaceBooks"
	VizMedia               PublisherType = "VizMedia"
	HanashiMedia           PublisherType = "HanashiMedia"
)

func IsPublisherType(val string) bool {
	switch val {
	case string(YenPress):
		return true
	case string(JNovelClub):
		return true
	case string(SevenSeasEntertainment):
		return true
	case string(OnePeaceBooks):
		return true
	case string(VizMedia):
		return true
	case string(HanashiMedia):
		return true
	default:
		return false
	}
}
