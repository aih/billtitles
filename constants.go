package billtitles

import (
	"path"

	"github.com/aih/billtitles/internal/projectpath"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type LogLevel int8

type LogLevels map[string]zerolog.Level

type billVersions map[string]int

// Constants for this package
var (
	TitleIndex          = "titles.json"
	MainTitleIndex      = "maintitles.json"
	SampleTitleIndex    = "sampletitles.json"
	PathToDataDir       = "data"
	TitlesPath          = path.Join(PathToDataDir, TitleIndex)
	MainTitlePath       = path.Join(PathToDataDir, MainTitleIndex)
	SampleTitlePath     = path.Join(PathToDataDir, SampleTitleIndex)
	BillVersionsOrdered = billVersions{"ih": 0, "rh": 1, "rfs": 2, "eh": 3, "es": 4, "enr": 5}
	ZLogLevels          = LogLevels{"Debug": zerolog.DebugLevel, "Info": zerolog.InfoLevel, "Error": zerolog.ErrorLevel}
)

func LoadEnv() (err error) {
	return godotenv.Load(path.Join(projectpath.Root, ".env"))
}
