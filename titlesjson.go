package billtitles

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

// Opens maintitle file and has functions to add titles and billnumbers

func AddTitle(titleMap TitleMap, title string) (TitleMap, error) {
	titleMap[title] = make([]string, 0)
	return titleMap, nil
}

func RemoveTitle(titleMap TitleMap, title string) (TitleMap, error) {
	_, ok := titleMap[title]
	if ok {
		delete(titleMap, title)
	}
	return titleMap, nil
}

func GetTitle(titleMap TitleMap, title string) ([]string, error) {
	_, ok := titleMap[title]
	if ok {
		return titleMap[title], nil
	} else {
		return nil, errors.New("Title not found")
	}
}

func AddBillNumbersToTitle(titleMap TitleMap, title string, billnumbers []string) (TitleMap, error) {
	_, ok := titleMap[title]
	if ok {
		titleMap[title] = RemoveDuplicates(append(titleMap[title], billnumbers...))
	} else {
		titleMap[title] = billnumbers
	}
	return titleMap, nil

}

func LoadTitlesMap(titlePath string) (titleMap TitleMap, err error) {
	if _, err := os.Stat(titlePath); os.IsNotExist(err) {
		titlePath = TitlesPath
		if _, err := os.Stat(titlePath); os.IsNotExist(err) {
			return nil, errors.New("titles file file not found")
		}
	}
	jsonFile, err := os.Open(titlePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("Successfully opened Titles JSON file")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &titleMap)

	if err != nil {
		return nil, err
	} else {
		return titleMap, nil
	}
}

func SaveTitlesMap(titleMap TitleMap, titlePath string) (err error) {
	jsonFile, err := os.Create(titlePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	jsonByte, err := json.Marshal(titleMap)
	if err != nil {
		return err
	}
	jsonFile.Write(jsonByte)
	return nil
}
