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
	"io"
	"os"
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

	file, err := os.Open("testdata/masked.jsonl")
	if err != nil {
		b.FailNow()
	}

	maskedReader := infra.NewDataRowReaderJSONLine(file, io.Discard)

	driver := mimo.NewDriver(
		realReader,
		maskedReader,
		func(fieldname string) mimo.Multimap {
			return mimo.Multimap{Backend: mimo.InMemoryMultimapBackend{}}
		},
		func(fieldname string) mimo.CounterBackend {
			return &mimo.InMemoryCounterBackend{
				TotalCount:   0,
				NilCount:     0,
				IgnoredCount: 0,
				MaskedCount:  0,
			}
		},
		infra.NewSubscriberLogger(),
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

	file, err := os.Open("testdata/masked.jsonl")
	if err != nil {
		b.FailNow()
	}

	maskedReader := infra.NewDataRowReaderJSONLine(file, io.Discard)

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
		func(fieldname string) mimo.CounterBackend {
			factory, err := infra.PebbleCounterBackendFactory("")
			if err != nil {
				b.FailNow()
			}

			return factory
		},
		infra.NewSubscriberLogger(),
	)

	defer driver.Close()

	for n := 0; n < b.N; n++ {
		if _, err := driver.Analyze(); err != nil {
			b.FailNow()
		}
	}
}

//nolint:funlen
func BenchmarkAllOptions(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)

	realReader, err := infra.NewDataRowReaderJSONLineFromFile("testdata/single-100-1.jsonl")
	if err != nil {
		b.FailNow()
	}

	file, err := os.Open("testdata/single-100-2.jsonl")
	if err != nil {
		b.FailNow()
	}

	maskedReader := infra.NewDataRowReaderJSONLine(file, io.Discard)

	driver := mimo.NewDriver(
		realReader,
		maskedReader,
		func(fieldname string) mimo.Multimap {
			return mimo.Multimap{Backend: mimo.InMemoryMultimapBackend{}}
		},
		func(fieldname string) mimo.CounterBackend {
			return &mimo.InMemoryCounterBackend{
				TotalCount:   0,
				NilCount:     0,
				IgnoredCount: 0,
				MaskedCount:  0,
			}
		},
		infra.NewSubscriberLogger(),
	)

	defer driver.Close()

	driver.Configure(mimo.Config{
		ColumnNames: []string{"value"},
		ColumnConfigs: map[string]mimo.ColumnConfig{
			"value": {
				Exclude:         []any{"Odile", "Tiffany"},
				ExcludeTemplate: `{{uuidv4 | contains "a"}}`,
				CoherentWith:    []string{"name", "surname"},
				CoherentSource:  "{{.name | NoAccent | title}} {{.surname | NoAccent | upper}}",
				Constraints: []mimo.Constraint{
					{
						Target: mimo.MaskingRate,
						Type:   mimo.ShouldBeGreaterThan,
						Value:  .9,
					},
				},
				Alias: "",
			},
		},
		PreprocessConfigs: []mimo.PreprocessConfig{
			{
				Path:  "email",
				Value: "{{.name | NoAccent | lower}}.{{.surname | NoAccent | lower}}@{{uuidv4}}.com",
			},
		},
	})

	for n := 0; n < b.N; n++ {
		if _, err := driver.Analyze(); err != nil {
			b.FailNow()
		}
	}
}
