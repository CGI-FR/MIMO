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

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/cockroachdb/pebble"
)

const (
	KeyPrefix   = "K_"
	CountPrefix = "C_"
	MinPrefix   = "M_"
)

type PebbleMultimap struct {
	db      *pebble.DB
	path    string
	tempory bool
}

func PebbleMultimapFactory(path string) (PebbleMultimap, error) {
	var (
		err     error
		tempory bool
	)

	originalPath := path

	if path == "" {
		path, err = os.MkdirTemp("", "mimo-pebble")
		if err != nil {
			return PebbleMultimap{}, fmt.Errorf("%w", err)
		}
	}

	log.Trace().Str("path", path).Msg("open pebble data base")
	//nolint:exhaustruct
	database, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return PebbleMultimap{}, fmt.Errorf("unable to open database %v : %w", path, err)
	}

	if originalPath == "" {
		tempory = true
	}

	return PebbleMultimap{db: database, path: path, tempory: tempory}, nil
}

// Close the database.
func (m PebbleMultimap) Close() error {
	err := m.db.Close()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// remove database if temporary
	if m.tempory {
		log.Info().Str("path", m.path).Msg("Remove database")
		err = os.RemoveAll(m.path)

		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

// Add a key/value pair to the multimap.
func (m PebbleMultimap) Add(key string, value string) {
	var set map[string]int

	item, closer, err := m.db.Get([]byte(KeyPrefix + key))

	if err != nil {
		set = make(map[string]int)
	} else {
		defer closer.Close()
		err = json.Unmarshal(item, &set)
		if err != nil {
			return
		}
	}

	set[value]++

	rawValue, err := json.Marshal(set)
	if err != nil {
		return
	}

	err = m.db.Set([]byte(KeyPrefix+key), rawValue, pebble.NoSync)
	if err != nil {
		return
	}

	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, int64(len(set)))
	err = m.db.Set([]byte(CountPrefix+key), buf[:n], pebble.NoSync)

	if err != nil {
		return
	}
}

// Count the number of values associated to key.
func (m PebbleMultimap) Count(key string) int {
	var count int64

	item, closer, err := m.db.Get([]byte(CountPrefix + key))
	if err != nil {
		count = 0
	} else {
		defer closer.Close()
		count, _ = binary.Varint(item)
	}

	return int(count)
}

// Rate return the percentage of keys that have a count of 1.
func (m PebbleMultimap) Rate() float64 {
	entries := 0
	entriesWithOneValue := 0

	iter, _ := m.db.NewIter(m.prefixIterOptions([]byte(CountPrefix)))
	for iter.First(); iter.Valid(); iter.Next() {
		localCount, _ := binary.Varint(iter.Value())

		if localCount == 1 {
			entriesWithOneValue++
		}
		entries++
	}

	if err := iter.Close(); err != nil {
		log.Fatal().Err(err).Msg("cant't close database iterator")
	}

	return float64(entriesWithOneValue) / float64(entries)
}

// CountMin returns the minimum count of values associated to a key across the map.
func (m PebbleMultimap) CountMin() int {
	minimum := 0

	iter, _ := m.db.NewIter(m.prefixIterOptions([]byte(CountPrefix)))
	for iter.First(); iter.Valid(); iter.Next() {
		localCount, _ := binary.Varint(iter.Value())

		if minimum == 0 || localCount < int64(minimum) {
			minimum = int(localCount)
		}

		if minimum == 1 {
			break
		}
	}

	if err := iter.Close(); err != nil {
		log.Fatal().Err(err).Msg("cant't close database iterator")
	}

	return minimum
}

func (m PebbleMultimap) keyUpperBound(b []byte) []byte {
	end := make([]byte, len(b))
	copy(end, b)

	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end[:i+1]
		}
	}

	return nil // no upper-bound
}

func (m PebbleMultimap) prefixIterOptions(prefix []byte) *pebble.IterOptions {
	//nolint:exhaustruct
	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: m.keyUpperBound(prefix),
	}
}
