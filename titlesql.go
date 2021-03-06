package billtitles

import (
	"fmt"
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Title struct {
	gorm.Model
	Title string  `gorm:"index:,unique"`
	Bills []*Bill `gorm:"many2many:bill_titles;"`
}
type Bill struct {
	gorm.Model
	Billnumber        string   `gorm:"not null" json:"billnumber"`
	Billnumberversion string   `gorm:"index:,unique" json:"billnumberversion"`
	Titles            []*Title `gorm:"many2many:bill_titles;" json:"titles"`
	TitlesWhole       []*Title `gorm:"many2many:bill_titleswhole" json:"titles_whole"`
}

type BillTitle struct {
	gorm.Model
	Billnumber     string `gorm:"index:,not null" json:"billnumber"`
	Title          string `gorm:"index:,not null" json:"title"`
	IsForWholeBill bool   `gorm:"not null" json:"is_for_whole_bill"`
}

const BILLTITLES_DB = "billtitles.db"

func BillNumberVersionToBillNumber(billNumberVersion string) string {
	return BillnumberRegexCompiled.ReplaceAllString(billNumberVersion, "$1$2$3")
}

func BillNumberVersionsToBillNumbers(billNumberVersions []string) (billNumbers []string) {
	for _, billNumberVersion := range billNumberVersions {
		billNumber := BillNumberVersionToBillNumber(billNumberVersion)
		billNumbers = append(billNumbers, billNumber)
	}
	billNumbers = RemoveDuplicates(billNumbers)
	return billNumbers
}

func GetDb_old(dbname string) *gorm.DB {
	if dbname == "" {
		dbname = BILLTITLES_DB
	}
	var newLogger = logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	var db, _ = gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: newLogger,
	})
	// May not be necessary; applies all associated changes
	// when updating a title or bill.
	db.AutoMigrate(&Bill{}, &Title{})
	db.Session(&gorm.Session{FullSaveAssociations: true})
	return db
}

func GetDb(dbname string) *gorm.DB {
	if dbname == "" {
		dbname = BILLTITLES_DB
	}
	var newLogger = logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	var db, _ = gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: newLogger,
	})

	db.AutoMigrate(&BillTitle{})
	//db.Session(&gorm.Session{FullSaveAssociations: true})
	return db
}

func AddTitleDb(db *gorm.DB, title string) *gorm.DB {
	tx := db.Create(&Title{Title: title})
	return tx
}

func AddTitleStructDb(db *gorm.DB, title *Title) {
	// For example
	// bills = [&Bill{Billnumberversion: "117hr100ih", Billnumber: "117hr100"},
	// &Bill{Billnumberversion: "117hr222ih", Billnumber: "117hr222"}]
	//newTitle := &Title{Title: titleString, Bills: bills}

	db.Create(title)
}

func AddBillnumberversionsDb(db *gorm.DB, billnumberversions []string) {
	billnumberversions = RemoveDuplicates(billnumberversions)
	for _, billnumberversion := range billnumberversions {
		db.Create(&Bill{Billnumber: BillNumberVersionToBillNumber(billnumberversion), Billnumberversion: billnumberversion})
	}
}

func AddBillStructDb(db *gorm.DB, billStruct *Bill) {
	db.Create(&billStruct)
}

func RemoveTitleDb(db *gorm.DB, title string) {
	db.Model(&Bill{}).Association("Titles").Delete(title)
	db.Model(&Bill{}).Association("TitlesWhole").Delete(title)
	db.Where("Title = ?", title).Delete(&Title{})
}

func GetBillsByTitleDb(db *gorm.DB, title string) (bills []*Bill) {
	var titleStruct Title
	db.First(&titleStruct, "Title = ?", title)
	db.Model(titleStruct).Association("Bills").Find(&bills)
	return bills
}

func GetTitlesByBillnumberDb(db *gorm.DB, billnumber string) (titles []*Title) {
	var bills []*Bill
	db.Where("Billnumber = ?", billnumber).Find(&bills)
	db.Model(bills).Association("Titles").Find(&titles)
	log.Debug().Msgf("titles: %+v", titles)
	return titles
}

func GetTitlesByBillnumberVersionDb(db *gorm.DB, billnumberversion string) (titles []*Title) {
	var bills []*Bill
	db.Where("Billnumberversion = ?", billnumberversion).Find(&bills)
	db.Model(bills).Association("Titles").Find(&titles)
	log.Debug().Msgf("titles: %+v", titles)
	return titles
}

