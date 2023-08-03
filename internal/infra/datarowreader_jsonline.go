package infra

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/adrienaury/mimo/pkg/mimo"
)

type DataRowReaderJSONLine struct {
	input  *bufio.Scanner
	output io.Writer
}

func NewDataRowReaderJSONLineFromFile(filename string) (*DataRowReaderJSONLine, error) {
	source, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &DataRowReaderJSONLine{input: bufio.NewScanner(source), output: io.Discard}, nil
}

func NewDataRowReaderJSONLine(input io.Reader, output io.Writer) *DataRowReaderJSONLine {
	return &DataRowReaderJSONLine{input: bufio.NewScanner(input), output: output}
}

func (drr *DataRowReaderJSONLine) ReadDataRow() (mimo.DataRow, error) {
	var data mimo.DataRow

	if drr.input.Scan() {
		if _, err := drr.output.Write(append(drr.input.Bytes(), '\n')); err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		data = mimo.DataRow{}
		if err := json.Unmarshal(drr.input.Bytes(), &data); err != nil {
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
