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

	pipe := infra.CreateDataRowPipeReader("/tmp/myFifo")
	input := infra.NewDataRowScanner()

	driver := mimo.NewDriver()
	if report, err := driver.Analyze(pipe, input); err != nil {
		log.Fatal().AnErr("error", err).Msg("end of program")
	} else {
		Print(report)
	}

	fmt.Println()
}

func Print(r mimo.Report) {
	fmt.Println("Metrics")
	fmt.Println("=======")

	for key, metrics := range r {
		fmt.Println(key, metrics.MaskingRate()*100, "%") //nolint:gomnd
	}
}
