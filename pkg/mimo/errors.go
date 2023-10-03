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

import "errors"

var (
	// ErrReadingDataRow is returned when an error occurs while reading a data row.
	ErrReadingDataRow = errors.New("error while reading datarow")

	// ErrOrphanRow is returned when a original row does not have a masked version, or the other way around.
	ErrOrphanRow = errors.New("error datarow is orphan")

	// ErrKeyNotFound is returned when key is not in multimap.
	ErrKeyNotFound = errors.New("key not found in multimap")

	// ErrDisparityType is returned when real and masked value have different types.
	ErrDisparityType = errors.New("can't validate value because type has changed")

	// ErrAnalyze is returned when a fatal error occurs during analyze.
	ErrAnalyze = errors.New("fatal error")
)
