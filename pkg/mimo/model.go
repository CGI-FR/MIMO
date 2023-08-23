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
	TotalCount  int64        // TotalCount is the number of values analyzed
	NilCount    int64        // NilCount is the number of null values in real data
	EmptyCount  int64        // EmptyCount is the number of empty values in real data (empty string or numbers at 0 value)
	MaskedCount int64        // MaskedCount is the number of non-blank real values masked
	Coherence   Multimap     // Coherence is a multimap used to compute the coherence rate
	Identifiant Multimap     // Identifiant is a multimap used to compute the identifiable rate
	Constraints []Constraint // Constraints is the set of rules to validate
}

type MultimapFactory func() Multimap

func NewMetrics(multimapFactory MultimapFactory, constraints ...Constraint) Metrics {
	return Metrics{
		TotalCount:  0,
		NilCount:    0,
		EmptyCount:  0,
		MaskedCount: 0,
		Coherence:   multimapFactory(),
		Identifiant: multimapFactory(),
		Constraints: constraints,
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

// MaskedRateValidate returns :
//   - -1 if at least one constraint fail on the MaskedRate,
//   - 0 if no constraint exist on the MaskedRate,
//   - 1 if all constraints succeed on the MaskedRate,
func (m Metrics) MaskedRateValidate() int {
	result := 0

	for _, constraint := range m.Constraints {
		if constraint.Target == MaskingRate {
			if !validate(constraint.Type, constraint.Value, m.MaskedRate()) {
				return -1
			}

			result = 1
		}
	}

	return result
}

// CoherenceRateValidate returns :
//   - -1 if at least one constraint fail on the CoherenceRate,
//   - 0 if no constraint exist on the CoherenceRate,
//   - 1 if all constraints succeed on the CoherenceRate,
func (m Metrics) CoherenceRateValidate() int {
	result := 0

	for _, constraint := range m.Constraints {
		if constraint.Target == CohenrentRate {
			if !validate(constraint.Type, constraint.Value, m.Coherence.Rate()) {
				return -1
			}

			result = 1
		}
	}

	return result
}

// IdentifiantRateValidate returns :
//   - -1 if at least one constraint fail on the IdentifiantRate,
//   - 0 if no constraint exist on the IdentifiantRate,
//   - 1 if all constraints succeed on the IdentifiantRate,
func (m Metrics) IdentifiantRateValidate() int {
	result := 0

	for _, constraint := range m.Constraints {
		if constraint.Target == IdentifiantRate {
			if !validate(constraint.Type, constraint.Value, m.Identifiant.Rate()) {
				return -1
			}

			result = 1
		}
	}

	return result
}

// Validate returns :
//   - -1 if at least one constraint fail,
//   - 0 if no constraint exist,
//   - 1 if all constraints succeed ,
func (m Metrics) Validate() int {
	resultMaskedRate := m.MaskedRateValidate()
	if resultMaskedRate < 0 {
		return -1
	}

	resultCoherentRate := m.CoherenceRateValidate()
	if resultCoherentRate < 0 {
		return -1
	}

	resultIdentifiantRate := m.IdentifiantRateValidate()
	if resultIdentifiantRate < 0 {
		return -1
	}

	if resultMaskedRate > 0 || resultCoherentRate > 0 || resultIdentifiantRate > 0 {
		return 1
	}

	return 0
}

type Report struct {
	Metrics         map[string]Metrics
	subs            Suscribers
	config          Config
	multiMapFactory MultimapFactory
}

func NewReport(subs []EventSubscriber, config Config, multiMapFactory MultimapFactory) Report {
	return Report{make(map[string]Metrics), subs, config, multiMapFactory}
}

func (r Report) Update(realRow DataRow, maskedRow DataRow) {
	for key, realValue := range realRow {
		metrics, exists := r.Metrics[key]
		if !exists {
			metrics = NewMetrics(r.multiMapFactory, r.config.ColumnConfigs[key].Constraints...)

			r.subs.PostNewField(key)
		}

		config := NewDefaultColumnConfig(key)
		if cfg, ok := r.config.ColumnConfigs[key]; ok {
			config = cfg
		}

		coherenceValues := make([]any, len(config.CoherentWith))
		for i, coherentColumn := range config.CoherentWith {
			coherenceValues[i] = realRow[coherentColumn]
		}

		if len(coherenceValues) == 0 {
			coherenceValues = []any{realValue}
		}

		if metrics.Update(key, realValue, maskedRow[key], coherenceValues, r.subs, config) {
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
	result := &strings.Builder{}

	for _, value := range values {
		if str, ok := toString(value); ok {
			result.WriteString(str)
		}

		result.WriteString("_")
	}

	return result.String()
}

func validate(constraint ConstraintType, reference float64, value float64) bool {
	switch constraint {
	case ShouldEqual:
		return value == reference
	case ShouldBeGreaterThan:
		return value > reference
	case ShouldBeGreaterThanOrEqualTo:
		return value >= reference
	case ShouldBeLowerThan:
		return value < reference
	case ShouldBeLessThanOrEqualTo:
		return value <= reference
	default:
		return false
	}
}
