package lint

import (
	"github.com/glennsarti/sentinel-parser/filetypes"
	scast "github.com/glennsarti/sentinel-parser/sentinel_config/ast"
)

var _ File = ConfigPrimaryFile{}

// Primary Configuration file
type ConfigPrimaryFile struct {
	ConfigFile         *scast.File
	ResolvedConfigFile *scast.File
	FilePath           string
}

func (f ConfigPrimaryFile) Type() filetypes.FileType {
	return filetypes.ConfigPrimaryFileType
}

func (f ConfigPrimaryFile) Path() string {
	return f.FilePath
}
