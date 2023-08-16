// Copyright (C) 2023 CGI France
//
// This file is part of MIMO.
//
// MIMO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// MIMO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with MIMO.  If not, see <http://www.gnu.org/licenses/>.

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
		return fmt.Errorf("%w", err)
	} else if err := e.tmpl.Execute(file, report); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
