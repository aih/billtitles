package billtitles

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Title struct {
	gorm.Model
	Name  string      `gorm:"index:,unique"`
	Bills []*BillItem `gorm:"many2many:bill_titles;"`
}

func MakeBillAndTitle() {

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
	db.First(&product, 1)                 // find product with integer primary key
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
}
