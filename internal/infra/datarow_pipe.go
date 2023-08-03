package infra

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/adrienaury/mimo/pkg/mimo"
)

type DataRowPipe struct {
	pipe     *os.File
	filePath string
}

func CreateDataRowPipeWriter(filePath string) *DataRowPipe {
	err := syscall.Mkfifo(filePath, 0o640) //nolint:gomnd
	if err != nil {
		fmt.Println("Failed to create pipe")
		panic(err)
	}

	pipe, err := os.OpenFile(filePath, os.O_WRONLY, 0o600) //nolint:gomnd
	if err != nil {
		panic(err)
	}

	return &DataRowPipe{pipe, filePath}
}

func CreateDataRowPipeReader(filePath string) *DataRowPipe {
	pipe, err := os.OpenFile(filePath, os.O_RDONLY, 0o600) //nolint:gomnd

	if errors.Is(err, os.ErrNotExist) {
		time.Sleep(time.Second)

		pipe, err = os.OpenFile(filePath, os.O_RDONLY, 0o600) //nolint:gomnd
	}

	if err != nil {
		fmt.Println("Couldn't open pipe with error: ", err)
	}

	return &DataRowPipe{pipe, filePath}
}

func (p *DataRowPipe) ReadDataRow() (mimo.DataRow, error) {
	reader := bufio.NewReader(p.pipe)

	jsondata, err := reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}

		return nil, err
	}

	data := mimo.DataRow{}
	if err := json.Unmarshal(jsondata, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (p *DataRowPipe) Write(bytes []byte) error {
	if _, err := p.pipe.Write(bytes); err != nil {
		return err
	}

	return nil
}

func (p *DataRowPipe) Close() {
	_ = p.pipe.Close()
}

func (p *DataRowPipe) Remove() {
	os.Remove(p.filePath)
}
