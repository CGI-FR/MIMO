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

package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/cgi-fr/mimo/internal/infra"
	"github.com/cgi-fr/mimo/pkg/mimo"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	name      string // provisioned by ldflags
	version   string // provisioned by ldflags
	commit    string // provisioned by ldflags
	buildDate string // provisioned by ldflags
	builtBy   string // provisioned by ldflags

	verbosity string
	jsonlog   bool
	debug     bool
	colormode string

	configfile  string
	diskStorage bool
	persist     string
)

func main() {
	cobra.OnInitialize(initLog)

	rootCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   fmt.Sprintf("%v real-data-file.jsonl", name),
		Short: "Masked Input Metrics Output",
		Long:  `MIMO is a purpose-built tool designed for assessing the quality of a pseudonymization process.`,
		Example: `  # create a pipe file to store the real json stream before pseudonymization
  > mkfifo real.jsonl
  # pseudonymize data with LINO and PIMO and verify the result with MIMO
  > lino pull prod | tee real.jsonl | pimo | mimo real.jsonl | lino push dev`,
		Version: fmt.Sprintf(`%v (commit=%v date=%v by=%v)
Copyright (C) 2021 CGI France
License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.`, version, commit, buildDate, builtBy),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Info().
				Str("verbosity", verbosity).
				Bool("log-json", jsonlog).
				Bool("debug", debug).
				Str("color", colormode).
				Msg("start MIMO")
		},
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := run(cmd, args[0]); err != nil {
				log.Fatal().Err(err).Msg("end MIMO")
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			log.Info().Int("return", 0).Msg("end MIMO")
		},
	}

	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", "warn",
		"set level of log verbosity : none (0), error (1), warn (2), info (3), debug (4), trace (5)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "add debug information to logs (very slow)")
	rootCmd.PersistentFlags().BoolVar(&jsonlog, "log-json", false, "output logs in JSON format")
	rootCmd.PersistentFlags().StringVar(&colormode, "color", "auto", "use colors in log outputs : yes, no or auto")
	rootCmd.PersistentFlags().StringVar(&configfile, "config", "", "name of the YAML configuration file to use")

	rootCmd.PersistentFlags().BoolVar(&diskStorage, "disk-storage", false, "enable data storage on disk")
	rootCmd.PersistentFlags().StringVar(&persist, "persist", "",
		"persist data in the specified directory (implies --disk-storage)")

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).Msg("error when executing command")
		os.Exit(1)
	}
}

