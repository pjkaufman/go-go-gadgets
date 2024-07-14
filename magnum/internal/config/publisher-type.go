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

func PublisherToDisplayString(val PublisherType) string {
	switch val {
	case YenPress:
		return "Yen Press"
	case JNovelClub:
		return "JNovel Club"
	case SevenSeasEntertainment:
		return "Seven Seas Entertainment"
	case OnePeaceBooks:
		return "One Peace Books"
	case VizMedia:
		return "Viz Media"
	case HanashiMedia:
		return "Hanashi Media"
	default:
		return ""
	}
}
