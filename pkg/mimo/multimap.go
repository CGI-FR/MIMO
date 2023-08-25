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
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Multimap struct {
	Backend MultimapBackend
}

// Close the database.
func (m Multimap) Close() error {
	err := m.Backend.Close()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Add a key/value pair to the multimap.
func (m Multimap) Add(key string, value string) {
	set, err := m.Backend.GetKey(key)
	if errors.Is(ErrKeyNotFound, err) {
		set = make(map[string]int)
	}

	set[value]++

	err = m.Backend.SetKey(key, set)
	if err != nil {
		return
	}
}

// Count the number of values associated to key.
func (m Multimap) Count(key string) int {
	return m.Backend.GetSize(key)
}

// Rate return the percentage of keys that have a count of 1.
func (m Multimap) Rate() float64 {
	entries := 0
	entriesWithOneValue := 0

	iter := m.Backend.NewSizeIterator()

	for iter.First(); iter.Valid(); iter.Next() {
		if iter.Value() == 1 {
			entriesWithOneValue++
		}
		entries++
	}

	if err := iter.Close(); err != nil {
		log.Warn().Err(err).Msg("cant't close database iterator")
	}

	return float64(entriesWithOneValue) / float64(entries)
}

// CountMin returns the minimum count of values associated to a key across the map.
func (m Multimap) CountMin() int {
	minimum := 0

	iter := m.Backend.NewSizeIterator()
	for iter.First(); iter.Valid(); iter.Next() {
		localCount := iter.Value()

		if minimum == 0 || localCount < minimum {
			minimum = localCount
		}

		if minimum == 1 {
			break
		}
	}

	if err := iter.Close(); err != nil {
		log.Warn().Err(err).Msg("cant't close database iterator")
	}

	return minimum
}
