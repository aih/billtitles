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

func GetDb(dbname string) *gorm.DB {
	if dbname == "" {
		dbname = BILLTITLES_DB
	}
	var newLogger = logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	var db, _ = gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: newLogger,
	})
	// May not be necessary; applies all associated changes
	// when updating a title or bill.
	db.Session(&gorm.Session{FullSaveAssociations: true})
	return db
}

func AddTitleDb(db *gorm.DB, title string) {
	db.Create(&Title{Title: title})
}
func AddBillnumberversionsDb(db *gorm.DB, billnumberversions []string) {
	billnumberversions = RemoveDuplicates(billnumberversions)
	for _, billnumberversion := range billnumberversions {
		db.Create(&Bill{Billnumber: BillNumberVersionToBillNumber(billnumberversion), Billnumberversion: billnumberversion})
	}
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

func AddBillsToTitleDb(db *gorm.DB, title string, bills []*Bill) {
	db.Model(&Title{Title: title}).Association("Bills").Append(bills)
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

		db.Raw("SELECT bills.billnumber FROM bills, titles, bill_titles ON bill_titles.title_id = titles.id WHERE titles.Title IN ?", titleStrings).Scan(&bills)
		db.Raw("SELECT bills.* FROM bills, titles, bill_titleswhole ON bill_titleswhole.title_id = titles.id WHERE titles.Title IN ?", titleWholeStrings).Scan(&bills_whole)
	} else {
		return bills, bills_whole, fmt.Errorf("no bills found for billnumber %s", billnumber)
	}
	return bills, bills_whole, nil
}
