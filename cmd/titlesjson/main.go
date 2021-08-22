package main

import (
	"flag"
	"os"
	"sync"

	"github.com/aih/billtitles"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const sampleTitle = "21st Century Energy Workforce Act"

func makeSampleTitlesFile(titleMap *sync.Map) {
	defer func() { log.Info().Msg("done making samples file") }()
	sampleTitles := new(sync.Map)
	count := 0
	titleMap.Range(func(key, value interface{}) bool {
		count += 1
		if count > 4 {
			log.Debug().Msgf("Returning 'false' from range loop")
			return false
		}
		log.Debug().Msgf("Adding a title")
		sampleTitles.Store(key.(string), value.([]string))
		return true
	})
	billtitles.SaveTitlesMap(sampleTitles, "data/sampletitles.json")
}

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

	titleMap, error := billtitles.LoadTitlesMap(billtitles.MainTitlePath)
	if error != nil {
		log.Info().Msgf("%v", error)
		panic(error)
	} else {
		makeSampleTitlesFile(titleMap)
		billnumbers, ok := titleMap.Load(sampleTitle)
		if !ok {
			log.Error().Msgf("No bills found for: %s", sampleTitle)
		} else {
			log.Info().Msgf("Bill numbers for sample bill (%s): %v", sampleTitle, billnumbers)
		}
		billtitles.AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr222"})
		billtitles.AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr999"})
		//billtitles.SaveTitlesMap(titleMap, billtitles.MainTitlePath)
		newbillnumbers, err := billtitles.GetBillnumbersByTitle(titleMap, "A test title")
		if err != nil {
			log.Error().Msgf("Error getting bill numbers for title: %s", err)
		} else {
			log.Info().Msgf("Bill numbers for title: %s", newbillnumbers)
		}
		billtitles.RemoveTitle(titleMap, "A test title")
		newbillnumbers, err = billtitles.GetBillnumbersByTitle(titleMap, "A test title")
		if err != nil {
			log.Error().Msgf("Error getting bill numbers for title: %s", err)
		} else {
			log.Info().Msgf("Bill numbers for title: %s", newbillnumbers)
		}
	}

}
