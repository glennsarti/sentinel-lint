package sentinel_config_basic

import (
	"fmt"
	"sort"

	"github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-parser/features"
	"github.com/glennsarti/sentinel-parser/sentinel_config/ast"
)

var RuleSet lint.RuleSet = sentinelConfigBasicRuleSet{}

const (
	conflictNameRuleID  = "Lint/ConflictingName"
	duplicateNameRuleID = "Lint/DuplicateName"
)

type sentinelConfigBasicRuleSet struct{}

func (r sentinelConfigBasicRuleSet) Rules() *lint.RuleDescriptions {
	return &lint.RuleDescriptions{
		{
			ID:          conflictNameRuleID,
			Name:        "Names should not conflict with standard imports.",
			Description: "TODO",
		},
		{
			ID:          duplicateNameRuleID,
			Name:        "The same name should not be use for different purposes.",
			Description: "TODO",
		},
	}
}

func (r sentinelConfigBasicRuleSet) Run(ctx lint.RuleSetContext) (*lint.Issues, error) {
	issues := make(lint.Issues, 0)
	var err error = nil

	for _, f := range ctx.Files() {
		switch actual := f.(type) {
		case lint.ConfigPrimaryFile:
			if actual.ConfigFile != nil {
				i, err := r.runSentinelConfigFile(actual.ConfigFile, ctx.SentinelVersion())
				if err != nil {
					break
				}
				issues = append(issues, i...)
			}
		}
	}

	return &issues, err
}

func (r sentinelConfigBasicRuleSet) runSentinelConfigFile(f *ast.File, sentinelVersion string) (lint.Issues, error) {
	if f == nil {
		return nil, nil
	}

	issues := make(lint.Issues, 0)

	for _, imp := range f.Imports {
		if level, ver := r.isStdImportName(imp.BlockName(), sentinelVersion); level != nil {
			issues = append(issues, r.newConflictNameIssue(imp, *level, ver))
		}
	}

	allNames := findNameLocations(f)
	if allNames != nil {
		for name, list := range *allNames {
			if len(list) > 1 {
				issues = append(issues, r.newDuplicateNameIssue(name, list))
			}
		}
	}

	return issues, nil
}

func (r sentinelConfigBasicRuleSet) newDuplicateNameIssue(name string, list nameLocationList) *lint.Issue {
	sort.Slice(list, func(i, j int) bool {
		return list[j].Range.Start.Byte > list[i].Range.Start.Byte
	})

	related := make([]lint.RelatedInformation, len(list)-1)
	for i := 1; i < len(list); i++ {
		related[i-1] = lint.RelatedInformation{
			Range:   &list[i].Range,
			Summary: fmt.Sprintf("%s '%s' definition", list[i].Type, name),
		}
	}

	pos := list[0].Range.Clone()

	return &lint.Issue{
		RuleId:   duplicateNameRuleID,
		Summary:  "Block uses a duplicate name",
		Severity: lint.Warning,
		Detail:   fmt.Sprintf("The %s with name %q uses a duplicate name", list[0].Type, name),
		Range:    &pos,
		Related:  &related,
	}
}

func (r sentinelConfigBasicRuleSet) isStdImportName(name, sentinelVersion string) (*lint.SeverityLevel, string) {
	if ver, ok := importToVersion[name]; ok {
		if features.SupportedVersion(sentinelVersion, ver) {
			return &lint.Error, ver
		} else {
			return &lint.Warning, ver
		}
	}
	return nil, ""
}

func (r sentinelConfigBasicRuleSet) newConflictNameIssue(block ast.HCLNode, level lint.SeverityLevel, fromSentinelVersion string) *lint.Issue {
	detail := fmt.Sprintf("The %s uses the name %q which conflicts with a standard import.", block.BlockType(), block.BlockName())
	if level == lint.Warning {
		detail = fmt.Sprintf("The %s uses the name %q, which will conflict with a standard import in Sentinel version %s onwards", block.BlockType(), block.BlockName(), fromSentinelVersion)
	}

	pos := block.Range().Clone()

	return &lint.Issue{
		RuleId:   conflictNameRuleID,
		Summary:  "Block uses a conflicting name",
		Severity: level,
		Detail:   detail,
		Range:    &pos,
	}
}
