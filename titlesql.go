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
	BillCongressTypeNumber        string   `gorm:"not null" json:"bill_congress_type_number"`
	BillCongressTypeNumberVersion string   `gorm:"index:,unique" json:"bill_congress_type_number_version"`
	Titles                        []*Title `gorm:"many2many:bill_titles;" json:"titles"`
	TitlesWholeBill               []*Title `gorm:"many2many:bill_titleswhole" json:"titles_whole_bill"`
}
