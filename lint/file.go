package lint

import (
	"github.com/glennsarti/sentinel-parser/filetypes"
)

type File interface {
	Type() filetypes.FileType
	Path() string
}
