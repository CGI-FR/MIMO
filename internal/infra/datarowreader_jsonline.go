package infra

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/adrienaury/mimo/pkg/mimo"
)

type DataRowReaderJSONLine struct {
	input  io.Reader
	output io.Writer
}

func NewDataRowReaderJSONLineFromFile(filename string) (*DataRowReaderJSONLine, error) {
	source, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &DataRowReaderJSONLine{input: source, output: os.NewFile(0, os.DevNull)}, nil
}

func NewDataRowReaderJSONLine(input io.Reader, output io.Writer) *DataRowReaderJSONLine {
	return &DataRowReaderJSONLine{input: input, output: output}
}

func (drr *DataRowReaderJSONLine) ReadDataRow() (mimo.DataRow, error) {
	reader := bufio.NewReader(drr.input)

	jsondata, err := reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	drr.output.Write(append(jsondata, '\n'))

	data := mimo.DataRow{}
	if err := json.Unmarshal(jsondata, &data); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}
