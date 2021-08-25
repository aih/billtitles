package billtitles

import (
	"testing"

	"github.com/aih/billtitles/internal/testutils"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const sampleTitle1 = "Expressing the sense of the House of Representatives that the United States remains committed to the North Atlantic Treaty Organization (NATO)."

func TestLoadTitles(t *testing.T) {
	testutils.SetLogLevel()
	log.Debug().Msg("Test loading titles from json")
	_, error := LoadTitlesMap(SampleTitlesPath)
	assert.Nil(t, error)
}
func TestGetTitle(t *testing.T) {
	testutils.SetLogLevel()
	log.Debug().Msg("Test getting title from loaded json")
	titleMap, _ := LoadTitlesMap(SampleTitlesPath)
	billnumbers, err := GetBillnumbersByTitle(titleMap, sampleTitle1)
	assert.Nil(t, err)
	assert.Equal(t, billnumbers, []string{"111hres152"})
}
func TestAddBillNumbersToTitle(t *testing.T) {
	testutils.SetLogLevel()
	log.Debug().Msg("Test adding a sample title ")
	titleMap, _ := LoadTitlesMap(SampleTitlesPath)
	AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr222"})
	AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr999"})
	newbillnumbers, err := GetBillnumbersByTitle(titleMap, "A test title")
	assert.Nil(t, err)
	assert.Equal(t, newbillnumbers, []string{"118hr222", "118hr999"})
}

func TestRemoveTitle(t *testing.T) {
	testutils.SetLogLevel()
	log.Debug().Msg("Test removing a title ")
	titleMap, _ := LoadTitlesMap(SampleTitlesPath)
	AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr222"})
	AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr999"})
	_, err := GetBillnumbersByTitle(titleMap, "A test title")
	assert.Nil(t, err)
	RemoveTitle(titleMap, "A test title")
	_, err = GetBillnumbersByTitle(titleMap, "A test title")
	assert.NotNil(t, err)
}
