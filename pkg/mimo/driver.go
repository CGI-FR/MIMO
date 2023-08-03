package mimo

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Driver struct {
	realDataSource DataRowReader
	maskDataSource DataRowReader
	subscribers    Suscribers
	report         Report
}

func NewDriver(realReader DataRowReader, maskedReader DataRowReader, subs ...EventSubscriber) Driver {
	return Driver{
		realDataSource: realReader,
		maskDataSource: maskedReader,
		subscribers:    subs,
		report:         NewReport(subs),
	}
}

func (d Driver) Analyze() (Report, error) {
	for {
		realRow, err := d.realDataSource.ReadDataRow()
		if err != nil {
			return d.report, fmt.Errorf("%w: %w", ErrReadingDataRow, err)
		}

		maskedRow, err := d.maskDataSource.ReadDataRow()
		if err != nil {
			return d.report, fmt.Errorf("%w: %w", ErrReadingDataRow, err)
		}

		if realRow == nil && maskedRow == nil {
			break
		}

		if realRow == nil || maskedRow == nil {
			return d.report, ErrOrphanRow
		}

		log.Trace().Interface("real", realRow).Interface("masked", maskedRow).Msg("read rows")

		d.report.Update(realRow, maskedRow)
	}

	return d.report, nil
}
