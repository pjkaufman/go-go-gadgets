package potentiallyfixableissue

type PotentiallyFixableIssue struct {
	Name                        string
	GetSuggestions              func(string) map[string]string
	IsEnabled                   *bool
	UpdateAllInstances          bool
	AddCssSectionBreakIfMissing bool
	AddCssPageBreakIfMissing    bool
}
