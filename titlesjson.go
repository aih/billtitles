package billtitles

import (
	"errors"
	"io/ioutil"
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
	jsonFile, err := os.Open(titlePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("Successfully opened Titles JSON file")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	titleMap, err = UnmarshalJSON(byteValue)

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
