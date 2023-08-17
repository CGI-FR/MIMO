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

package infra

type InMemoryMultimap map[string]map[string]int

// Add a key/value pair to the multimap.
func (m InMemoryMultimap) Add(key string, value string) {
	set, exists := m[key]
	if !exists {
		set = make(map[string]int)
	}

	set[value]++

	m[key] = set
}

// Count the number of values associated to key.
func (m InMemoryMultimap) Count(key string) int {
	set, exists := m[key]
	if !exists {
		return 0
	}

	return len(set)
}

// Rate return the percentage of keys that have a count of 1.
func (m InMemoryMultimap) Rate() float64 {
	cnt := 0

	for _, set := range m {
		if len(set) == 1 {
			cnt++
		}
	}

	return float64(cnt) / float64(len(m))
}

// CountMin returns the minimum count of values associated to a key across the map.
func (m InMemoryMultimap) CountMin() int {
	var cnt int

	for _, set := range m {
		if cnt == 0 || len(set) < cnt {
			cnt = len(set)
		}

		if cnt == 1 {
			break
		}
	}

	return cnt
}
