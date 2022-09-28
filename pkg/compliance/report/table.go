package report

import (
	"io"
	"sync"

	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	pkgReport "github.com/aquasecurity/trivy/pkg/report/table"
	"github.com/aquasecurity/trivy/pkg/types"
	"golang.org/x/xerrors"
)

type TableWriter struct {
	Report        string
	Output        io.Writer
	Severities    []dbTypes.Severity
	ColumnHeading []string
}

const (
	ControlIDColumn   = "ID"
	SeverityColumn    = "Severity"
	ControlNameColumn = "Control Name"
	ComplianceColumn  = "Compliance"
)

func (tw TableWriter) Columns() []string {
	return []string{ControlIDColumn, SeverityColumn, ControlNameColumn, ComplianceColumn}
}

func (tw TableWriter) Write(report *ComplianceReport) error {
	switch tw.Report {
	case allReport:
		t := pkgReport.Writer{Output: tw.Output, Severities: tw.Severities, ShowMessageOnce: &sync.Once{}}
		for _, cr := range report.ControlResults {
			r := types.Report{Results: cr.Checks}
			err := t.Write(r)
			if err != nil {
				return err
			}
		}
	case summaryReport:
		writer := NewSummaryWriter(tw.Output, tw.Severities, tw.Columns())
		return writer.Write(report)
	default:
		return xerrors.Errorf(`report %q not supported. Use "summary" or "all"`, tw.Report)
	}

	return nil
}
