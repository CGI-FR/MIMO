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
	"errors"
	"fmt"
	"os"

	"github.com/cgi-fr/mimo/pkg/mimo"
	"github.com/rs/zerolog/log"

	"github.com/cockroachdb/pebble"
)

const (
	KeyPrefix   = "K_"
	CountPrefix = "C_"
	MinPrefix   = "M_"
)

type PebbleMultimapBackend struct {
	db      *pebble.DB
	path    string
	tempory bool
}

func (b PebbleMultimapBackend) Close() error {
	err := b.db.Close()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	log.Info().Str("path", b.path).Msg("database closed")

	// remove database if temporary
	if b.tempory {
		log.Info().Str("path", b.path).Msg("Remove database")
		err = os.RemoveAll(b.path)

		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func PebbleMultimapFactory(path string) (mimo.Multimap, error) {
	var (
		err     error
		tempory bool
	)

	originalPath := path

	if path == "" {
		path, err = os.MkdirTemp("", "mimo-pebble")
		if err != nil {
			return mimo.Multimap{}, fmt.Errorf("%w", err)
		}
	}

	log.Trace().Str("path", path).Msg("open pebble data base")
	//nolint:exhaustruct
	database, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return mimo.Multimap{}, fmt.Errorf("unable to open database %v : %w", path, err)
	}

	if originalPath == "" {
		tempory = true
	}

	return mimo.Multimap{Backend: PebbleMultimapBackend{db: database, path: path, tempory: tempory}}, nil
}

func (b PebbleMultimapBackend) GetKey(key string) (map[string]int, error) {
	var (
		set map[string]int
		err error
	)

	item, closer, err := b.db.Get([]byte(KeyPrefix + key))
	if errors.Is(err, pebble.ErrNotFound) {
		return nil, mimo.ErrKeyNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer closer.Close()

	err = json.Unmarshal(item, &set)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return set, nil
}

func (b PebbleMultimapBackend) SetKey(key string, value map[string]int) error {
	rawValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = b.db.Set([]byte(KeyPrefix+key), rawValue, pebble.NoSync)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// save size in a different namespace
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, int64(len(value)))

	err = b.db.Set([]byte(CountPrefix+key), buf[:n], pebble.NoSync)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (b PebbleMultimapBackend) GetSize(key string) int {
	var count int64

	item, closer, err := b.db.Get([]byte(CountPrefix + key))
	if err != nil {
		count = 0
	} else {
		defer closer.Close()
		count, _ = binary.Varint(item)
	}

	return int(count)
}

func (b PebbleMultimapBackend) NewSizeIterator() mimo.SizeIterator { //nolint: ireturn
	iter, _ := b.db.NewIter(b.prefixIterOptions([]byte(CountPrefix)))

	return PebbleSizeIterator{iter}
}

type PebbleSizeIterator struct {
	iterator *pebble.Iterator
}

func (i PebbleSizeIterator) First() bool {
	return i.iterator.First()
}

func (i PebbleSizeIterator) Next() bool {
	return i.iterator.Next()
}

func (i PebbleSizeIterator) Valid() bool {
	return i.iterator.Valid()
}

func (i PebbleSizeIterator) Value() int {
	localCount, _ := binary.Varint(i.iterator.Value())

	return int(localCount)
}

func (i PebbleSizeIterator) Close() error {
	err := i.iterator.Close()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (b PebbleMultimapBackend) keyUpperBound(upper []byte) []byte {
	end := make([]byte, len(upper))
	copy(end, upper)

	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end[:i+1]
		}
	}

	return nil // no upper-bound
}

func (b PebbleMultimapBackend) prefixIterOptions(prefix []byte) *pebble.IterOptions {
	//nolint:exhaustruct
	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: b.keyUpperBound(prefix),
	}
}
