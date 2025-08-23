package potentiallyfixableissue

type PotentiallyFixableIssue struct {
	Name                        string
	GetSuggestions              func(string) (map[string]string, error)
	IsEnabled                   *bool
	UpdateAllInstances          bool
	AddCssSectionBreakIfMissing bool
	AddCssPageBreakIfMissing    bool
}
