package sentinel_config_overrides

import (
	"fmt"

	"github.com/glennsarti/sentinel-lint/lint"
	"github.com/glennsarti/sentinel-parser/position"
	"github.com/glennsarti/sentinel-parser/sentinel_config/ast"
)

func (r sentinelConfigOverridesRuleSet) lintNoChanges(pri, ovr ast.File, addFunc issueAdder) {
	// Globals
	yieldPtrNodes(pri.Globals, ovr.Globals, func(p, o ast.Global) {
		if ast.DynamicValuePtrComparer(p.Value, o.Value) {
			uselessAttributeOverrideIssue("value", o.ValueRange, addFunc)
		}
	})

	// Imports
	yieldNodes(pri.Imports, ovr.Imports, func(p, o ast.Import) {
		noChangeForImport(p, o, addFunc)
	})

	// Mocks
	yieldPtrNodes(pri.Mocks, ovr.Mocks, func(p, o ast.Mock) {
		if !mapChanged(p.Data, o.Data, func(pNode, oNode ast.Parameter) bool {
			return ast.DynamicValuePtrComparer(pNode.Value, oNode.Value)
		}) {
			uselessAttributeOverrideIssue("data", o.DataRange, addFunc)
		}

		if p.Module != nil && o.Module != nil {
			if noStringChange(p.Module.Source, o.Module.Source) {
				uselessAttributeOverrideIssue("source", o.Module.SourceRange, addFunc)
			}
		}
	})

	// Params
	yieldPtrNodes(pri.Params, ovr.Params, func(p, o ast.Parameter) {
		if ast.DynamicValuePtrComparer(p.Value, o.Value) {
			uselessAttributeOverrideIssue("value", o.ValueRange, addFunc)
		}
	})

	// Policies
	yieldPtrNodes(pri.Policies, ovr.Policies, func(p, o ast.Policy) {
		if noStringChange(p.Source, o.Source) {
			uselessAttributeOverrideIssue("source", o.SourceRange, addFunc)
		}
		if noStringChange(p.EnforcementLevel, o.EnforcementLevel) {
			uselessAttributeOverrideIssue("enforcement_level", o.EnforcementLevelRange, addFunc)
		}

		if !mapChanged(p.Params, o.Params, func(pNode, oNode ast.Parameter) bool {
			return ast.DynamicValuePtrComparer(pNode.Value, oNode.Value)
		}) {
			uselessAttributeOverrideIssue("params", o.ParamsRange, addFunc)
		}
	})

	// SentinelOptions
	if pri.SentinelOptions != nil && ovr.SentinelOptions != nil {
		noChangeForSentinelOptions(*pri.SentinelOptions, *ovr.SentinelOptions, addFunc)
	}

	// Test
	if pri.Test != nil && ovr.Test != nil {
		addFunc(&lint.Issue{
			RuleId:   uselessOverrideRuleID,
			Summary:  "Block has no effect",
			Severity: lint.Warning,
			Detail:   "A test block cannot be overridden.",
			Range:    ovr.Test.TestRange,
		})
	}
}

func uselessAttributeOverrideIssue(name string, attrRange *position.SourceRange, addFunc issueAdder) {
	addFunc(&lint.Issue{
		RuleId:   uselessOverrideRuleID,
		Summary:  "Attribute has no effect",
		Severity: lint.Information,
		Detail:   fmt.Sprintf("The %q attribute has no effect.", name),
		Range:    attrRange,
	})
}

func noChangeForSentinelOptions(pri, ovr ast.SentinelOptions, addFunc issueAdder) {
	if len(pri.Features) == 0 || len(ovr.Features) == 0 {
		return
	}

	priMap := map[string]ast.DynamicValue{}
	for _, feat := range pri.Features {
		if feat == nil || feat.Value == nil {
			continue
		}
		priMap[feat.Name] = *feat.Value
	}

	for _, feat := range ovr.Features {
		if feat == nil || feat.Value == nil {
			continue
		}

		priValue, ok := priMap[feat.Name]
		if !ok {
			continue
		}

		if ast.DynamicValuePtrComparer(&priValue, feat.Value) {
			r := position.SourceRange{
				Filename: feat.ValueRange.Filename,
				Start:    feat.NameRange.Start,
				End:      feat.ValueRange.End,
			}
			uselessAttributeOverrideIssue(feat.Name, &r, addFunc)
		}
	}
}

