package billtitles

import (
	"errors"
	stdlog "log"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const BILLSRELATED_DB = "billsrelated.db"

type BillToBill struct {
	gorm.Model
	Billnumber    string `gorm:"index:,not null" json:"billnumber"`
	Billnumber_to string `gorm:"index:,not null" json:"billnumber_to"`
	Reason        string `gorm:"not null" json:"reason"`
	Identified_by string `gorm:"index:,not null" json:"identified_by"`
}

func GetRelatedDb(dbname string) *gorm.DB {
	if dbname == "" {
		dbname = BILLTITLES_DB
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

func LoadRelatedMap(jsonPath string) (*sync.Map, error) {
	relatedMap := new(sync.Map)
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		jsonPath = TitlesPath
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			return relatedMap, errors.New("related bills file file not found")
		}
	}
	log.Debug().Msgf("Path to JSON file: %s", jsonPath)
	var err error
	relatedMap, err = UnmarshalJsonFile(jsonPath)
	if err != nil {
		return nil, err
	} else {
		return relatedMap, nil
	}
}

//TODO create related sync.Map from individual json files

// jsonPath := RelatedJsonPath
// db := GetDb(BILLSRELATED_DB)
func LoadBillsRelatedToDBFromJson(db *gorm.DB, jsonPath string) {
	log.Info().Msgf("Loading titles to database from json file: %s", jsonPath)
	relatedMap, error := LoadRelatedMap(jsonPath)
	if error != nil {
		log.Fatal().Msgf("Error loading titles: %s", error)
	}
	relatedMap.Range(func(key, value interface{}) bool {
		billnumbers := value.([]string)
		if len(billnumbers) > 0 {
			log.Info().Msgf("Adding billnumbers %+v to title `%s`", billnumbers, key)
			bills := []*Bill{}
			for _, billnumber := range billnumbers {
				billnumberversion := billnumber + "ih"
				bills = append(bills, &Bill{Billnumber: billnumber, Billnumberversion: billnumberversion})
			}
			title := Title{Title: key.(string), Bills: bills}
			db.Create(&title)

		}
		return true
	})
}
