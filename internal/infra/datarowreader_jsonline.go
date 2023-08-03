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
	source io.Reader
}

func NewDataRowReaderJSONLineFromFile(filename string) (*DataRowReaderJSONLine, error) {
	source, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return &DataRowReaderJSONLine{source}, nil
}

func NewDataRowReaderJSONLine(source io.Reader) *DataRowReaderJSONLine {
	return &DataRowReaderJSONLine{source}
}

func (drr *DataRowReaderJSONLine) ReadDataRow() (mimo.DataRow, error) {
	reader := bufio.NewReader(drr.source)

	jsondata, err := reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}

		return nil, fmt.Errorf("%w", err)
	}

	data := mimo.DataRow{}
	if err := json.Unmarshal(jsondata, &data); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return data, nil
}
