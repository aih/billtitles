package billtitles

type RelatedBillItem struct {
	BillId                        string   `json:"bill_id"`
	IdentifiedBy                  string   `json:"identified_by"`
	Reason                        string   `json:"reason"`
	Type                          string   `json:"type"`
	BillCongressTypeNumber        string   `json:"bill_congress_type_number"`
	BillCongressTypeNumberVersion string   `json:"bill_congress_type_number_version"`
	Titles                        []string `json:"titles"`
	TitlesWholeBill               []string `json:"titles_whole_bill"`
}

// Used for db functions
type BillItem struct {
	BillId                        string   `json:"bill_id"`
	IdentifiedBy                  string   `json:"identified_by"`
	Reason                        string   `json:"reason"`
	Type                          string   `json:"type"`
	BillCongressTypeNumber        string   `gorm:"not null" json:"bill_congress_type_number"`
	BillCongressTypeNumberVersion string   `gorm:"primary_key" json:"bill_congress_type_number_version"`
	Titles                        []*Title `gorm:"many2many:bill_titles;" json:"titles"`
	TitlesWholeBill               []*Title `gorm:"many2many:bill_titleswhole" json:"titles_whole_bill"`
}

type RelatedBillMap map[string]RelatedBillItem
