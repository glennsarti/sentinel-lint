package helpers

import (
	"github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-parser/diagnostics"
)

func DiagToIssue(d *diagnostics.Diagnostic) *lint.Issue {
	if d == nil {
		return nil
	}
	sev := lint.Unknown
	switch d.Severity {
	case diagnostics.Error:
		sev = lint.Error
	case diagnostics.Warning:
		sev = lint.Warning
	}

	return &lint.Issue{
		RuleId:   lint.SyntaxErrorRuleID,
		Range:    d.Range,
		Detail:   d.Detail,
		Summary:  d.Summary,
		Severity: sev,
	}
}
