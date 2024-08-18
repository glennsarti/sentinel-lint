package lint

type RuleSetList = []RuleSet

type RuleSet interface {
	Run(ctx RuleSetContext) (*Issues, error)

	Rules() *RuleDescriptions
}
