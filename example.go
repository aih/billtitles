package billtitles

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func RunExample() {

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

	db2, err := gorm.Open(sqlite.Open("billtitles.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db2.AutoMigrate(&BillItem{}, &Title{})

	// Create
	db2.Create(&Title{Name: "This is a test title"})

	// Read
	var title Title
	//db2.First(&title, 1)                                  // find product with integer primary key
	db2.First(&title, "name = ?", "This is a test title") // find product with code D42
	log.Info().Msgf("Got title item: %v", title)
	log.Info().Msgf("Got title item: %v", title.Name)

	// Update - update title
	db2.Model(&title).Update("Name", "This is a new test title")
	log.Info().Msgf("Title %v updated", title)

	// Delete - delete title
	db2.Delete(&title, 1)

}
