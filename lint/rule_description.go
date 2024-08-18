package lint

type RuleDescriptions []RuleDescription

type RuleDescription struct {
	ID          string
	Name        string
	Description string
	URL         string
}
