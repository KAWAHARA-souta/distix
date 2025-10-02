package data

import (
	"github.com/protobom/protobom/pkg/sbom"
)


type Extractor interface {
	Extract() (SbomData, error)
}


type SbomData interface {
	Convert2ProtobomDocument() (*sbom.Document, error)
	Convert2MultiProtobomDocument() (*sbom.Document, []*sbom.Document, error)
}

