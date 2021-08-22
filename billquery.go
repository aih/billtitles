package billtitles

import (
	"sync"

	"github.com/rs/zerolog/log"
)

func GetRelatedBills(titleMap *sync.Map, titleWholeBillMap *sync.Map, relatedBillItem RelatedBillItem) (relatedBills RelatedBillMap, err error) {
	titles := relatedBillItem.Titles
	titlesWholeBill := relatedBillItem.TitlesWholeBill
	relatedBills = make(RelatedBillMap)
	for _, title := range titles {
		log.Debug().Msgf("Getting bills with title: %s", title)
		billnumbers, err := GetBillnumbersByTitle(titleMap, title)
		if err != nil {
			log.Error().Msgf("Error getting bills with title: %s", title)
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
		log.Debug().Msgf("Getting bills with main title: %s", title)
		billnumbers, err := GetBillnumbersByTitle(titleWholeBillMap, title)
		if err != nil {
			log.Error().Msgf("Error getting bills with title: %s", title)
			continue
		}
		for _, billnumber := range billnumbers {
			if billItem, ok := relatedBills[billnumber]; ok {
				tempBillItem := billItem
				tempBillItem.TitlesWholeBill = RemoveDuplicates(append(billItem.TitlesWholeBill, titlesWholeBill...))
				relatedBills[billnumber] = tempBillItem
			} else {
				relatedBills[billnumber] = RelatedBillItem{TitlesWholeBill: []string{title}}
			}
		}
	}
	return relatedBills, err
}
