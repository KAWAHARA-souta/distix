package format

import (
	"github.com/protobom/protobom/pkg/sbom"
)

type Writer interface {
	Write(*sbom.Document, SbomFileFormatType) error
}
