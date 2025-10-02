package format

import (
	"io"
	"fmt"
	"log/slog"

	"github.com/protobom/protobom/pkg/writer"
	"github.com/protobom/protobom/pkg/formats"
	"github.com/protobom/protobom/pkg/sbom"
)

var protobomFormatMap = map[SbomType]formats.Format{
	SbomType{SPDX, JSON}: formats.SPDX23JSON,
	SbomType{CYCLONEDX, JSON}: formats.CDX16JSON,
}

func IsSupportedByProtobom(sbomType SbomType) bool {
	_, exists := protobomFormatMap[sbomType]
	return exists
}

func ValidateForProtobom(sbomType SbomType) error {
	if !IsSupportedByProtobom(sbomType) {
		return fmt.Errorf("SbomType %+v is not supported by protobom format", sbomType)
	}
	return nil
}


type ProtobomWriter struct {
  doc *sbom.Document
	sbomType SbomType
}

func NewProtobomWriter(doc *sbom.Document, sbomType SbomType) (*ProtobomWriter, error) {
  if err := ValidateForProtobom(sbomType); err != nil {
		return nil, err
	}
	return &ProtobomWriter{
		doc: doc,
		sbomType: sbomType,
	}, nil
}

func (pw *ProtobomWriter) Write(outputFile io.Writer) error {
	w := writer.New()
	format := protobomFormatMap[pw.sbomType]

	slog.Debug("Output SBOM",
		"path:", outputFile,
		"format:", pw.sbomType,
	)
	return w.WriteStreamWithOptions(
		pw.doc, outputFile, &writer.Options{Format: format},
	)
}

