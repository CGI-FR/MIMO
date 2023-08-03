package main

import (
	"fmt"
	"os"

	"github.com/adrienaury/mimo/internal/infra"
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

func main() {
	//nolint: exhaustruct
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msgf("%v %v (commit=%v date=%v by=%v)", name, version, commit, buildDate, builtBy)

	realReader, err := infra.NewDataRowReaderJSONLineFromFile("input.jsonl")
	maskedReader := infra.NewDataRowReaderJSONLine(os.Stdin)

	if err != nil {
		log.Fatal().Err(err).Msg("end of program")
	}

	driver := mimo.NewDriver()
	if report, err := driver.Analyze(realReader, maskedReader, infra.SubscriberLogger{}); err != nil {
		log.Error().Err(err).Msg("end of program")
	} else {
		report.Print()
	}

	fmt.Println()
}
