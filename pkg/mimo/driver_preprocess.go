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
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
)

func (d *Driver) preprocess(row DataRow) {
	for _, preprocess := range d.report.config.PreprocessConfigs {
		preprocessValue(row, strings.Split(preprocess.Path, "."), []any{}, preprocess.Value, row)
	}
}

func preprocessValue(value any, paths []string, stack []any, tmpl *template.Template, root DataRow) {
	path := paths[0]

	var err error

	if len(paths) == 1 {
		if obj, ok := value.(map[string]any); ok {
			obj[path], err = applyTemplate(tmpl, root, append(stack, obj))
		}

		if obj, ok := value.(DataRow); ok {
			obj[path], err = applyTemplate(tmpl, root, append(stack, obj))
		}

		if err != nil {
			log.Error().AnErr("error", err).Msg("failed to generate preprocessed field")
		}

		return
	}

	if path == "[]" {
		if array, ok := value.([]any); ok {
			for _, item := range array {
				preprocessValue(item, paths[1:], append(stack, array), tmpl, root) //nolint:asasalint
			}
		}

		return
	}

	if obj, ok := value.(map[string]any); ok {
		preprocessValue(obj[path], paths[1:], append(stack, obj), tmpl, root)
	}

	if obj, ok := value.(DataRow); ok {
		preprocessValue(obj[path], paths[1:], append(stack, obj), tmpl, root)
	}
}
