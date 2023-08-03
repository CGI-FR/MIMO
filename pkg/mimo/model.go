package mimo

import (
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
	TotalCount  int64 // Number of values analyzed
	BlankCount  int64 // Number of blank (null or empty) values in real data
	MaskedCount int64 // Number of non-blank real values masked
}

func NewMetrics() Metrics {
	return Metrics{TotalCount: 0, BlankCount: 0, MaskedCount: 0}
}

func (m *Metrics) Update(fieldname string, realValue any, maskedValue any, subs Suscribers) {
	log.Trace().Interface("real", realValue).Interface("masked", maskedValue).Msg("update metric")

	switch {
	case realValue == nil:
		m.BlankCount++
	case realValue != maskedValue:
		m.MaskedCount++
	case m.MaskedCount == (m.TotalCount - m.BlankCount):
		subs.PostFirstNonMaskedValue(fieldname, realValue)
	}

	m.TotalCount++
}

// MaskingRate is equal to
//
//	Number of non-blank real values masked
//	  / (Number of values analyzed - Number of blank (null or empty) values in real data) ).
func (m Metrics) MaskingRate() float64 {
	return float64(m.MaskedCount) / float64(m.TotalCount-m.BlankCount)
}

type Report struct {
	metrics map[string]Metrics
	subs    Suscribers
}

func NewReport(subs []EventSubscriber) Report {
	return Report{make(map[string]Metrics), subs}
}

func (r Report) Update(realRow DataRow, maskedRow DataRow) {
	for key, realValue := range realRow {
		metrics, exists := r.metrics[key]
		if !exists {
			metrics = NewMetrics()

			r.subs.PostNewField(key)
		}

		metrics.Update(key, realValue, maskedRow[key], r.subs)
		r.metrics[key] = metrics
	}
}

func (r Report) Columns() []string {
	columns := make([]string, 0, len(r.metrics))
	for colname := range r.metrics {
		columns = append(columns, colname)
	}

	return columns
}

func (r Report) ColumnMetric(colname string) Metrics {
	return r.metrics[colname]
}
