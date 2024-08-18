package lint

type RuleSetContext interface {
	Files() []File
	SentinelVersion() string
}
