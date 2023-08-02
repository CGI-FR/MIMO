package infra

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/adrienaury/mimo/pkg/mimo"
)

type DataRowScanner struct {
	*bufio.Scanner
}

func NewDataRowScanner() *DataRowScanner {
	return &DataRowScanner{bufio.NewScanner(os.Stdin)}
}

func (s *DataRowScanner) ReadDataRow() (mimo.DataRow, error) {
	if s.Scan() {
		data := mimo.DataRow{}
		if err := json.Unmarshal(s.Bytes(), &data); err != nil {
			return nil, err
		}

		return data, nil
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}
