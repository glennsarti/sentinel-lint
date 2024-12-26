package rules

import (
	"github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-lint/rules/sentinel_assignments"
	"github.com/glennsarti/sentinel-lint/rules/sentinel_config_basic"
	"github.com/glennsarti/sentinel-lint/rules/sentinel_config_overrides"
)

func NewDefaultRuleSet() lint.RuleSetList {
	return lint.RuleSetList{
		sentinel_assignments.RuleSet,
		sentinel_config_basic.RuleSet,
		sentinel_config_overrides.RuleSet,
	}
}
