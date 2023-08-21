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

package infra_test

import (
	"log"
	"testing"

	"github.com/cgi-fr/mimo/internal/infra"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestBadgerMultimap(t *testing.T) {
	t.Parallel()

	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	multimap := infra.BadgerMultimap{DB: db}

	multimap.Add("A", "X")
	multimap.Add("A", "Y")
	multimap.Add("B", "Z")
	multimap.Add("B", "Z")
	multimap.Add("C", "Z")
	multimap.Add("C", "Z")
	multimap.Add("D", "W")
	multimap.Add("D", "W")
	multimap.Add("E", "V")
	multimap.Add("E", "V")

	assert.Equal(t, 2, multimap.Count("A"))
	assert.Equal(t, 1, multimap.Count("B"))
	assert.Equal(t, 1, multimap.Count("C"))
	assert.Equal(t, 1, multimap.Count("D"))
	assert.Equal(t, 1, multimap.Count("E"))
	assert.Equal(t, 0, multimap.Count("F"))

	assert.Equal(t, 0.8, multimap.Rate())
}
