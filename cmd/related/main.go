package main

import (
	"flag"
	"os"

	"github.com/aih/billtitles"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type flagDef struct {
	value string
	usage string
}

func main() {
	flagDefs := map[string]flagDef{
		"parentPath": {"../congress", "Absolute path to the parent directory for json metadata files"},
		"log":        {"Info", "Sets Log level. Options: Error, Info, Debug"},
	}
	debug := flag.Bool("debug", false, "sets log level to debug")

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	var parentPath string
	flag.StringVar(&parentPath, "parentPath", flagDefs["parentPath"].value, flagDefs["parentPath"].usage)
	flag.StringVar(&parentPath, "p", flagDefs["parentPath"].value, flagDefs["parentPath"].usage+" (shorthand)")

	var logLevel string
	flag.StringVar(&logLevel, "logLevel", flagDefs["log"].value, flagDefs["log"].usage)
	flag.StringVar(&logLevel, "l", flagDefs["log"].value, flagDefs["log"].usage+" (shorthand)")

	flag.Parse()
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msg("Log level set to Debug")
	db := GetRelatedDb(BILLTITLES_DB)
	billtitles.LoadBillsRelatedToDBFromJson(parentPath)

}
