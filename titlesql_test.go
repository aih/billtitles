package billtitles

import (
	stdlog "log"
	"os"
	"testing"
	"time"

	"github.com/aih/billtitles/internal/testutils"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var newLogger = logger.New(
	stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Warn, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		Colorful:                  false,       // Disable color
	},
)

func TestBillNumberVersionToBillNumber(t *testing.T) {
	assert.Equal(t, "116hr200", BillNumberVersionToBillNumber("116hr200ih"))
}

func TestBillNumberVersionsToBillNumbers(t *testing.T) {
	assert.Equal(t, []string{"116hr200", "117hjres20"}, BillNumberVersionsToBillNumbers([]string{"116hr200ih", "117hjres20ih"}))
}

func TestCreateAndGetTitle(t *testing.T) {
	// Setup
	var db, _ = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})

	// Migrate the schema
	db.AutoMigrate(&Bill{}, &Title{})
	db.Session(&gorm.Session{FullSaveAssociations: true})
	testutils.SetLogLevel()
	log.Debug().Msg("Test setting and getting titles and bills from sql db")

	const (
		titleString  = "This is a test title"
		titleString2 = "This is a new test title"
		titleString3 = "This Title Added with AddTitleDB function"
		titleString4 = "This is a title for a whole bill"
		titleString5 = "Title to append bills to"
	)
	newBill1 := &Bill{Billnumberversion: "116hr1500ih", Billnumber: "116hr1500"}
	newBill2 := &Bill{Billnumberversion: "117hr200ih", Billnumber: "117hr200"}
	newBill3 := &Bill{Billnumberversion: "117hr100ih", Billnumber: "117hr100"}
	newBill4 := &Bill{Billnumberversion: "117hr222ih", Billnumber: "117hr222"}

	// Associates titleString with newBill1 and newBill3
	newTitle := &Title{Title: titleString, Bills: []*Bill{newBill1, newBill3}}
	db.Session(&gorm.Session{FullSaveAssociations: true})

	// Create
	db.Create(newTitle)

	// Associates titleString with newBill4, so now: newBill1, newBill3, newBill4
	db.Model(newTitle).Association("Bills").Append([]Bill{*newBill4})

	// newBill2 is not associated with any bill, so it should not be in the list of associated bills
	db.Create(newBill2)

	// Read
	var title Title
	var associatedBills []*Bill

	t.Run("Get Title", func(t *testing.T) {
		db.First(&title, "title = ?", titleString) // find title
		log.Debug().Msgf("Got title item: %v", title)
		log.Debug().Msgf("Title item has the title: %v", title.Title)
		assert.Equal(t, titleString, title.Title)
	})
	t.Run("Get Associated Bills", func(t *testing.T) {
		db.First(&title, "title = ?", titleString) // find title with Title = "This is a test title"
		// Should be associated with newBill1, newBill3, newBill4
		db.Model(title).Association("Bills").Find(&associatedBills)
		assert.NotEqual(t, 0, len(associatedBills))
		log.Debug().Msgf("Associated bills: %+v", associatedBills)
		assert.Equal(t, 3, len(associatedBills))
	})

	t.Run("Get Associated bills using GetBillsByTitleDb", func(t *testing.T) {
		// Should be associated with newBill1, newBill3, newBill4
		associatedBills := GetBillsByTitleDb(db, titleString)
		assert.NotEqual(t, 0, len(associatedBills))
		log.Debug().Msgf("Associated bills: %+v", associatedBills)
		assert.Equal(t, 3, len(associatedBills))
	})

	t.Run("Test GetTitlesByBillnumberDb", func(t *testing.T) {
		// Should be associated with titleString only
		queryString := "116hr1500"
		titles := GetTitlesByBillnumberDb(db, queryString)
		assert.NotEqual(t, 0, len(titles))
		log.Debug().Msgf("Associated titles: %+v", titles)
		assert.Equal(t, 1, len(titles))
		assert.Equal(t, titleString, titles[0].Title)
	})
	t.Run("Test GetTitlesByBillnumberVersionDb", func(t *testing.T) {
		// Should be associated with titleString only
		queryString := "116hr1500ih"
		titles := GetTitlesByBillnumberVersionDb(db, queryString)
		if len(titles) > 0 {
			log.Info().Msgf("Associated titles length: %+v", len(titles))
			log.Info().Msgf("Associated titles: %+v", titles[0].Title)
		}
		assert.NotEqual(t, 0, len(titles))
		assert.Equal(t, 1, len(titles))
		assert.Equal(t, titleString, titles[0].Title)
	})

	t.Run("Test GetTitlesByBillnumberVersionDb where no associated titles", func(t *testing.T) {
		// Should not have associated titles
		queryString := "117hr200ih"
		titles := GetTitlesByBillnumberVersionDb(db, queryString)
		assert.Equal(t, 0, len(titles))
	})

	t.Run("Test GetTitlesWholeByBillnumberVersionDb", func(t *testing.T) {
		queryString := "117hr100ih"
		titles := GetTitlesByBillnumberVersionDb(db, queryString)
		assert.NotEqual(t, 0, len(titles))
		log.Debug().Msgf("Associated titles: %+v", titles)
		assert.Equal(t, 1, len(titles))
		assert.Equal(t, titleString, titles[0].Title)
	})

	t.Run("Add a title entry with db.Model", func(t *testing.T) {
		// Update - update title from "This is a test title" to "This is a new test title"
		db.Model(&title).Update("Title", titleString2)
		log.Debug().Msgf("Title '%v' updated", title.Title)
		assert.Equal(t, titleString2, title.Title)
	})

	t.Run("Add a title (whole) entry with db.Model", func(t *testing.T) {
		newTitle4 := &Title{Title: titleString4, Bills: []*Bill{newBill1, newBill3}}
		db.Create(&newTitle4)
		log.Debug().Msgf("Title '%v' added, with bills %s and %s", newTitle4.Title, newBill1.Billnumber, newBill3.Billnumber)
		db.Model(&newBill1).Association("TitlesWhole").Append(newTitle4)
		var associatedTitles []*Title
		db.Model(&newBill1).Association("TitlesWhole").Find(&associatedTitles)
		assert.NotNil(t, associatedTitles)
		assert.NotEqual(t, 0, len(associatedTitles))
	})

	var bill Bill
	var titles []Title

	t.Run("Query by ID", func(t *testing.T) {
		db.First(&bill, "ID = ?", "1") // find item with any of the billnumbers in bctnvs
		log.Debug().Msgf("Got bill item for ID=1: %+v", bill)
		assert.Equal(t, "116hr1500ih", bill.Billnumberversion)
	})

	t.Run("Query by billnumber (no titles)", func(t *testing.T) {
		var bill2 Bill
		querystring := "117hr200"
		db.First(&bill2, "Billnumber = ?", querystring) // find item with billnumber = 117hr200
		log.Debug().Msgf("Got bill item for billnumber=%s: %+v", querystring, bill2)
		assert.Equal(t, "117hr200ih", bill2.Billnumberversion)
		db.Model(&bill2).Association("Titles").Find(&titles)
		assert.Equal(t, 0, len(titles))
	})

	t.Run("Query by billnumber (has titles)", func(t *testing.T) {
		var bill *Bill
		querystring := "116hr1500"
		db.First(&bill, "Billnumber = ?", querystring) // find item with billnumber = 116hr1500
		log.Debug().Msgf("Got bill item for billnumber=%s: %+v", querystring, bill)
		assert.Equal(t, "116hr1500ih", bill.Billnumberversion)
		db.Model(&bill).Association("Titles").Find(&titles)
		log.Debug().Msgf("Associated titles: %+v", titles)
		assert.Equal(t, 2, len(titles))
		var titleswhole []*Title
		db.Model(&bill).Association("TitlesWhole").Find(&titleswhole)
		//log.Debug().Msgf("Associated titles (whole): %+v", titleswhole[0].Title)
		assert.Equal(t, 1, len(titleswhole))
		//assert.Equal(t, titleString4, titleswhole[0].Title)
	})

	t.Run("Query by billnumberversion", func(t *testing.T) {
		var bill2 Bill
		querystring := "117hr200ih"
		db.First(&bill2, "Billnumberversion = ?", querystring) // find item with billnumber = 117hr200
		db.Model(&bill2).Association("Titles").Find(&titles)
		log.Debug().Msg("Got no associated title item.")
		assert.Equal(t, 0, len(titles))
	})

	t.Run("Add Title using AddTitleDB", func(t *testing.T) {
		AddTitleDb(db, titleString3)
		var newTitle Title
		db.First(&newTitle, "title = ?", titleString3) // find title
		assert.Equal(t, titleString3, newTitle.Title)
	})
	t.Run("Add bills using AddBillnumberversionsDb", func(t *testing.T) {
		AddBillnumberversionsDb(db, []string{"115hr100ih", "116hr999rh"})
		var bill3 Bill
		var bill4 Bill
		db.First(&bill3, "Billnumber = ?", "115hr100") // find bill
		assert.Equal(t, "115hr100ih", bill3.Billnumberversion)
		db.First(&bill4, "Billnumberversion = ?", "116hr999rh") // find bill
		assert.Equal(t, "116hr999", bill4.Billnumber)
	})

	t.Run("Add bill by Struct (with titles)", func(t *testing.T) {
		var title5 Title
		var associatedBills []*Bill
		var associatedTitles []*Title
		bill5 := Bill{Billnumber: "111hr100", Billnumberversion: "111hr100ih", Titles: []*Title{{Title: titleString5}, {Title: titleString5 + "2"}}}
		AddBillStructDb(db, &bill5)
		// Equivalent to:
		//db.Create(&bill5)
		log.Debug().Msgf("Bill '%v' added, with titles '%s' and '%s'", bill5.Billnumberversion, titleString5, titleString5+"2")

		var newbill Bill
		db.First(&newbill, "Billnumberversion = ?", "111hr100ih") // find bill with billnumberversion = 111hr100ih

		// Should be associated with title5 and title5+2
		db.Model(newbill).Association("Titles").Find(&associatedTitles)
		assert.Equal(t, 2, len(associatedTitles))
		log.Info().Msgf("Titles for bill '%s': %+v", newbill.Billnumberversion, associatedBills)

		db.First(&title5, "title = ?", titleString5) // find title with Title = titleString5
		// Should be associated with bill5
		db.Model(title5).Association("Bills").Find(&associatedBills)
		assert.Equal(t, 1, len(associatedBills))
		log.Info().Msgf("Bills with title '%s': %+v", titleString5, associatedBills)
	})

	t.Run("Remove Title by string", func(t *testing.T) {
		var newTitle1 Title
		var newTitle2 Title
		err := db.Where("title = ?", titleString3).First(&newTitle1).Error
		assert.Nil(t, err)
		RemoveTitleDb(db, titleString3)
		err = db.Where("title = ?", titleString3).First(&newTitle2).Error
		assert.NotNil(t, err)
	})

	t.Run("Get bills using GetBillsWithSameTitleDb", func(t *testing.T) {
		querystring := "116hr1500"
		bills, bills_whole, err := GetBillsWithSameTitleDb(db, querystring)
		log.Warn().Msgf("Bills: %d", len(bills))
		log.Warn().Msgf("Bills whole: %d", len(bills_whole))
		assert.Nil(t, err)
		log.Debug().Msgf("Related bills for %s: %+v", querystring, bills)
		log.Debug().Msgf("Related bills (whole) for %s: %+v", querystring, bills_whole)
	})

	t.Run("Get bills using GetBillsWithSameTitleDb (only same bill)", func(t *testing.T) {
		querystring := "117hr200"
		bills, bills_whole, err := GetBillsWithSameTitleDb(db, querystring)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(bills))
		assert.Equal(t, querystring, bills[0].Billnumber)
		log.Debug().Msgf("Related bills for %s: %+v", querystring, bills)
		log.Debug().Msgf("Related bills (whole) for %s: %+v", querystring, bills_whole)
	})

	// Delete - delete all items in bill and title tables
	db.Exec("DELETE FROM bills")
	db.Exec("DELETE FROM titles")
}
