package billtitles

import (
	"errors"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

// Opens maintitle file and has functions to add titles and billnumbers

func AddTitle(titleMap *sync.Map, title string) (*sync.Map, error) {
	titleMap.LoadOrStore(title, make([]string, 0))
	return titleMap, nil
}

func RemoveTitle(titleMap *sync.Map, title string) (*sync.Map, error) {
	titleMap.Delete(title)
	return titleMap, nil
}

func GetBillnumbersByTitle(titleMap *sync.Map, title string) (billnumbers []string, err error) {
	results, ok := titleMap.Load(title)
	if ok {
		billnumbers := results.([]string)
		return billnumbers, nil
	} else {
		return nil, errors.New("Title not found")
	}
}

func AddBillNumbersToTitle(titleMap *sync.Map, title string, billnumbers []string) (*sync.Map, error) {
	if titleBills, loaded := titleMap.LoadOrStore(title, billnumbers); loaded {
		titleBills = RemoveDuplicates(append(titleBills.([]string), billnumbers...))
		titleMap.Store(title, titleBills)
	}
	return titleMap, nil

}

func LoadTitlesMap(titlePath string) (*sync.Map, error) {
	titleMap := new(sync.Map)
	if _, err := os.Stat(titlePath); os.IsNotExist(err) {
		titlePath = TitlesPath
		if _, err := os.Stat(titlePath); os.IsNotExist(err) {
			return titleMap, errors.New("titles file file not found")
		}
	}
	log.Debug().Msgf("Path to JSON file: %s", titlePath)
	var err error
	titleMap, err = UnmarshalTitlesJsonFile(titlePath)
	if err != nil {
		return nil, err
	} else {
		return titleMap, nil
	}
}

func SaveTitlesMap(titleMap *sync.Map, titlePath string) (err error) {
	jsonFile, err := os.Create(titlePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	jsonByte, err := MarshalJSONStringArray(titleMap)
	if err != nil {
		return err
	}
	jsonFile.Write(jsonByte)
	return nil
}

func MakeSampleTitlesFile(titleMap *sync.Map) {
	defer func() { log.Info().Msg("done making samples file") }()
	sampleTitles := new(sync.Map)
	count := 0
	titleMap.Range(func(key, value interface{}) bool {
		count += 1
		if count > 4 {
			log.Debug().Msgf("Returning 'false' from range loop")
			return false
		}
		log.Debug().Msgf("Adding a title")
		sampleTitles.Store(key.(string), value.([]string))
		return true
	})
	SaveTitlesMap(sampleTitles, "data/sampletitles.json")
}
