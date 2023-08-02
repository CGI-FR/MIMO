package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/adrienaury/mimo/internal/infra"
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

const maxCapacity int = 64 * 1024

func main() {
	//nolint: exhaustruct
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msgf("%v %v (commit=%v date=%v by=%v)", name, version, commit, buildDate, builtBy)

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	pipe := infra.CreateDataRowPipeWriter("/tmp/myFifo")
	defer pipe.Close()

	defer os.Remove("/tmp/myFifo")

	for scanner.Scan() {
		if err := pipe.Write(append(scanner.Bytes(), '\n')); err != nil {
			log.Error().Err(err).Msg("")
		}

		os.Stdout.Write(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("")
	}

	fmt.Println()
	time.Sleep(10 * time.Second)
}
