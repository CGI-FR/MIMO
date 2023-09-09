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
	"errors"
	"fmt"
	"os"

	"github.com/cgi-fr/mimo/pkg/mimo"
	"github.com/cockroachdb/pebble"
	"github.com/rs/zerolog/log"
)

var (
	TotalKey   = []byte("TOTAL")   //nolint:gochecknoglobals
	NilKey     = []byte("NIL")     //nolint:gochecknoglobals
	IgnoredKey = []byte("IGNORED") //nolint:gochecknoglobals
	MaskedKey  = []byte("MASKED")  //nolint:gochecknoglobals
)

type PebbleCounterBackend struct {
	db      *pebble.DB
	path    string
	tempory bool
}

func (b PebbleCounterBackend) Close() error {
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

func PebbleCounterBackendFactory(path string) (mimo.CounterBackend, error) { //nolint:ireturn
	var (
		err     error
		tempory bool
	)

	originalPath := path

	if path == "" {
		path, err = os.MkdirTemp("", "mimo-pebble")
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	log.Trace().Str("path", path).Msg("open pebble data base")
	//nolint:exhaustruct
	database, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, fmt.Errorf("unable to open database %v : %w", path, err)
	}

	if originalPath == "" {
		tempory = true
	}

	return PebbleCounterBackend{db: database, path: path, tempory: tempory}, nil
}

func (b PebbleCounterBackend) IncreaseTotalCount() {
	var count int64

	item, closer, err := b.db.Get(TotalKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		panic(err)
	}

	if closer != nil {
		defer closer.Close()
	}

	if !errors.Is(err, pebble.ErrNotFound) {
		count = int64(binary.BigEndian.Uint64(item))
	}

	count++

	countbin := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(countbin, uint64(count))

	err = b.db.Set(TotalKey, countbin, pebble.NoSync)
	if err != nil {
		log.Error().AnErr("error", err).Int64("total-count", count).Msg("failed to store counter in database")
	}
}

func (b PebbleCounterBackend) GetTotalCount() int64 {
	var count int64

	item, closer, err := b.db.Get(TotalKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		panic(err)
	}

	if closer != nil {
		defer closer.Close()
	}

	if !errors.Is(err, pebble.ErrNotFound) {
		count = int64(binary.BigEndian.Uint64(item))
	}

	return count
}

func (b PebbleCounterBackend) IncreaseNilCount() {
	var count int64

	item, closer, err := b.db.Get(NilKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		panic(err)
	}

	if closer != nil {
		defer closer.Close()
	}

	if !errors.Is(err, pebble.ErrNotFound) {
		count = int64(binary.BigEndian.Uint64(item))
	}

	count++

	countbin := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(countbin, uint64(count))

	err = b.db.Set(NilKey, countbin, pebble.NoSync)
	if err != nil {
		log.Error().AnErr("error", err).Int64("nil-count", count).Msg("failed to store counter in database")
	}
}

func (b PebbleCounterBackend) GetNilCount() int64 {
	var count int64

	item, closer, err := b.db.Get(NilKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		panic(err)
	}

	if closer != nil {
		defer closer.Close()
	}

	if !errors.Is(err, pebble.ErrNotFound) {
		count = int64(binary.BigEndian.Uint64(item))
	}

	return count
}

func (b PebbleCounterBackend) IncreaseIgnoredCount() {
	var count int64

	item, closer, err := b.db.Get(IgnoredKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		panic(err)
	}

	if closer != nil {
		defer closer.Close()
	}

	if !errors.Is(err, pebble.ErrNotFound) {
		count = int64(binary.BigEndian.Uint64(item))
	}

	count++

	countbin := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(countbin, uint64(count))

	err = b.db.Set(IgnoredKey, countbin, pebble.NoSync)
	if err != nil {
		log.Error().AnErr("error", err).Int64("ignored-count", count).Msg("failed to store counter in database")
	}
}

func (b PebbleCounterBackend) GetIgnoredCount() int64 {
	var count int64

	item, closer, err := b.db.Get(IgnoredKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		panic(err)
	}

	if closer != nil {
		defer closer.Close()
	}

	if !errors.Is(err, pebble.ErrNotFound) {
		count = int64(binary.BigEndian.Uint64(item))
	}

	return count
}

func (b PebbleCounterBackend) IncreaseMaskedCount() {
	var count int64

	item, closer, err := b.db.Get(MaskedKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		panic(err)
	}

	if closer != nil {
		defer closer.Close()
	}

	if !errors.Is(err, pebble.ErrNotFound) {
		count = int64(binary.BigEndian.Uint64(item))
	}

	count++

	countbin := make([]byte, 8) //nolint:gomnd
	binary.BigEndian.PutUint64(countbin, uint64(count))

	err = b.db.Set(MaskedKey, countbin, pebble.NoSync)
	if err != nil {
		log.Error().AnErr("error", err).Int64("masked-count", count).Msg("failed to store counter in database")
	}
}

func (b PebbleCounterBackend) GetMaskedCount() int64 {
	var count int64

	item, closer, err := b.db.Get(MaskedKey)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		panic(err)
	}

	if closer != nil {
		defer closer.Close()
	}

	if !errors.Is(err, pebble.ErrNotFound) {
		count = int64(binary.BigEndian.Uint64(item))
	}

	return count
}
