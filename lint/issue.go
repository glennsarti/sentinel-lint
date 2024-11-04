package lint

import (
	"github.com/glennsarti/sentinel-parser/position"
)

type Issues []*Issue

func (i Issues) HasErrors() bool {
	for _, issue := range i {
		if issue != nil && issue.Severity == Error {
			return true
		}
	}
	return false
}

type Issue struct {
	Severity SeverityLevel

	RuleId string

	Summary string
	Detail  string
	Range   *position.SourceRange

	Related *[]RelatedInformation
}

type RelatedInformation struct {
	Range   *position.SourceRange
	Summary string
}
