package billtitles

import (
	"testing"

	"github.com/aih/billtitles/internal/testutils"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const sampleTitle = "Expressing the sense of the House of Representatives that the United States remains committed to the North Atlantic Treaty Organization (NATO)."

func TestLoadTitles(t *testing.T) {
	testutils.SetLogLevel()
	log.Info().Msg("Test loading titles from json")
	_, error := LoadTitlesMap(SampleTitlePath)
	assert.Nil(t, error)
}
func TestGetTitle(t *testing.T) {
	testutils.SetLogLevel()
	log.Info().Msg("Test getting title from loaded json")
	titleMap, _ := LoadTitlesMap(SampleTitlePath)
	billnumbers, err := GetBillnumbersByTitle(titleMap, sampleTitle)
	assert.Nil(t, err)
	assert.Equal(t, billnumbers, []string{"111hres152"})
}
func TestAddBillNumbersToTitle(t *testing.T) {
	testutils.SetLogLevel()
	log.Info().Msg("Test adding a sample title ")
	titleMap, _ := LoadTitlesMap(SampleTitlePath)
	AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr222"})
	AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr999"})
	newbillnumbers, err := GetBillnumbersByTitle(titleMap, "A test title")
	assert.Nil(t, err)
	assert.Equal(t, newbillnumbers, []string{"118hr222", "118hr999"})
}

func TestRemoveTitle(t *testing.T) {
	testutils.SetLogLevel()
	log.Info().Msg("Test removing a title ")
	titleMap, _ := LoadTitlesMap(SampleTitlePath)
	AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr222"})
	AddBillNumbersToTitle(titleMap, "A test title", []string{"118hr999"})
	_, err := GetBillnumbersByTitle(titleMap, "A test title")
	assert.Nil(t, err)
	RemoveTitle(titleMap, "A test title")
	_, err = GetBillnumbersByTitle(titleMap, "A test title")
	assert.NotNil(t, err)
}
