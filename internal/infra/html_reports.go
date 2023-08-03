package infra

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"

	"github.com/Masterminds/sprig/v3"
	"github.com/adrienaury/mimo/pkg/mimo"
)

//go:embed template/default.html
var defaultTemplate string

type ReportExporter struct {
	tmpl *template.Template
}

func NewReportExporter() ReportExporter {
	t, err := template.New("MIMO Report").Funcs(sprig.FuncMap()).Parse(defaultTemplate)
	if err != nil {
		panic(err)
	}

	return ReportExporter{tmpl: t}
}

func (e ReportExporter) Export(report mimo.Report, filename string) error {
	if file, err := os.Create(filename); err != nil {
		return err
	} else if err := e.tmpl.Execute(file, report); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
