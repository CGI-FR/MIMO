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

package mimo_test

import (
	"testing"

	"github.com/cgi-fr/mimo/internal/infra"
	"github.com/cgi-fr/mimo/pkg/mimo"
	"github.com/rs/zerolog"
)

func BenchmarkInMemory(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	realReader, err := infra.NewDataRowReaderJSONLineFromFile("testdata/real.jsonl")
	if err != nil {
		b.FailNow()
	}

	maskedReader, err := infra.NewDataRowReaderJSONLineFromFile("testdata/masked.jsonl")
	if err != nil {
		b.FailNow()
	}

	driver := mimo.NewDriver(
		realReader,
		maskedReader,
		func(fieldname string) mimo.Multimap {
			return mimo.Multimap{Backend: mimo.InMemoryMultimapBackend{}}
		},
		infra.SubscriberLogger{},
	)

	defer driver.Close()

	for n := 0; n < b.N; n++ {
		if _, err := driver.Analyze(); err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkOnDisk(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	realReader, err := infra.NewDataRowReaderJSONLineFromFile("testdata/real.jsonl")
	if err != nil {
		b.FailNow()
	}

	maskedReader, err := infra.NewDataRowReaderJSONLineFromFile("testdata/masked.jsonl")
	if err != nil {
		b.FailNow()
	}

	driver := mimo.NewDriver(
		realReader,
		maskedReader,
		func(fieldname string) mimo.Multimap {
			factory, err := infra.PebbleMultimapFactory("")
			if err != nil {
				b.FailNow()
			}

			return factory
		},
		infra.SubscriberLogger{},
	)

	defer driver.Close()

	for n := 0; n < b.N; n++ {
		if _, err := driver.Analyze(); err != nil {
			b.FailNow()
		}
	}
}
