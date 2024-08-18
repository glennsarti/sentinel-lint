package lint

import (
	"github.com/glennsarti/sentinel-parser/position"
)

type Issues []*Issue

// TODO : JSON serialise

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
