package mimo

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Driver struct{}

func NewDriver() Driver {
	return Driver{}
}

func (a Driver) Analyze(realReader DataRowReader, maskedReader DataRowReader, subs ...EventSubscriber) (Report, error) {
	report := NewReport(subs)

	for {
		realRow, err := realReader.ReadDataRow()
		if err != nil {
			return report, fmt.Errorf("%w: %w", ErrReadingDataRow, err)
		}

		maskedRow, err := maskedReader.ReadDataRow()
		if err != nil {
			return report, fmt.Errorf("%w: %w", ErrReadingDataRow, err)
		}

		if realRow == nil && maskedRow == nil {
			break
		}

		if realRow == nil || maskedRow == nil {
			return report, ErrOrphanRow
		}

		log.Trace().Interface("real", realRow).Interface("masked", maskedRow).Msg("read rows")

		report.Update(realRow, maskedRow)
	}

	return report, nil
}
