package parsing

import (
	"github.com/glennsarti/sentinel-lint/lint"

	"github.com/glennsarti/sentinel-parser/diagnostics"
	sast "github.com/glennsarti/sentinel-parser/sentinel/ast"
	sp "github.com/glennsarti/sentinel-parser/sentinel/parser"
	"github.com/glennsarti/sentinel-parser/sentinel/token"
	scast "github.com/glennsarti/sentinel-parser/sentinel_config/ast"
	scp "github.com/glennsarti/sentinel-parser/sentinel_config/parser"
)

func diagsToIssues(diags diagnostics.Diagnostics) lint.Issues {
	list := make(lint.Issues, 0)
	for _, diag := range diags {
		if diag != nil && diag.Severity == diagnostics.Error {
			list = append(list, &lint.Issue{
				RuleId:   lint.SyntaxErrorRuleID,
				Summary:  diag.Summary,
				Detail:   diag.Detail,
				Range:    diag.Range,
				Severity: diagSeverityToIssueSeverity(diag.Severity),
			})
		}
	}
	return list
}

func diagSeverityToIssueSeverity(sev diagnostics.SeverityLevel) lint.SeverityLevel {
	switch sev {
	case diagnostics.Error:
		return lint.Error
	case diagnostics.Warning:
		return lint.Warning
	default:
		return lint.Unknown
	}
}

func ParseSentinelFile(
	sentinelVersion string,
	filename string,
	src []byte,
) (*sast.File, token.Locator, lint.Issues, error) {
	f, loc, diag, err := sp.ParseFile(sentinelVersion, filename, src)

	return f, loc, diagsToIssues(diag), err
}

func ParseSentinelConfigFile(
	sentinelVersion string,
	filename string,
	src []byte,
) (*scast.File, lint.Issues, error) {

	p, err := scp.New(sentinelVersion)
	if err != nil {
		return nil, nil, err
	}

	f, diag := p.ParseFile(filename, src)

	return f, diagsToIssues(diag), err
}
