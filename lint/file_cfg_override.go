package lint

import (
	"github.com/glennsarti/sentinel-parser/filetypes"
	scast "github.com/glennsarti/sentinel-parser/sentinel_config/ast"
)

var _ File = ConfigOverrideFile{}

// Configuration Override file
type ConfigOverrideFile struct {
	ConfigFile  *scast.File
	PrimaryFile *scast.File
	FilePath    string
}

func (f ConfigOverrideFile) Type() filetypes.FileType {
	return filetypes.ConfigOverrideFileType
}

func (f ConfigOverrideFile) Path() string {
	return f.FilePath
}
