package mimo

import (
	"fmt"

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

func (m *Metrics) Update(fieldname string, realValue any, maskedValue any, subs Suscribers) {
	log.Trace().Interface("real", realValue).Interface("masked", maskedValue).Msg("update metric")

	switch {
	case realValue == nil:
		m.NilCount++
	case realValue != maskedValue:
		m.MaskedCount++
	case m.MaskedCount == m.NonBlankCount():
		subs.PostFirstNonMaskedValue(fieldname, realValue)
	}

	m.TotalCount++

	m.Coherence.Add(fmt.Sprint(realValue), fmt.Sprint(maskedValue))
	m.Identifiant.Add(fmt.Sprint(maskedValue), fmt.Sprint(realValue))
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
	return m.TotalCount - m.MaskedCount
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
}

func NewReport(subs []EventSubscriber) Report {
	return Report{make(map[string]Metrics), subs}
}

func (r Report) Update(realRow DataRow, maskedRow DataRow) {
	for key, realValue := range realRow {
		metrics, exists := r.Metrics[key]
		if !exists {
			metrics = NewMetrics()

			r.subs.PostNewField(key)
		}

		metrics.Update(key, realValue, maskedRow[key], r.subs)
		r.Metrics[key] = metrics
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
