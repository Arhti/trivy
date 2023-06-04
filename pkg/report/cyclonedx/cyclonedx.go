package cyclonedx

import (
	"io"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"golang.org/x/xerrors"

	rep "github.com/aquasecurity/trivy/pkg/k8s/report"
	"github.com/aquasecurity/trivy/pkg/sbom/cyclonedx"
	"github.com/aquasecurity/trivy/pkg/types"
)

// Writer implements types.Writer
type Writer struct {
	output    io.Writer
	format    cdx.BOMFileFormat
	marshaler *cyclonedx.Marshaler
}

func NewWriter(output io.Writer, appVersion string) Writer {
	return Writer{
		output:    output,
		format:    cdx.BOMFileFormatJSON,
		marshaler: cyclonedx.NewMarshaler(appVersion),
	}
}

// Write writes the results in CycloneDX format
func (w Writer) Write(report types.Report) error {
	bom, err := w.marshaler.Marshal(report)
	if err != nil {
		return xerrors.Errorf("CycloneDX marshal error: %w", err)
	}
	return w.encodeBom(bom)
}

func (w Writer) WriteKbom(report rep.Report) error {
	bom, err := w.marshaler.MarshalKbom(report)
	if err != nil {
		return xerrors.Errorf("CycloneDX marshal error: %w", err)
	}
	return w.encodeBom(bom)
}

func (w Writer) encodeBom(bom *cdx.BOM) error {
	encoder := cdx.NewBOMEncoder(w.output, w.format)
	encoder.SetPretty(true)
	if err := encoder.Encode(bom); err != nil {
		return xerrors.Errorf("failed to encode bom: %w", err)
	}
	return nil
}
