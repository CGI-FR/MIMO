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

type MultimapBackend interface {
	Close() error
	GetKey(key string) (map[string]int, error)
	SetKey(key string, value map[string]int) error
	GetSize(key string) int
	NewSizeIterator() SizeIterator
	GetSamplesMono(n int) []Sample
	GetSamplesMulti(n int) []Sample
}

type SizeIterator interface {
	First() bool
	Next() bool
	Valid() bool
	Value() int
	Close() error
}

type CounterBackend interface {
	IncreaseTotalCount()
	GetTotalCount() int64

	IncreaseNilCount()
	GetNilCount() int64

	IncreaseIgnoredCount()
	GetIgnoredCount() int64

	IncreaseMaskedCount()
	GetMaskedCount() int64

	Close() error
}
