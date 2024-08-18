package lint

import (
	"github.com/glennsarti/sentinel-parser/filetypes"
	sast "github.com/glennsarti/sentinel-parser/sentinel/ast"
)

var _ File = ModuleFile{}

type ModuleFile struct {
	File     *sast.File
	FilePath string
}

func (f ModuleFile) Type() filetypes.FileType {
	return filetypes.ModuleFileType
}

func (f ModuleFile) Path() string {
	return f.FilePath
}
