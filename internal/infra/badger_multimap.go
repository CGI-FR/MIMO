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

	"github.com/dgraph-io/badger"
)

const (
	KeyPrefix   = "K_"
	CountPrefix = "C_"
	MinPrefix   = "M_"
)

type BadgerMultimap struct {
	DB *badger.DB
}

// Add a key/value pair to the multimap.
func (m BadgerMultimap) Add(key string, value string) {
	var set map[string]int

	(*m.DB).Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(KeyPrefix + key))
		if err != nil {
			set = make(map[string]int)
		} else {
			err = item.Value(func(val []byte) error {
				err = json.Unmarshal(val, &set)
				if err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				return err
			}
		}

		set[value]++

		value, err := json.Marshal(set)
		if err != nil {
			return err
		}

		err = txn.Set([]byte(KeyPrefix+key), value)
		if err != nil {
			return err
		}

		buf := make([]byte, binary.MaxVarintLen64)
		n := binary.PutVarint(buf, int64(len(set)))
		err = txn.Set([]byte(CountPrefix+key), buf[:n])
		if err != nil {
			return err
		}
		return nil
	})
}

// Count the number of values associated to key.
func (m BadgerMultimap) Count(key string) int {
	var count int64
	(*m.DB).View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(CountPrefix + key))
		if err != nil {
			count = 0
		} else {
			err = item.Value(func(val []byte) error {
				count, _ = binary.Varint(val)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	})

	return int(count)
}

// Rate return the percentage of keys that have a count of 1.
func (m BadgerMultimap) Rate() float64 {
	m.DB.View(func(txn *badger.Txn) error {
		return nil
	})
	return 0
}

// CountMin returns the minimum count of values associated to a key across the map.
func (m BadgerMultimap) CountMin() int {
	return 1
}
