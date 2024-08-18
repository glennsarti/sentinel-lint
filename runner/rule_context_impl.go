package runner

import (
	lint "github.com/glennsarti/sentinel-lint/lint"
)

var _ lint.RuleSetContext = ruleContext{}

type ruleContext struct {
	files           []lint.File
	sentinelVersion string
}

func (rc ruleContext) Files() []lint.File {
	return rc.files
}

func (rc ruleContext) SentinelVersion() string {
	return rc.sentinelVersion
}
