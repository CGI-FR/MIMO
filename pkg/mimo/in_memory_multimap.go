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

type InMemoryMultimapBackend map[string]map[string]int

func (m InMemoryMultimapBackend) SetKey(key string, value map[string]int) error {
	m[key] = value

	return nil
}

// Close the backend.
func (m InMemoryMultimapBackend) Close() error {
	return nil
}

func (m InMemoryMultimapBackend) GetKey(key string) (map[string]int, error) {
	value, ok := m[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return value, nil
}

func (m InMemoryMultimapBackend) GetSize(key string) int {
	return len(m[key])
}

// CountMin returns the minimum count of values associated to a key across the map.
func (m InMemoryMultimapBackend) NewSizeIterator() SizeIterator { //nolint: ireturn
	sizes := []int{}
	for _, array := range m {
		sizes = append(sizes, len(array))
	}

	return &InMemoryIterator{sizes, 0}
}

type InMemoryIterator struct {
	sizes  []int
	cursor int
}

func (i *InMemoryIterator) First() bool {
	return i.cursor == 0
}

func (i *InMemoryIterator) Next() bool {
	i.cursor++

	return i.cursor < len(i.sizes)
}

func (i *InMemoryIterator) Valid() bool {
	return i.cursor < len(i.sizes)
}

func (i *InMemoryIterator) Value() int {
	return i.sizes[i.cursor]
}

func (i *InMemoryIterator) Close() error {
	return nil
}
