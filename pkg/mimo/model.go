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
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type DataRow map[string]any

type Suscribers []EventSubscriber

func (subs Suscribers) PostNewField(fieldname string) {
	for _, sub := range subs {
		sub.NewField(fieldname)
	}
}

func (subs Suscribers) PostFirstNonMaskedValue(fieldname string, value any) {
	for _, sub := range subs {
		sub.FirstNonMaskedValue(fieldname, value)
	}
}

type Metrics struct {
	TotalCount  int64    // TotalCount is the number of values analyzed
	NilCount    int64    // NilCount is the number of null values in real data
	EmptyCount  int64    // EmptyCount is the number of empty values in real data (empty string or numbers at 0 value)
	MaskedCount int64    // MaskedCount is the number of non-blank real values masked
	Coherence   Multimap // Coherence is a multimap used to compute the coherence rate
	Identifiant Multimap // Identifiant is a multimap used to compute the identifiable rate
}

func NewMetrics() Metrics {
	return Metrics{
		TotalCount:  0,
		NilCount:    0,
		EmptyCount:  0,
		MaskedCount: 0,
		Coherence:   Multimap{},
		Identifiant: Multimap{},
	}
}

func (m *Metrics) Update(
	fieldname string,
	realValue any,
	maskedValue any,
	coherenceValue []any,
	subs Suscribers,
	config ColumnConfig,
) bool {
	nonBlankCount := m.NonBlankCount()

	realValueStr, realValueOk := toString(realValue)
	maskedValueStr, maskedValueOk := toString(maskedValue)

	if !realValueOk || !maskedValueOk {
		return false // special case (arrays, objects) are not covered right now
	}

	m.TotalCount++

	m.Coherence.Add(toStringSlice(coherenceValue), maskedValueStr)
	m.Identifiant.Add(maskedValueStr, realValueStr)

	if realValue == nil {
		m.NilCount++

		return true
	}

	if slices.Contains(config.Exclude, realValue) {
		m.EmptyCount++

		return true
	}

	if realValueOk && maskedValueOk {
		if realValueStr != maskedValueStr {
			m.MaskedCount++
		} else if m.MaskedCount == nonBlankCount {
			subs.PostFirstNonMaskedValue(fieldname, realValue)
		}
	}

	return true
}

// BlankCount is the number of blank (null or empty) values in real data.
func (m Metrics) BlankCount() int64 {
	return m.NilCount + m.EmptyCount
}

// NonBlankCount is the number of non-blank (non-null and non-empty) values in real data.
func (m Metrics) NonBlankCount() int64 {
	return m.TotalCount - m.BlankCount()
}

// NonMaskedCount is the number of non-blank (non-null and non-empty) values in real data that were not masked.
func (m Metrics) NonMaskedCount() int64 {
	return m.NonBlankCount() - m.MaskedCount
}

// K is the minimum number of value pseudonym was attributed.
func (m Metrics) K() int {
	return m.Identifiant.CountMin()
}

// MaskedRate is equal to
//
//	Number of non-blank real values masked
//	  / (Number of values analyzed - Number of blank (null or empty) values in real data) ).
func (m Metrics) MaskedRate() float64 {
	return float64(m.MaskedCount) / float64(m.NonBlankCount())
}

type Report struct {
	Metrics map[string]Metrics
	subs    Suscribers
	config  Config
}

func NewReport(subs []EventSubscriber, config Config) Report {
	return Report{make(map[string]Metrics), subs, config}
}

func (r Report) Update(realRow DataRow, maskedRow DataRow) {
	for key, realValue := range realRow {
		metrics, exists := r.Metrics[key]
		if !exists {
			metrics = NewMetrics()

			r.subs.PostNewField(key)
		}

		config := NewDefaultColumnConfig(key)
		if cfg, ok := r.config.ColumnConfigs[key]; ok {
			config = cfg
		}

		coherenceValue := make([]any, len(config.CoherentWith))
		for i, coherentColumn := range config.CoherentWith {
			coherenceValue[i] = realRow[coherentColumn]
		}

		if metrics.Update(key, realValue, maskedRow[key], coherenceValue, r.subs, config) {
			r.Metrics[key] = metrics
		}
	}
}

func (r Report) Columns() []string {
	columns := make([]string, 0, len(r.Metrics))
	for colname := range r.Metrics {
		columns = append(columns, colname)
	}

	return columns
}

func (r Report) ColumnMetric(colname string) Metrics {
	return r.Metrics[colname]
}

func toString(value any) (string, bool) {
	var str string
	switch tvalue := value.(type) {
	case string:
		str = strconv.Quote(tvalue)
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float32, float64, bool:
		str = fmt.Sprint(tvalue)
	case json.Number:
		str = string(tvalue)
	case nil:
		str = "nil"
	default:
		return "", false
	}

	return str, true
}

func toStringSlice(values []any) string {
	result := strings.Builder{}

	for _, value := range values {
		if str, ok := toString(value); ok {
			result.WriteString(str)
		}

		result.WriteString("_")
	}

	return result.String()
}
