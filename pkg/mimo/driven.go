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

type DataRowReader interface {
	ReadDataRow() (DataRow, error)
}

type EventSubscriber interface {
	NewField(fieldname string)
	FirstNonMaskedValue(fieldname string, value any)
}

type Multimap interface {
	Add(key string, value string)
	Count(key string) int
	Rate() float64
	CountMin() int
}
