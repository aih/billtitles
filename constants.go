package billtitles

import (
	"path"
	"regexp"

	"github.com/aih/billtitles/internal/projectpath"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type LogLevel int8

type LogLevels map[string]zerolog.Level

type billVersions map[string]int

// Constants for this package
var (
	TitleIndex              = "titles.json"
	MainTitleIndex          = "maintitles.json"
	SampleTitleIndex        = "sampletitles.json"
	PathToDataDir           = "data"
	TitlesPath              = path.Join(PathToDataDir, TitleIndex)
	MainTitlesPath          = path.Join(PathToDataDir, MainTitleIndex)
	SampleTitlesPath        = path.Join(PathToDataDir, SampleTitleIndex)
	BillnumberRegexCompiled = regexp.MustCompile(`(?P<congress>[1-9][0-9]*)(?P<stage>[a-z]{1,8})(?P<billnumber>[1-9][0-9]*)(?P<version>[a-z]+)?`)
	BillVersionsOrdered     = billVersions{"ih": 0, "rh": 1, "rfs": 2, "eh": 3, "es": 4, "enr": 5}
	ZLogLevels              = LogLevels{"Debug": zerolog.DebugLevel, "Info": zerolog.InfoLevel, "Error": zerolog.ErrorLevel}
)

func LoadEnv() (err error) {
	return godotenv.Load(path.Join(projectpath.Root, ".env"))
}
