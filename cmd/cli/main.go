// MIT License
//
// Copyright (c) 2021 Adrien Aury
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"os"

	"github.com/adrienaury/mimo/pkg/mimo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Provisioned by ldflags.
var (
	name      string //nolint: gochecknoglobals
	version   string //nolint: gochecknoglobals
	commit    string //nolint: gochecknoglobals
	buildDate string //nolint: gochecknoglobals
	builtBy   string //nolint: gochecknoglobals
)

//nolint:gomnd
func main() {
	//nolint: exhaustruct
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msgf("%v %v (commit=%v date=%v by=%v)", name, version, commit, buildDate, builtBy)

	realData := &Test{
		data: []mimo.DataRow{
			{"name": "Adrien", "age": 12},
			{"name": "Youen", "age": 12},
			{"name": nil, "age": 12},
		},
		index: 0,
	}

	maskedData := &Test{
		data: []mimo.DataRow{
			{"name": "Charles", "age": 12},
			{"name": "Youen", "age": 50},
			{"name": nil, "age": 12},
		},
		index: 0,
	}

	driver := mimo.NewDriver()
	if report, err := driver.Analyze(realData, maskedData); err != nil {
		log.Fatal().AnErr("error", err).Msg("end of program")
	} else {
		Print(report)
	}

	fmt.Println()
}

type Test struct {
	data  []mimo.DataRow
	index int
}

func (t *Test) ReadDataRow() (mimo.DataRow, error) {
	t.index++
	if t.index > len(t.data) {
		return nil, nil
	}

	return t.data[t.index-1], nil
}

func Print(r mimo.Report) {
	fmt.Println("Metrics")
	fmt.Println("=======")

	for key, metrics := range r {
		fmt.Println(key, metrics.MaskingRate()*100, "%") //nolint:gomnd
	}
}
