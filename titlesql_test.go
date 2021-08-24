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
	testutils.SetLogLevel()
	log.Info().Msg("Test setting and getting titles and bills from sql db")

	newBill1 := &Bill{Billnumberversion: "116hr1500ih", Billnumber: "116hr1500"}
	newBill2 := &Bill{Billnumberversion: "117hr200ih", Billnumber: "117hr200"}
	newBill3 := &Bill{Billnumberversion: "117hr100ih", Billnumber: "117hr100"}
	newBill4 := &Bill{Billnumberversion: "117hr222ih", Billnumber: "117hr222"}
	newTitle := &Title{Title: "This is a test title", Bills: []*Bill{newBill1, newBill3}}
	db.Session(&gorm.Session{FullSaveAssociations: true})

	// Create
	db.Create(newTitle)
	db.Create(newBill2)
	db.Save(newTitle)
	db.Save(newBill2)

	db.Model(newTitle).Association("Bills").Append([]Bill{*newBill4})
	//db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&newTitle)

	// Read
	var title Title
	var associatedBills []*Bill

	t.Run("Get Title", func(t *testing.T) {
		db.First(&title, "title = ?", "This is a test title") // find title
		log.Debug().Msgf("Got title item: %v", title)
		log.Debug().Msgf("Title item has the title: %v", title.Title)
		assert.Equal(t, "This is a test title", title.Title)
	})
	t.Run("Get Associated Bills", func(t *testing.T) {
		db.First(&title, "title = ?", "This is a test title") // find title with Title = "This is a test title"
		db.Model(title).Association("Bills").Find(&associatedBills)
		assert.NotEqual(t, 0, len(associatedBills))
		log.Debug().Msgf("Associated bills: %+v", associatedBills)
		assert.Equal(t, 3, len(associatedBills))
	})

	t.Run("Add a title entry with db.Model", func(t *testing.T) {
		// Update - update title
		db.Model(&title).Update("Title", "This is a new test title")
		log.Debug().Msgf("Title '%v' updated", title.Title)
		var sampleBill Bill
		var associatedTitles []*Title
		db.Take(&sampleBill)
		assert.NotNil(t, sampleBill)
		db.Model(&sampleBill).Association("Titles").Find(&associatedTitles)
		assert.NotNil(t, associatedTitles)
		assert.NotEqual(t, 0, len(associatedTitles))
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

	t.Run("Query by billnumber", func(t *testing.T) {
		db.First(&bill2, "Billnumber = ?", "117hr200") // find item with billnumber = 117hr200
		log.Debug().Msgf("Got bill item for billnumber=117hr200: %+v", bill2)
		assert.Equal(t, "117hr200ih", bill2.Billnumberversion)
	})

	t.Run("Query by billnumberversion", func(t *testing.T) {
		db.Model(&bill2).Association("Titles").Find(&titles)
		log.Debug().Msg("Got no associated title item.")
		assert.Equal(t, 0, len(titles))
	})

	t.Run("Add Title using AddTitleDB", func(t *testing.T) {
		newTitle := "This Title Added with AddTitleDB function"
		AddTitleDb(db, newTitle)
		var newTitle2 Title
		db.First(&newTitle2, "title = ?", newTitle) // find title
		assert.Equal(t, newTitle, newTitle2.Title)
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

	// Delete - delete all items in bill and title tables
	db.Exec("DELETE FROM bills")
	db.Exec("DELETE FROM titles")
}
