package billtitles

import (
	"testing"

	"github.com/aih/billtitles/internal/testutils"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const (
	sampleTitle = "21st Century Energy Workforce Act"
	// declared in titlesjson_test.go; sampleTitle1 = "Expressing the sense of the House of Representatives that the United States remains committed to the North Atlantic Treaty Organization (NATO)."
	sampleTitle2 = "To allow the Miami Tribe of Oklahoma to lease or transfer certain lands."
)

var sampleTitles = []string{sampleTitle1, sampleTitle2}
var sampleWholeTitles = []string{sampleTitle}

func TestGetRelatedBills(t *testing.T) {
	testutils.SetLogLevel()
	log.Debug().Msg("Test querying title maps for related bills")
	titleMap, error := LoadTitlesMap(TitlesPath)
	if error != nil {
		log.Info().Msgf("%v", error)
		panic(error)
	}
	titleWholeBillMap, error := LoadTitlesMap(MainTitlesPath)
	if error != nil {
		log.Info().Msgf("%v", error)
		panic(error)
	}
	billItem := RelatedBillItem{
		BillId:          "1500-116hr",
		Titles:          sampleTitles,
		TitlesWholeBill: sampleWholeTitles,
	}
	relatedBills, err := GetRelatedBills(titleMap, titleWholeBillMap, billItem)
	if err != nil {
		log.Info().Msgf("Error getting related bills: %v", err)
	}
	log.Debug().Msgf("%v", relatedBills)

	assert.Nil(t, err)
	relatedBill, ok := relatedBills["115hr1837"]
	assert.True(t, ok)
	assert.Equal(t, relatedBill.TitlesWholeBill, []string{"21st Century Energy Workforce Act"})
}
