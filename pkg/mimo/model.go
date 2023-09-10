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

	"github.com/ohler55/ojg/jp"
	"github.com/rs/zerolog/log"
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
	Fieldname   string         // Fieldname is name of column analyzed
	Coherence   Multimap       // Coherence is a multimap used to compute the coherence rate
	Identifiant Multimap       // Identifiant is a multimap used to compute the identifiable rate
	Constraints []Constraint   // Constraints is the set of rules to validate
	backend     CounterBackend // Backend counters storage
}

type (
	MultimapFactory func(fieldname string) Multimap
	CounterFactory  func(fieldname string) CounterBackend
)

func NewMetrics(
	fieldname string, multimapFactory MultimapFactory, counterFactory CounterFactory, constraints ...Constraint,
) Metrics {
	return Metrics{
		Fieldname:   fieldname,
		Coherence:   multimapFactory(fieldname + "-coherence"),
		Identifiant: multimapFactory(fieldname + "-identifiant"),
		Constraints: constraints,
		backend:     counterFactory(fieldname + "-counters"),
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

	m.backend.IncreaseTotalCount()

	m.Coherence.Add(toStringSlice(coherenceValue), maskedValueStr)
	m.Identifiant.Add(maskedValueStr, realValueStr)

	log.Trace().
		Str("masked", maskedValueStr).
		Str("real", realValueStr).
		Str("coherence", toStringSlice(coherenceValue)).
		Msg("stored")

	if realValue == nil {
		m.backend.IncreaseNilCount()

		return true
	}

	if isExcluded(config.Exclude, realValue, realValueStr) {
		m.backend.IncreaseIgnoredCount()

		return true
	}

	if realValueOk && maskedValueOk {
		if realValueStr != maskedValueStr {
			m.backend.IncreaseMaskedCount()
		} else if m.backend.GetMaskedCount() == nonBlankCount {
			subs.PostFirstNonMaskedValue(fieldname, realValue)
		}
	}

	return true
}

func (m Metrics) NilCount() int64 {
	return m.backend.GetNilCount()
}

func (m Metrics) IgnoredCount() int64 {
	return m.backend.GetIgnoredCount()
}

func (m Metrics) MaskedCount() int64 {
	return m.backend.GetMaskedCount()
}

// BlankCount is the number of blank (null or ignored) values in real data.
func (m Metrics) BlankCount() int64 {
	return m.backend.GetNilCount() + m.backend.GetIgnoredCount()
}

// NonBlankCount is the number of non-blank (non-null and non-ignored) values in real data.
func (m Metrics) NonBlankCount() int64 {
	return m.backend.GetTotalCount() - m.BlankCount()
}

// NonMaskedCount is the number of non-blank (non-null and non-ignored) values in real data that were not masked.
func (m Metrics) NonMaskedCount() int64 {
	return m.NonBlankCount() - m.backend.GetMaskedCount()
}

// K is the minimum number of value pseudonym was attributed.
func (m Metrics) K() int {
	return m.Identifiant.CountMin()
}

// MaskedRate is equal to
//
//	Number of non-blank real values masked
//	  / (Number of values analyzed - Number of blank (null or ignored) values in real data) ).
func (m Metrics) MaskedRate() float64 {
	return float64(m.backend.GetMaskedCount()) / float64(m.NonBlankCount())
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
		if constraint.Target == CoherentRate {
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
	counterFactory  CounterFactory
}

func NewReport(
	subs []EventSubscriber, config Config, multiMapFactory MultimapFactory, counterFactory CounterFactory,
) Report {
	return Report{make(map[string]Metrics), subs, config, multiMapFactory, counterFactory}
}

func (r Report) UpdateDeep(root DataRow, realRow DataRow, maskedRow DataRow, stack []any, path ...string) {
	for key, realValue := range realRow {
		newpath := append(path, key) //nolint:gocritic

		switch typedRealValue := realValue.(type) {
		case map[string]any:
			if typedMaskedValue, ok := maskedRow[key].(map[string]any); ok {
				r.UpdateDeep(root, typedRealValue, typedMaskedValue, append(stack, realValue), newpath...)
			} else {
				log.Warn().
					Strs("path", newpath).
					Msg("ignored path because structure is different between real and masked data")
			}
		case []any:
			if typedMaskedValue, ok := maskedRow[key].([]any); ok {
				r.UpdateArray(root, typedRealValue, typedMaskedValue, append(stack, realValue), newpath...)
			}
		case nil, any:
			r.UpdateValue(root, typedRealValue, maskedRow[key], append(stack, realValue), newpath...)
		default:
			log.Warn().
				Strs("path", newpath).
				Msg("ignored path because structure is not supported")
		}
	}
}

func (r Report) UpdateArray(root DataRow, realArray []any, maskedArray []any, stack []any, path ...string) {
	for index := 0; index < len(realArray) && index < len(maskedArray); index++ {
		newpath := append(path, "[]") //nolint:gocritic

		switch typedRealItem := realArray[index].(type) {
		case map[string]any:
			if typedMaskedItem, ok := maskedArray[index].(map[string]any); ok {
				r.UpdateDeep(root, typedRealItem, typedMaskedItem, append(stack, realArray[index]), newpath...)
			} else {
				log.Warn().
					Strs("path", newpath).
					Int("index", index).
					Msg("ignored item in array at path because structure is different between real and masked data")
			}
		case []any:
			if typedMaskedItem, ok := maskedArray[index].([]any); ok {
				r.UpdateArray(root, typedRealItem, typedMaskedItem, append(stack, realArray[index]), newpath...)
			}
		case nil, any:
			r.UpdateValue(root, typedRealItem, maskedArray[index], append(stack, realArray[index]), newpath...)
		default:
			log.Warn().
				Strs("path", newpath).
				Int("index", index).
				Msg("ignored item in array at path because structure is not supported")
		}
	}
}

func (r Report) UpdateValue(root DataRow, realValue any, maskedValue any, stack []any, path ...string) {
	key := strings.Join(path, ".")

	config := NewDefaultColumnConfig()
	if cfg, ok := r.config.ColumnConfigs[key]; ok {
		config = cfg
	}

	if len(config.Alias) > 0 {
		key = config.Alias
	}

	metrics, exists := r.Metrics[key]
	if !exists {
		metrics = NewMetrics(key, r.multiMapFactory, r.counterFactory, r.config.ColumnConfigs[key].Constraints...)
		r.subs.PostNewField(key)
	}

	coherenceValues := computeCoherenceValues(config, root, stack)
	if len(coherenceValues) == 0 {
		coherenceValues = []any{realValue}
	}

	if !metrics.Update(key, realValue, maskedValue, coherenceValues, r.subs, config) && !exists {
		metrics.Coherence.Close()
		metrics.Identifiant.Close()
	} else {
		r.Metrics[key] = metrics
	}
}

func computeCoherenceValues(config ColumnConfig, root DataRow, stack []any) []any {
	coherenceValues := make([]any, len(config.CoherentWith))

	for i, coherentColumn := range config.CoherentWith {
		if jpexp, err := jp.ParseString(coherentColumn); err != nil {
			coherenceValues[i] = root[coherentColumn]
		} else {
			coherenceValues[i] = jpexp.Get(root)[0]
		}
	}

	if len(config.CoherentSource) > 0 {
		source, err := generateCoherentSource(config.CoherentSource, root, stack)

		log.Err(err).Str("result", source).Msg("generating coherence source from template")

		if err == nil {
			coherenceValues = append(coherenceValues, source) //nolint:makezero
		}
	}

	return coherenceValues
}

func (r Report) Update(realRow DataRow, maskedRow DataRow) {
	r.UpdateDeep(realRow, realRow, maskedRow, []any{})
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

func isExcluded(exclude []any, value any, valueStr string) bool {
	if slices.Contains(exclude, value) {
		return true
	}

	for _, exVal := range exclude {
		if exValStr, ok := toString(exVal); ok && valueStr == exValStr {
			return true
		}
	}

	return false
}
