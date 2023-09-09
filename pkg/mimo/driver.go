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
	"fmt"

	"github.com/rs/zerolog/log"
)

type Driver struct {
	realDataSource DataRowReader
	maskDataSource DataRowReader
	subscribers    Suscribers
	report         Report
}

func NewDriver(
	realReader DataRowReader,
	maskedReader DataRowReader,
	multimapFactory MultimapFactory,
	counterFactory CounterFactory,
	subs ...EventSubscriber,
) Driver {
	return Driver{
		realDataSource: realReader,
		maskDataSource: maskedReader,
		subscribers:    subs,
		report:         NewReport(subs, NewConfig(), multimapFactory, counterFactory),
	}
}

func (d *Driver) Configure(c Config) {
	d.report.config = c
}

func (d *Driver) Analyze() (Report, error) {
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

		d.preprocess(realRow)
		d.preprocess(maskedRow)
		d.report.Update(realRow, maskedRow)
	}

	return d.report, nil
}

func (d Driver) Close() error {
	errors := []error{}

	for key, metric := range d.report.Metrics {
		log.Info().Str("key", key).Msg("close metrics")

		err := metric.Coherence.Close()
		if err != nil {
			errors = append(errors, err)
		}

		err = metric.Identifiant.Close()
		if err != nil {
			errors = append(errors, err)
		}

		err = metric.backend.Close()
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple errors %w (%v)", errors[0], errors[1:])
	}

	return nil
}
