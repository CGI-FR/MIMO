package mimo

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Driver struct{}

func NewDriver() Driver {
	return Driver{}
}

func (a Driver) Analyze(realRowReader DataRowReader, maskedRowReader DataRowReader) (Report, error) {
	report := NewReport()

	for {
		realRow, err := realRowReader.ReadDataRow()
		if err != nil {
			return report, fmt.Errorf("%w: %w", ErrReadingDataRow, err)
		}

		maskedRow, err := maskedRowReader.ReadDataRow()
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