func run(_ *cobra.Command, realJSONLineFileName string) error {
	realReader, err := infra.NewDataRowReaderJSONLineFromFile(realJSONLineFileName)
	maskedReader := infra.NewDataRowReaderJSONLine(os.Stdin, os.Stdout)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var config mimo.Config

	if configfile != "" {
		if config, err = infra.LoadConfig(configfile); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	driver := mimo.NewDriver(
		realReader,
		maskedReader,
		selectMultimapFactory(),
		selectCounterFactory(),
		infra.SubscriberLogger{},
	)

	defer driver.Close()

	driver.Configure(config)

	haserror := false

	var report mimo.Report

	if report, err = driver.Analyze(); err != nil {
		return fmt.Errorf("%w", err)
	}

	columns := report.Columns()
	sort.Strings(columns)

	for _, colname := range columns {
		haserror = appendColumnMetric(report, colname, haserror)
	}

	if err = infra.NewReportExporter().Export(report, "report.html"); err != nil {
		return fmt.Errorf("%w", err)
	}

	if haserror {
		return errUnsatisfiedConstraint
	}

	return nil
}

func selectMultimapFactory() mimo.MultimapFactory {
	var multimapFactory mimo.MultimapFactory

	if !diskStorage && persist == "" {
		multimapFactory = func(fieldname string) mimo.Multimap {
			return mimo.Multimap{Backend: mimo.InMemoryMultimapBackend{}}
		}
	} else {
		multimapFactory = func(fieldname string) mimo.Multimap {
			path := ""
			if persist != "" {
				path = persist + "/" + fieldname
			}
			factory, err := infra.PebbleMultimapFactory(path)
			if err != nil {
				log.Fatal().Err(err).Msg("End MIMO")
			}

			return factory
		}
	}

	return multimapFactory
}

func selectCounterFactory() mimo.CounterFactory {
	var counterFactory mimo.CounterFactory

	if persist == "" {
		counterFactory = func(fieldname string) mimo.CounterBackend {
			return &mimo.InMemoryCounterBackend{
				TotalCount:   0,
				NilCount:     0,
				IgnoredCount: 0,
				MaskedCount:  0,
			}
		}
	} else {
		counterFactory = func(fieldname string) mimo.CounterBackend {
			path := ""
			if persist != "" {
				path = persist + "/" + fieldname
			}
			factory, err := infra.PebbleCounterBackendFactory(path)
			if err != nil {
				log.Fatal().Err(err).Msg("End MIMO")
			}

			return factory
		}
	}

	return counterFactory
}

func appendColumnMetric(report mimo.Report, colname string, haserror bool) bool {
	metrics := report.ColumnMetric(colname)
	if metrics.Validate() >= 0 {
		log.Info().
			Str("field", colname).
			Int64("count-nil", metrics.NilCount()).
			Int64("count-ignored", metrics.IgnoredCount()).
			Int64("count-masked", metrics.MaskedCount()).
			Int64("count-missed", metrics.NonMaskedCount()).
			Float64("rate-masking", metrics.MaskedRate()).
			Float64("rate-coherence", metrics.Coherence.Rate()).
			Float64("rate-identifiable", metrics.Identifiant.Rate()).
			Msg("summmary for column " + colname)
	} else {
		log.Error().
			Str("field", colname).
			Int64("count-nil", metrics.NilCount()).
			Int64("count-ignored", metrics.IgnoredCount()).
			Int64("count-masked", metrics.MaskedCount()).
			Int64("count-missed", metrics.NonMaskedCount()).
			Float64("rate-masking", metrics.MaskedRate()).
			Float64("rate-coherence", metrics.Coherence.Rate()).
			Float64("rate-identifiable", metrics.Identifiant.Rate()).
			Msg("summmary for column " + colname)
		haserror = true

		logSamples("coherence", "real-value", "pseudonyms", metrics.GetInvalidSamplesForCoherentRate(10))      //nolint:gomnd
		logSamples("identifiant", "pseudonym", "real-values", metrics.GetInvalidSamplesForIdentifiantRate(10)) //nolint:gomnd
	}

	return haserror
}

func logSamples(target, labelForValue, labelForAssigned string, samples []mimo.Sample) {
	for _, sample := range samples {
		lenMax := fmt.Sprintf("%d", len(sample.AssignedValues))

		if len(sample.AssignedValues) > 10 { //nolint:gomnd
			sample.AssignedValues = sample.AssignedValues[:10]
		}

		log.Error().
			Str(labelForValue, sample.OriginalValue).
			Strs(labelForAssigned, sample.AssignedValues).
			Msg("sample value that failed " + target + " because it was attributed " + lenMax + " " + labelForAssigned)
	}
}

func initLog() {
	color := false

	switch strings.ToLower(colormode) {
	case "auto":
		if isatty.IsTerminal(os.Stdout.Fd()) && runtime.GOOS != "windows" {
			color = true
		}
	case "yes", "true", "1", "on", "enable":
		color = true
	}

	if jsonlog {
		log.Logger = zerolog.New(os.Stderr)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !color}) //nolint:exhaustruct
	}

	if debug {
		log.Logger = log.Logger.With().Caller().Logger()
	}

	setVerbosity()
}

func setVerbosity() {
	switch verbosity {
	case "trace", "5":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug", "4":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info", "3":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn", "2":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error", "1":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
