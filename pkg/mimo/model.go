package mimo

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type DataRow map[string]any

type Metrics struct {
	TotalCount  int64 // Number of values analyzed
	BlankCount  int64 // Number of blank (null or empty) values in real data
	MaskedCount int64 // Number of non-blank real values masked
}

type Report map[string]*Metrics

func NewReport() Report {
	return Report{}
}

func (r Report) Update(realRow DataRow, maskedRow DataRow) {
	for key, realValue := range realRow {
		metrics := r[key]
		if metrics == nil {
			metrics = NewMetrics()
		}

		metrics.Update(realValue, maskedRow[key])
		r[key] = metrics
	}
}

func (r Report) Print() {
	fmt.Println("Metrics")
	fmt.Println("=======")

	for key, metrics := range r {
		fmt.Println(key, metrics.MaskingRate()*100, "%") //nolint:gomnd
	}
}

func NewMetrics() *Metrics {
	return &Metrics{TotalCount: 0, BlankCount: 0, MaskedCount: 0}
}

func (m *Metrics) Update(realValue any, maskedValue any) {
	log.Trace().Interface("real", realValue).Interface("masked", maskedValue).Msg("update metric")

	m.TotalCount++
	if realValue == nil {
		m.BlankCount++
	} else if realValue != maskedValue {
		m.MaskedCount++
	}
}

// MaskingRate is equal to
//
//	Number of non-blank real values masked
//	  / (Number of values analyzed - Number of blank (null or empty) values in real data) ).
func (m *Metrics) MaskingRate() float64 {
	return float64(m.MaskedCount) / float64(m.TotalCount-m.BlankCount)
}
