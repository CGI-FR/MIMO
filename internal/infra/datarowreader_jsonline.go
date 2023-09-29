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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-json"

	"github.com/cgi-fr/mimo/pkg/mimo"
)

type DataRowReaderJSONLine struct {
	input  *bufio.Scanner
	output *bufio.Writer
}

func NewDataRowReaderJSONLineFromFile(filename string) (*DataRowReaderJSONLine, error) {
	source, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &DataRowReaderJSONLine{input: bufio.NewScanner(source), output: nil}, nil
}

func NewDataRowReaderJSONLine(input io.Reader, output io.Writer) *DataRowReaderJSONLine {
	return &DataRowReaderJSONLine{input: bufio.NewScanner(input), output: bufio.NewWriter(output)}
}

func (drr *DataRowReaderJSONLine) ReadDataRow() (mimo.DataRow, error) {
	var data mimo.DataRow

	if drr.input.Scan() {
		data = mimo.DataRow{}
		if err := json.UnmarshalNoEscape(drr.input.Bytes(), &data); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	if err := drr.input.Err(); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

func (drr *DataRowReaderJSONLine) ReadDataRowAndWrite() (mimo.DataRow, error) {
	var data mimo.DataRow

	if drr.input.Scan() {
		if _, err := drr.output.Write(append(drr.input.Bytes(), '\n')); err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		data = mimo.DataRow{}
		if err := json.UnmarshalNoEscape(drr.input.Bytes(), &data); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	if err := drr.input.Err(); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}

func (drr *DataRowReaderJSONLine) Flush() error {
	return fmt.Errorf("%w", drr.output.Flush())
}
