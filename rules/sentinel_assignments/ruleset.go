package sentinel_assignments

import (
	"github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-parser/sentinel/ast"
)

var RuleSet lint.RuleSet = sentinelAssignmentsRuleSet{}

const (
	assignmentsAfterRulesRuleID = "Lint/AssignmentsAfterRules"
)

type sentinelAssignmentsRuleSet struct{}

func (r sentinelAssignmentsRuleSet) Rules() *lint.RuleDescriptions {
	return &lint.RuleDescriptions{
		{
			ID:          assignmentsAfterRulesRuleID,
			Name:        "Avoid assignments after rules",
			Description: "TODO",
		},
	}
}

func (r sentinelAssignmentsRuleSet) Run(ctx lint.RuleSetContext) (*lint.Issues, error) {
	issues := make(lint.Issues, 0)
	var err error = nil

	for _, f := range ctx.Files() {
		switch actual := f.(type) {
		case lint.PolicyFile:
			if actual.File != nil {
				i, err := r.runPolicyFile(actual)
				if err != nil {
					break
				}
				issues = append(issues, i...)
			}
		}
	}

	return &issues, err
}

func (r sentinelAssignmentsRuleSet) runPolicyFile(f lint.PolicyFile) (lint.Issues, error) {
	visitor := &walker{
		issues: make(lint.Issues, 0),
	}

	if err := ast.Walk(visitor.Visit, f.File); err != nil {
		return nil, err
	}

	return visitor.issues, nil
}
