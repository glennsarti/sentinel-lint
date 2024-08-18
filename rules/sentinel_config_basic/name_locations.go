package sentinel_config_basic

import (
	"github.com/glennsarti/sentinel-parser/position"
	"github.com/glennsarti/sentinel-parser/sentinel_config/ast"
)

type nameLocations map[string]nameLocationList

type nameLocationList []nameLocation

type nameLocation struct {
	Range position.SourceRange
	Type  string
}

func findNameLocations(f *ast.File) *nameLocations {
	if f == nil {
		return nil
	}

	allNames := make(nameLocations, 0)

	// Imports
	for name, item := range f.Imports {
		if item != nil {
			loc := nameLocation{
				Range: item.Range().Clone(),
				Type:  item.BlockType(),
			}
			if list, ok := allNames[name]; ok {
				allNames[name] = append(list, loc)
			} else {
				allNames[name] = nameLocationList{loc}
			}
		}
	}

	// Globals
	for name, item := range f.Globals {
		if item != nil {
			loc := nameLocation{
				Range: item.Range().Clone(),
				Type:  item.BlockType(),
			}
			if list, ok := allNames[name]; ok {
				allNames[name] = append(list, loc)
			} else {
				allNames[name] = []nameLocation{loc}
			}
		}
	}

	// Params
	for name, item := range f.Params {
		if item != nil {
			loc := nameLocation{
				Range: item.Range().Clone(),
				Type:  item.BlockType(),
			}
			if list, ok := allNames[name]; ok {
				allNames[name] = append(list, loc)
			} else {
				allNames[name] = []nameLocation{loc}
			}
		}
	}

	// Policies
	for name, item := range f.Policies {
		if item != nil {
			loc := nameLocation{
				Range: item.Range().Clone(),
				Type:  item.BlockType(),
			}
			if list, ok := allNames[name]; ok {
				allNames[name] = append(list, loc)
			} else {
				allNames[name] = []nameLocation{loc}
			}
		}
	}

	return &allNames
}
