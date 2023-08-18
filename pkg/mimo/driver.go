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
	config         Config
}

func NewDriver(realReader DataRowReader, maskedReader DataRowReader, subs ...EventSubscriber) Driver {
	return Driver{
		realDataSource: realReader,
		maskDataSource: maskedReader,
		subscribers:    subs,
		report:         NewReport(subs),
		config:         NewConfig(),
	}
}

func (d *Driver) Configure(c Config) {
	d.config = c
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

		d.report.Update(realRow, maskedRow)
	}

	return d.report, nil
}
