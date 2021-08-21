package main

import (
	"flag"
	"os"

	"github.com/aih/billtitles"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const sampleTitle = "21st Century Energy Workforce Act"

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msg("Log level set to Debug")
	flag.Parse()

	titleMap, error := billtitles.LoadTitlesMap(billtitles.MainTitlePath)
	if error != nil {
		log.Info().Msgf("%v", error)
		panic(error)
	} else {
		billnumbers, ok := titleMap.Load(sampleTitle)
		if !ok {
			log.Error().Msgf("No bills found for: %s", sampleTitle)
		} else {
			log.Info().Msgf("Bill numbers for sample bill (%s): %v", sampleTitle, billnumbers)
		}
		billtitles.RemoveTitle(titleMap, "")
		billtitles.SaveTitlesMap(titleMap, billtitles.MainTitlePath)
	}

}
