package config

type SeriesStatus string

const (
	Ongoing   SeriesStatus = "O"
	Hiatus    SeriesStatus = "H"
	Completed SeriesStatus = "C"
)

const (
	OngoingDisplay   string = "Ongoing"
	HiatusDisplay    string = "Hiatus"
	CompletedDisplay string = "Completed"
)

func IsSeriesStatus(val string) bool {
	switch val {
	case string(Ongoing):
		return true
	case string(Hiatus):
		return true
	case string(Completed):
		return true
	default:
		return false
	}
}

func SeriesStatusToDisplayText(val SeriesStatus) string {
	switch val {
	case Ongoing:
		return OngoingDisplay
	case Hiatus:
		return HiatusDisplay
	case Completed:
		return CompletedDisplay
	default:
		return ""
	}
}
