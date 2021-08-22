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

type RelatedBillMap map[string]RelatedBillItem
