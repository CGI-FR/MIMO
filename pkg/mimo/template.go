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

package mimo

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Template struct {
	tpml *template.Template
}

func NewTemplate(tmplstr string) (*Template, error) {
	if len(tmplstr) > 0 {
		tmpl, err := template.New("").Funcs(generateFuncMap()).Funcs(sprig.TxtFuncMap()).Parse(tmplstr)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return &Template{tmpl}, nil
	}

	return nil, nil //nolint:nilnil
}

func (t *Template) Execute(root DataRow, stack []any) (string, error) {
	funcstack := template.FuncMap{
		"Stack": generateStackFunc(stack),
	}

	tmpl := t.tpml.Funcs(funcstack)

	result := &strings.Builder{}
	err := tmpl.Execute(result, root)

	return result.String(), err
}

func generateFuncMap() template.FuncMap {
	return template.FuncMap{
		"Stack":    func(index int) any { return nil },
		"ToUpper":  strings.ToUpper,
		"ToLower":  strings.ToLower,
		"NoAccent": rmAcc,
	}
}

func generateStackFunc(theStack []any) func(index int) any {
	return func(index int) any {
		if index > 0 {
			return theStack[index-1]
		}

		return theStack[len(theStack)+index-1]
	}
}

// rmAcc removes accents from string
// Function derived from: http://blog.golang.org/normalization
func rmAcc(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)

	return result
}
