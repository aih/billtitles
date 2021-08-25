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
	log.Info().Msg("Test setting and getting titles and bills from sql db")

	const (
		titleString  = "This is a test title"
		titleString2 = "This is a new test title"
		titleString3 = "This Title Added with AddTitleDB function"
		titleString4 = "This is a title for a whole bill"
	)
	newBill1 := &Bill{Billnumberversion: "116hr1500ih", Billnumber: "116hr1500"}
	newBill2 := &Bill{Billnumberversion: "117hr200ih", Billnumber: "117hr200"}
	newBill3 := &Bill{Billnumberversion: "117hr100ih", Billnumber: "117hr100"}
	newBill4 := &Bill{Billnumberversion: "117hr222ih", Billnumber: "117hr222"}
	newTitle := &Title{Title: titleString, Bills: []*Bill{newBill1, newBill3}}
	db.Session(&gorm.Session{FullSaveAssociations: true})

	// Create
	db.Create(newTitle)
	db.Create(newBill2)

	db.Model(newTitle).Association("Bills").Append([]Bill{*newBill4})

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
		db.Model(title).Association("Bills").Find(&associatedBills)
		assert.NotEqual(t, 0, len(associatedBills))
		log.Debug().Msgf("Associated bills: %+v", associatedBills)
		assert.Equal(t, 3, len(associatedBills))
	})

	t.Run("Get Associated bills using GetBillsByTitleDb", func(t *testing.T) {
		associatedBills := GetBillsByTitleDb(db, titleString)
		assert.NotEqual(t, 0, len(associatedBills))
		log.Debug().Msgf("Associated bills: %+v", associatedBills)
		assert.Equal(t, 3, len(associatedBills))
	})

	t.Run("Test GetTitlesByBillnumberDb", func(t *testing.T) {
		titles := GetTitlesByBillnumberDb(db, "116hr1500")
		assert.NotEqual(t, 0, len(titles))
		log.Debug().Msgf("Associated titles: %+v", titles)
		assert.Equal(t, 1, len(titles))
		assert.Equal(t, titleString, titles[0].Title)
	})

	t.Run("Test GetTitlesByBillnumberVersionDb", func(t *testing.T) {
		titles := GetTitlesByBillnumberVersionDb(db, "117hr100ih")
		assert.NotEqual(t, 0, len(titles))
		log.Debug().Msgf("Associated titles: %+v", titles)
		assert.Equal(t, 1, len(titles))
		assert.Equal(t, titleString, titles[0].Title)
	})

	t.Run("Test GetTitlesWholeByBillnumberVersionDb", func(t *testing.T) {
		titles := GetTitlesByBillnumberVersionDb(db, "117hr100ih")
		assert.NotEqual(t, 0, len(titles))
		log.Debug().Msgf("Associated titles: %+v", titles)
		assert.Equal(t, 1, len(titles))
		assert.Equal(t, titleString, titles[0].Title)
	})

	t.Run("Add a title entry with db.Model", func(t *testing.T) {
		// Update - update title
		db.Model(&title).Update("Title", titleString2)
		log.Debug().Msgf("Title '%v' updated", title.Title)
		var sampleBill Bill
		var associatedTitles []*Title
		db.Take(&sampleBill)
		assert.NotNil(t, sampleBill)
		db.Model(&sampleBill).Association("Titles").Find(&associatedTitles)
		assert.NotNil(t, associatedTitles)
		assert.NotEqual(t, 0, len(associatedTitles))
	})

	t.Run("Add a title (whole) entry with db.Model", func(t *testing.T) {
		newTitle4 := &Title{Title: titleString4, Bills: []*Bill{newBill1, newBill3}}
		db.Create(&newTitle4)
		log.Debug().Msgf("Title '%v' added", newTitle4.Title)
		db.Model(&newBill1).Association("TitlesWhole").Append(newTitle4)
		var associatedTitles []*Title
		db.Model(&newBill1).Association("TitlesWhole").Find(&associatedTitles)
		log.Debug().Msgf("Title '%s' associated with %s", associatedTitles[0].Title, newBill1.Billnumber)
		assert.NotNil(t, associatedTitles)
		assert.NotEqual(t, 0, len(associatedTitles))
	})

	t.Run("Get bills using GetBillsWithSameTitleDb", func(t *testing.T) {
		querystring := "116hr1500"
		bills, bills_whole, err := GetBillsWithSameTitleDb(db, querystring)
		assert.Nil(t, err)
		log.Debug().Msgf("Related bills for %s: %+v", querystring, bills)
		log.Debug().Msgf("Related bills (whole) for %s: %+v", querystring, bills_whole)
	})

	t.Run("Get bills using GetBillsWithSameTitleDb (fewer results)", func(t *testing.T) {
		querystring := "117hr200"
		bills, bills_whole, err := GetBillsWithSameTitleDb(db, querystring)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(bills))
		assert.Equal(t, querystring, bills[0].Billnumber)
		log.Debug().Msgf("Related bills for %s: %+v", querystring, bills)
		log.Debug().Msgf("Related bills (whole) for %s: %+v", querystring, bills_whole)
	})

	var bill Bill
	var bill2 Bill
	var titles []Title
	//bctnvs := []string{"116hr1500ih", "116hr200rh", "115hr2200ih"}
	//bctns := []string{"116hr1500", "116hr200", "115hr2200"}

	t.Run("Query by ID", func(t *testing.T) {
		db.First(&bill, "ID = ?", "1") // find item with any of the billnumbers in bctnvs
		log.Debug().Msgf("Got bill item for ID=1: %+v", bill)
		assert.Equal(t, "116hr1500ih", bill.Billnumberversion)
	})

	t.Run("Query by billnumber (no titles)", func(t *testing.T) {
		querystring := "117hr200"
		db.First(&bill2, "Billnumber = ?", querystring) // find item with billnumber = 117hr200
		log.Debug().Msgf("Got bill item for billnumber=%s: %+v", querystring, bill2)
		assert.Equal(t, "117hr200ih", bill2.Billnumberversion)
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
		log.Debug().Msgf("Associated titles (whole): %+v", titleswhole[0].Title)
		assert.Equal(t, 1, len(titleswhole))
		assert.Equal(t, titleString4, titleswhole[0].Title)
	})

	t.Run("Query by billnumberversion", func(t *testing.T) {
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

	t.Run("Remove Title by string", func(t *testing.T) {
		var newTitle1 Title
		var newTitle2 Title
		err := db.Where("title = ?", titleString3).First(&newTitle1).Error
		assert.Nil(t, err)
		RemoveTitleDb(db, titleString3)
		err = db.Where("title = ?", titleString3).First(&newTitle2).Error
		assert.NotNil(t, err)
	})

	// Delete - delete all items in bill and title tables
	db.Exec("DELETE FROM bills")
	db.Exec("DELETE FROM titles")
}
