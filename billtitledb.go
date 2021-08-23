package billtitles

import (
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func RunDbExample() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	//db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	log.Info().Msgf("product %v updated", product)
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	log.Info().Msgf("product %v updated", product)

	// Delete - delete product
	db.Delete(&product, 1)

	newLogger := logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db2, err2 := gorm.Open(sqlite.Open("billtitles.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err2 != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db2.AutoMigrate(&Bill{}, &Title{})

	newBill1 := &Bill{BillCongressTypeNumberVersion: "116hr1500ih", BillCongressTypeNumber: "116hr1500"}
	newBill2 := &Bill{BillCongressTypeNumberVersion: "117hr200ih", BillCongressTypeNumber: "117hr200"}
	newBill3 := &Bill{BillCongressTypeNumberVersion: "117hr100ih", BillCongressTypeNumber: "117hr100"}
	newBill4 := &Bill{BillCongressTypeNumberVersion: "117hr222ih", BillCongressTypeNumber: "117hr222"}
	newTitle := &Title{Title: "This is a test title", Bills: []*Bill{newBill1, newBill3}}
	db2.Session(&gorm.Session{FullSaveAssociations: true})

	// Create
	db2.Create(newTitle)
	db2.Create(newBill2)
	db2.Save(newTitle)
	db2.Save(newBill2)

	db2.Model(newTitle).Association("Bills").Append([]Bill{*newBill4})
	//db2.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&newTitle)

	// Read
	var title Title
	var associatedBills []*Bill
	//db2.First(&title, 1)                                  // find product with integer primary key
	db2.First(&title, "title = ?", "This is a test title") // find product with code D42
	log.Info().Msgf("Got title item: %v", title)
	log.Info().Msgf("Title item has the title: %v", title.Title)
	db2.Model(title).Association("Bills").Find(&associatedBills)
	if len(associatedBills) == 0 {
		log.Info().Msg("Got no associated bill item.")
	} else {
		for _, associatedBill := range associatedBills {
			log.Info().Msgf("Got bill associated with the title item: %s", associatedBill.BillCongressTypeNumberVersion)
		}
	}

	// Update - update title
	db2.Model(&title).Update("Title", "This is a new test title")
	log.Info().Msgf("Title '%v' updated", title.Title)

	var sampleBill Bill
	var associatedTitles []*Title
	db2.Take(&sampleBill)
	log.Info().Msgf("Got sample bill item: %v", sampleBill.BillCongressTypeNumberVersion)
	db2.Model(&sampleBill).Association("Titles").Find(&associatedTitles)
	if len(associatedTitles) == 0 {
		log.Info().Msg("Got no associated title item.")
	} else {
		for _, associatedTitle := range associatedTitles {
			log.Info().Msgf("Got title associated with the bill item: %s", associatedTitle.Title)
		}
	}

	var bill Bill
	var titles []Title
	//bctnvs := []string{"116hr1500ih", "116hr200rh", "115hr2200ih"}
	//bctns := []string{"116hr1500", "116hr200", "115hr2200"}

	db2.First(&bill, "ID = ?", "1") // find item with any of the billnumbers in bctnvs
	log.Info().Msgf("Got bill item: %+v", bill)

	// TODO: Getting error 'No such column: BillCongressTypeNumber'
	db2.First(&bill, "billcongresstypenumber = ?", "117hr200") // find item with billnumber = 117hr200
	log.Info().Msgf("Got bill item: %+v", bill)

	db2.Model(&bill).Association("Titles").Find(&titles)
	if len(titles) == 0 {
		log.Info().Msg("Got no associated title item.")
	} else {
		for _, titleItem := range titles {
			log.Info().Msgf("Got title associated with the bill item: %s", titleItem.Title)
		}
	}

	// Delete - delete title
	db2.Delete(&title, 1)

}
