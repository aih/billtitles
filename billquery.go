package billtitles

import "sync"

func GetRelatedBills(titleMap *sync.Map, relatedBillItem RelatedBillItem) (relatedBills RelatedBillMap, err error) {
	titles := relatedBillItem.Titles
	titlesWholeBill := relatedBillItem.TitlesWholeBill
	relatedBills = make(RelatedBillMap)
	for _, title := range titles {
		billnumbers, err := GetBillnumbersByTitle(titleMap, title)
		if err != nil {
			continue
		}
		for _, billnumber := range billnumbers {
			if billItem, ok := relatedBills[billnumber]; ok {
				tempBillItem := billItem
				tempBillItem.Titles = RemoveDuplicates(append(billItem.Titles, title))
				relatedBills[billnumber] = tempBillItem
			} else {
				relatedBills[billnumber] = RelatedBillItem{Titles: []string{title}}
			}
		}
	}
	for _, title := range titlesWholeBill {
		billnumbers, err := GetBillnumbersByTitle(titleMap, title)
		if err != nil {
			continue
		}
		for _, billnumber := range billnumbers {
			if billItem, ok := relatedBills[billnumber]; ok {
				tempBillItem := billItem
				tempBillItem.TitlesWholeBill = RemoveDuplicates(append(billItem.Titles, titlesWholeBill...))
				relatedBills[billnumber] = tempBillItem
			} else {
				relatedBills[billnumber] = RelatedBillItem{Titles: []string{title}}
			}
		}
	}
	return relatedBills, err
}
