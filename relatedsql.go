package billtitles

import (
	"encoding/json"
	"io/fs"
	stdlog "log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const BILLSRELATED_DB = "billsrelated.db"

type BillToBill struct {
	gorm.Model
	Billnumber    string  `gorm:"UNIQUE_INDEX:compositeindex;index:,not null" json:"billnumber"`
	Billnumber_to string  `gorm:"UNIQUE_INDEX:compositeindex;index:,not null" json:"billnumber_to"`
	Reason        string  `gorm:"UNIQUE_INDEX:compositeindex;not null" json:"reason"`
	Score         float64 `json:"score"`
	ScoreOther    float64 `json:"score_other"` // score for other bill
	Identified_by string  `gorm:"index:,not null" json:"identified_by"`
}

type compareItem struct {
	Score        float64 `json:"score"`
	ScoreOther   float64 `json:"scoreother"` // score for other bill
	Explanation  string  `json:"explanation"`
	ComparedDocs string  `json:"compareddocs"`
}

type filterFunc func(string) bool

// Walk directory with a filter. Returns the filepaths that
// pass the 'testPath' function
// There is an exported function in the `bills` package that does this
func walkDirFilter(root string, testPath filterFunc) (filePaths []string, err error) {
	defer log.Info().Msg("Done collecting filepaths.")
	log.Info().Msgf("Getting all file paths in %s.  This may take a while.\n", root)
	filePaths = make([]string, 0)
	accumulate := func(fpath string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Err(err)
			return err
		}
		if testPath(fpath) {
			filePaths = append(filePaths, fpath)
		}
		return nil
	}
	err = filepath.WalkDir(root, accumulate)
	return
}

func GetRelatedDb(dbname string) *gorm.DB {
	if dbname == "" {
		dbname = BILLSRELATED_DB
	}
	var newLogger = logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	var db, _ = gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: newLogger,
	})
	// May not be necessary; applies all associated changes
	// when updating a title or bill.
	db.AutoMigrate(&BillToBill{})
	db.Session(&gorm.Session{FullSaveAssociations: true})
	return db
}

func similarCategoryJsonFilter(testPath string) bool {
	matched, err := regexp.MatchString(`esSimilarCategory`, testPath)
	matchedJson, err2 := regexp.MatchString(`\.json$`, testPath)
	if err != nil || err2 != nil {
		return false
	}
	return matched && matchedJson
}

func processRelatedJson(filePath string) (dat map[string]compareItem, err error) {
	defer func() {
		log.Info().Msgf("Finished processing: %s\n", filePath)
	}()
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Msgf("Error reading data.json: %s", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &dat)
	return dat, err
}

// jsonPath := RelatedJsonPath
// db := GetDb(BILLSRELATED_DB)
func LoadBillsRelatedToDBFromJson(db *gorm.DB, parentPath string) {
	log.Info().Msgf("Loading titles to database from json files in directory: %s", parentPath)
	defer log.Info().Msg("Done processing similar bill json")
	dataJsonFiles, error := walkDirFilter(parentPath, similarCategoryJsonFilter)
	if error != nil {
		log.Fatal().Msgf("Error getting files list: %s", error)
	}
	log.Info().Msgf("Found %d json files", len(dataJsonFiles))
	reportAt := 100
	count := 0
	for _, jpath := range dataJsonFiles {
		count++
		log.Debug().Msgf("Processing: %s", jpath)
		compareMap, err := processRelatedJson(jpath)
		if err != nil {
			log.Error().Msgf("Error processing json: %s", err)
			continue
		}

		log.Debug().Msgf("Got compare map: %+v\n", compareMap)

		if count%reportAt == 0 {
			log.Info().Msgf("Processed %d files", count)
		}

		for _, compare := range compareMap {
			bills := strings.Split(compare.ComparedDocs, "-")
			if len(bills) != 2 {
				log.Error().Msgf("Error parsing bills: %s", compare.ComparedDocs)
				continue
			}
			// TODO: can also test if the key of the map is equal to bills[1]

			// Create a billtobill item from compare and insert into db
			billtobill := BillToBill{
				Billnumber:    bills[0],
				Billnumber_to: bills[1],
				Reason:        compare.Explanation,
				Score:         compare.Score,
				ScoreOther:    compare.ScoreOther,
				Identified_by: "BillMap",
			}
			db.Create(&billtobill)
			log.Debug().Msgf("Saved compare item to db")
		}
	}
}
