package sentinel_assignments

import (
	"github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-parser/position"
	"github.com/glennsarti/sentinel-parser/sentinel/ast"
)

// Walker
type walker struct {
	firstRulePos     *position.SourceRange
	firstRuleEndByte int

	issues lint.Issues
}

func (w *walker) Visit(n ast.Node) ast.VisitFunc {
	switch actual := n.(type) {
	case *ast.AssignStatement:
		w.visitAssignStatement(actual)
	}

	return w.Visit
}

func (w *walker) visitAssignStatement(n *ast.AssignStatement) {
	switch actual := n.RightExpr.(type) {
	case *ast.RuleExpression:
		if w.firstRulePos == nil {
			p := n.LeftExpr.Position().Clone()
			w.firstRulePos = &p
			w.firstRuleEndByte = actual.Position().End.Byte
		}
		return
	}
	if w.firstRulePos == nil || n.NodePos.Start.Byte == 0 || n.NodePos.Start.Byte <= w.firstRuleEndByte {
		return
	}

	r := n.LeftExpr.Position().Clone()
	w.issues = append(w.issues, &lint.Issue{
		RuleId:   assignmentsAfterRulesRuleID,
		Severity: lint.Warning,
		Summary:  "Avoid assignment after rules",
		Detail:   "Avoid assignments after rules are defined as it may cause confusion due to rules being lazily evaluated.",
		Range:    &r,
		Related: &[]lint.RelatedInformation{
			{
				Summary: "First rule",
				Range:   w.firstRulePos,
			},
		},
	})
}
