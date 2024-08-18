package lint

import (
	"github.com/glennsarti/sentinel-parser/filetypes"
	sast "github.com/glennsarti/sentinel-parser/sentinel/ast"
	scast "github.com/glennsarti/sentinel-parser/sentinel_config/ast"
)

var _ File = PolicyFile{}

type PolicyFile struct {
	File       *sast.File
	ConfigFile *scast.File
	FilePath   string
}

func (f PolicyFile) Type() filetypes.FileType {
	return filetypes.PolicyFileType
}

func (f PolicyFile) Path() string {
	return f.FilePath
}
