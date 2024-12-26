package sentinel_config_overrides

import (
	"fmt"

	"github.com/glennsarti/sentinel-lint/internal/helpers"
	"github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-parser/features"
	"github.com/glennsarti/sentinel-parser/sentinel_config/ast"
	"github.com/glennsarti/sentinel-parser/sentinel_config/parser"
)

var RuleSet lint.RuleSet = sentinelConfigOverridesRuleSet{}

type issueAdder func(i *lint.Issue)

const (
	uselessOverrideRuleID = "Lint/UselessOverride"
)

type sentinelConfigOverridesRuleSet struct{}

func (r sentinelConfigOverridesRuleSet) Rules() *lint.RuleDescriptions {
	return &lint.RuleDescriptions{
		{
			ID:          uselessOverrideRuleID,
			Name:        "Overriding a value should have an effect.",
			Description: "TODO",
		},
	}
}

func (r sentinelConfigOverridesRuleSet) Run(ctx lint.RuleSetContext) (*lint.Issues, error) {
	issues := make(lint.Issues, 0)
	var err error = nil

	for _, f := range ctx.Files() {
		switch actual := f.(type) {
		case lint.ConfigOverrideFile:
			i, err := r.runSentinelOverrideFile(actual, ctx.SentinelVersion())
			if err != nil {
				break
			}
			issues = append(issues, i...)
		}
	}

	return &issues, err
}

func (r sentinelConfigOverridesRuleSet) runSentinelOverrideFile(f lint.ConfigOverrideFile, sentinelVersion string) (lint.Issues, error) {
	issues := make(lint.Issues, 0)
	if f.ConfigFile == nil {
		return issues, nil
	}

	if features.UnsupportedVersion(sentinelVersion, features.ConfigurationOverrideMinimumVersion) {
		// TODO: Should this raise a lint issue?
		return issues, nil
	}

	r.checkEmptyBlocks(*f.ConfigFile, func(i *lint.Issue) {
		issues = append(issues, i)
	})

	if f.PrimaryFile == nil {
		return issues, nil
	}

	// Attempt to override the primary file with the config file
	diags := parser.AttemptOverrideFileWith(f.PrimaryFile, f.ConfigFile, sentinelVersion)
	for _, d := range diags {
		if d != nil {
			issues = append(issues, helpers.DiagToIssue(d))
		}
	}
	if issues.HasErrors() {
		return issues, nil
	}

	// TODO Check for overrides with the same value

	return issues, nil
}

func (r sentinelConfigOverridesRuleSet) checkEmptyBlocks(f ast.File, addFunc issueAdder) {
	// Globals
	for _, blk := range f.Globals {
		if blk.ValueRange == nil {
			r.emptyBlockViolation(blk, addFunc)
		}
	}

	// Imports
	for _, blk := range f.Imports {
		switch actual := blk.(type) {
		case *ast.V1ModuleImport:
			if actual.SourceRange == nil {
				r.emptyBlockViolation(blk, addFunc)
			}
		case *ast.V1PluginImport:
			if actual.PathRange == nil &&
				actual.ArgsRange == nil &&
				actual.EnvRange == nil &&
				actual.ConfigRange == nil {
				r.emptyBlockViolation(blk, addFunc)
			}
		case *ast.V2ModuleImport:
			if actual.SourceRange == nil {
				r.emptyBlockViolation(blk, addFunc)
			}
		case *ast.V2PluginImport:
			if actual.SourceRange == nil &&
				actual.ArgsRange == nil &&
				actual.EnvRange == nil &&
				actual.ConfigRange == nil {
				r.emptyBlockViolation(blk, addFunc)
			}
		case *ast.V2StaticImport:
			if actual.SourceRange == nil && actual.FormatRange == nil {
				r.emptyBlockViolation(blk, addFunc)
			}
		}
	}

	// Mocks
	for _, blk := range f.Mocks {
		if blk.DataRange == nil && blk.Module == nil {
			r.emptyBlockViolation(blk, addFunc)
		}
	}

	// Params
	for _, blk := range f.Params {
		if blk.ValueRange == nil {
			r.emptyBlockViolation(blk, addFunc)
		}
	}

	// Policies
	for _, blk := range f.Policies {
		if blk.EnforcementLevelRange == nil &&
			blk.SourceRange == nil &&
			blk.ParamsRange == nil {
			r.emptyBlockViolation(blk, addFunc)
		}
	}

	// SentinelOptions
	if f.SentinelOptions != nil && f.SentinelOptions.FeaturesRange == nil {
		r.emptyBlockViolation(f.SentinelOptions, addFunc)
	}

	// Test
	if f.Test != nil && f.Test.TestRange != nil {
		r.emptyBlockViolation(f.Test, addFunc)
	}
}

func (r sentinelConfigOverridesRuleSet) emptyBlockViolation(block ast.HCLNode, addFunc issueAdder) {
	i := &lint.Issue{
		RuleId:   uselessOverrideRuleID,
		Summary:  "Block has no effect",
		Severity: lint.Information,
		Detail:   fmt.Sprintf("The %s block is empty and has no effect.", block.BlockType()),
		Range:    block.Range(),
	}

	addFunc(i)
}
