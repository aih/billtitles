package billtitles

import (
	"gorm.io/gorm"
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
