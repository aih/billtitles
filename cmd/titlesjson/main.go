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
	flag.Parse()
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msg("Log level set to Debug")

	titleMap, error := billtitles.LoadTitlesMap(billtitles.TitlesPath)
	if error != nil {
		log.Info().Msgf("%v", error)
		panic(error)
	}
	billnumbers, err := billtitles.GetBillnumbersByTitle(titleMap, sampleTitle)
	if err != nil {
		log.Info().Msgf("Error getting related bills: %v", err)
	}
	log.Debug().Msgf("%v", billnumbers)

	// TODO: create services to expose the map and add/remove functions

}
