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

type InMemoryCounterBackend struct {
	TotalCount   int64 // TotalCount is the number of values analyzed
	NilCount     int64 // NilCount is the number of null values in real data
	IgnoredCount int64 // IgnoredCount is the number of ignored values in real data
	MaskedCount  int64 // MaskedCount is the number of non-blank real values masked
}

func (b *InMemoryCounterBackend) IncreaseTotalCount() {
	b.TotalCount++
}

func (b *InMemoryCounterBackend) GetTotalCount() int64 {
	return b.TotalCount
}

func (b *InMemoryCounterBackend) IncreaseNilCount() {
	b.NilCount++
}

func (b *InMemoryCounterBackend) GetNilCount() int64 {
	return b.NilCount
}

func (b *InMemoryCounterBackend) IncreaseIgnoredCount() {
	b.IgnoredCount++
}

func (b *InMemoryCounterBackend) GetIgnoredCount() int64 {
	return b.IgnoredCount
}

func (b *InMemoryCounterBackend) IncreaseMaskedCount() {
	b.MaskedCount++
}

func (b *InMemoryCounterBackend) GetMaskedCount() int64 {
	return b.MaskedCount
}

func (b *InMemoryCounterBackend) Close() error {
	return nil
}