func GetTitlesWholeByBillnumberDb(db *gorm.DB, billnumber string) (titles []*Title) {
	var bills []*Bill
	db.Where("Billnumber = ?", billnumber).Find(&bills)
	db.Model(bills).Association("TitlesWhole").Find(&titles)
	log.Debug().Msgf("titles: %+v", titles)
	return titles
}

func GetTitlesWholeByBillnumberVersionDb(db *gorm.DB, billnumberversion string) (titles []*Title) {
	var bills []*Bill
	db.Where("Billnumberversion = ?", billnumberversion).Find(&bills)
	db.Model(bills).Association("TitlesWhole").Find(&titles)
	log.Debug().Msgf("titles: %+v", titles)
	return titles
}

func GetBillsWithSameTitleDb(db *gorm.DB, billnumber string) (bills, bills_whole []*Bill, err error) {
	var titles []*Title
	var titleswhole []*Title
	db.Where("Billnumber = ?", billnumber).Find(&bills)
	if len(bills) > 0 {
		//log.Info().Msgf("Found bills: %+v", bills[0].Billnumberversion)
		db.Model(&bills).Association("Titles").Find(&titles)
		if len(titles) > 0 {
			for _, title := range titles {
				log.Debug().Msgf("Found title: %+v", title.Title)
			}
		}
		db.Model(&bills).Association("TitlesWhole").Find(&titleswhole)
		if len(titleswhole) > 0 {
			for _, title := range titles {
				log.Debug().Msgf("Found title (whole): %+v", title.Title)
			}
		}
		var titleStrings []string
		var titleWholeStrings []string
		for _, title := range titles {
			titleStrings = append(titleStrings, title.Title)
		}
		for _, title := range titleswhole {
			titleWholeStrings = append(titleWholeStrings, title.Title)
		}

		log.Info().Msgf("Found titles related to this bill: %+v", titleStrings)

		// Working!
		db.Raw("SELECT bills.billnumber, titles.Title FROM bills INNER JOIN bill_titles ON bills.id = bill_titles.bill_id INNER JOIN titles ON titles.id = bill_titles.title_id WHERE titles.Title IN ?", titleStrings).Scan(&bills)
		log.Debug().Msgf("Found %d bills related to this bill's titles", len(bills))

		db.Raw("SELECT bills.billnumber, titles.Title FROM bills INNER JOIN bill_titleswhole ON bills.id = bill_titleswhole.bill_id INNER JOIN titles ON titles.id = bill_titleswhole.title_id WHERE titles.Title IN ?", titleWholeStrings).Scan(&bills_whole)
		log.Debug().Msgf("Found %d bills related to this bill's whole titles", len(bills_whole))
	} else {
		return bills, bills_whole, fmt.Errorf("no bills found for billnumber %s", billnumber)
	}
	return bills, bills_whole, nil
}

/*
   TODO: Create a service to:
    	- Add a bill+title and bill+titlewhole to the database.
    	- Remove a title or titlewhole from the database.
     	- Query the database by billnumber for related titles and titleswhole.
     	- Query the database by title (string) for bills that have that title or whole title
*/

// jsonPath := TitlesPath
// db := GetDb(BILLTITLES_DB)
// NOTE: This adds all bills at the 'ih' version,
// since we currently only have the json map with billnumbers, not billnumberversions
func LoadTitlesToDBFromJson(db *gorm.DB, jsonPath string) {
	log.Info().Msgf("Loading titles to database from json file: %s", jsonPath)
	titleMap, error := LoadTitlesMap(jsonPath)
	if error != nil {
		log.Fatal().Msgf("Error loading titles: %s", error)
	}
	titleMap.Range(func(key, value interface{}) bool {
		billnumbers := value.([]string)
		if len(billnumbers) > 0 {
			log.Debug().Msgf("Adding billnumbers %+v to title `%s`", billnumbers, key)
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

func LoadBillTitleToDBFromJson(db *gorm.DB, jsonPath string, isForWhole bool) {
	log.Info().Msgf("Loading titles to database from json file: %s", jsonPath)
	titleMap, error := LoadTitlesMap(jsonPath)
	if error != nil {
		log.Fatal().Msgf("Error loading titles: %s", error)
	}
	titleMap.Range(func(key, value interface{}) bool {
		title := key.(string)
		billnumbers := value.([]string)
		if len(billnumbers) > 0 {
			log.Debug().Msgf("Adding billnumbers %+v to title `%s`", billnumbers, key)
			for _, billnumber := range billnumbers {
				billTitle := BillTitle{Billnumber: billnumber, Title: title, IsForWholeBill: isForWhole}
				db.Create(&billTitle)
			}

		}
		return true
	})
}