// Import Comparisons
func noChangeForImport(pri, ovr ast.Import, addFunc issueAdder) {
	switch priImport := pri.(type) {
	case *ast.V1ModuleImport:
		noChangeForV1ModuleImport(priImport, ovr, addFunc)
	case *ast.V1PluginImport:
		noChangeForV1PluginImport(priImport, ovr, addFunc)
	case *ast.V2ModuleImport:
		noChangeForV2ModuleImport(priImport, ovr, addFunc)
	case *ast.V2PluginImport:
		noChangeForV2PluginImport(priImport, ovr, addFunc)
	case *ast.V2StaticImport:
		noChangeForV2StaticImport(priImport, ovr, addFunc)
	}
}

func noChangeForV1ModuleImport(pri *ast.V1ModuleImport, otherImport ast.Import, addFunc issueAdder) {
	if ovr, ok := otherImport.(*ast.V1ModuleImport); ok {
		if noStringChange(pri.Source, ovr.Source) {
			uselessAttributeOverrideIssue("source", ovr.SourceRange, addFunc)
		}
	}
}

func noChangeForV1PluginImport(pri *ast.V1PluginImport, otherImport ast.Import, addFunc issueAdder) {
	if ovr, ok := otherImport.(*ast.V1PluginImport); ok {
		if noStringChange(pri.Path, ovr.Path) {
			uselessAttributeOverrideIssue("path", ovr.PathRange, addFunc)
		}

		if noStringListChange(pri.Args, ovr.Args) {
			uselessAttributeOverrideIssue("args", ovr.ArgsRange, addFunc)
		}

		if !mapChanged(pri.Config, ovr.Config, func(pNode, oNode ast.Parameter) bool {
			return ast.DynamicValuePtrComparer(pNode.Value, oNode.Value)
		}) {
			uselessAttributeOverrideIssue("config", ovr.ConfigRange, addFunc)
		}

		if noStringListChange(pri.Env, ovr.Env) {
			uselessAttributeOverrideIssue("env", ovr.EnvRange, addFunc)
		}
	}
}

func noChangeForV2ModuleImport(pri *ast.V2ModuleImport, otherImport ast.Import, addFunc issueAdder) {
	if ovr, ok := otherImport.(*ast.V2ModuleImport); ok {
		if noStringChange(pri.Source, ovr.Source) {
			uselessAttributeOverrideIssue("source", ovr.SourceRange, addFunc)
		}
	}
}

func noChangeForV2PluginImport(pri *ast.V2PluginImport, otherImport ast.Import, addFunc issueAdder) {
	if ovr, ok := otherImport.(*ast.V2PluginImport); ok {
		if noStringChange(pri.Source, ovr.Source) {
			uselessAttributeOverrideIssue("source", ovr.SourceRange, addFunc)
		}

		if noStringListChange(pri.Args, ovr.Args) {
			uselessAttributeOverrideIssue("args", ovr.ArgsRange, addFunc)
		}

		if !mapChanged(pri.Config, ovr.Config, func(pNode, oNode ast.Parameter) bool {
			return ast.DynamicValuePtrComparer(pNode.Value, oNode.Value)
		}) {
			uselessAttributeOverrideIssue("config", ovr.ConfigRange, addFunc)
		}

		if !mapChanged(pri.Env, ovr.Env, func(pNode, oNode ast.Parameter) bool {
			return ast.DynamicValuePtrComparer(pNode.Value, oNode.Value)
		}) {
			uselessAttributeOverrideIssue("env", ovr.EnvRange, addFunc)
		}
	}
}

func noChangeForV2StaticImport(pri *ast.V2StaticImport, otherImport ast.Import, addFunc issueAdder) {
	if ovr, ok := otherImport.(*ast.V2StaticImport); ok {
		if noStringChange(pri.Source, ovr.Source) {
			uselessAttributeOverrideIssue("source", ovr.SourceRange, addFunc)
		}
		if noStringChange(pri.Format, ovr.Format) {
			uselessAttributeOverrideIssue("format", ovr.FormatRange, addFunc)
		}
	}
}
