package lint

import (
	"github.com/glennsarti/sentinel-parser/filetypes"
	scast "github.com/glennsarti/sentinel-parser/sentinel_config/ast"
)

var _ File = ConfigTestFile{}

// Configuration Override file
type ConfigTestFile struct {
	ConfigFile *scast.File
	FilePath   string
}

func (f ConfigTestFile) Type() filetypes.FileType {
	return filetypes.ConfigTestFileType
}

func (f ConfigTestFile) Path() string {
	return f.FilePath
}
