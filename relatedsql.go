package billtitles

import (
	"encoding/json"
	"errors"
	"io/fs"
	stdlog "log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const BILLSRELATED_DB = "billsrelated.db"

type BillToBill struct {
	gorm.Model
	Billnumber    string `gorm:"index:,not null" json:"billnumber"`
	Billnumber_to string `gorm:"index:,not null" json:"billnumber_to"`
	Reason        string `gorm:"not null" json:"reason"`
	Identified_by string `gorm:"index:,not null" json:"identified_by"`
}

type compareItem struct {
	Score        float64 `json:"score"`
	ScoreOther   float64 `json:"score_other"` // score for other bill
	Explanation  string  `json:"explanation"`
	ComparedDocs string  `json:"compared_docs"`
}

type filterFunc func(string) bool

// Walk directory with a filter. Returns the filepaths that
// pass the 'testPath' function
// There is an exported function in the `bills` package that does this
func walkDirFilter(root string, testPath filterFunc) (filePaths []string, err error) {
	defer log.Info().Msg("Done collecting filepaths.")
	log.Info().Msgf("Getting all file paths in %s.  This may take a while.\n", root)
	filePaths = make([]string, 0)
	accumulate := func(fpath string, entry fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Err(err)
			return err
		}
		if testPath(fpath) {
			filePaths = append(filePaths, fpath)
		}
		return nil
	}
	err = filepath.WalkDir(root, accumulate)
	return
}

func GetRelatedDb(dbname string) *gorm.DB {
	if dbname == "" {
		dbname = BILLTITLES_DB
	}
	var newLogger = logger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	var db, _ = gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: newLogger,
	})
	// May not be necessary; applies all associated changes
	// when updating a title or bill.
	db.AutoMigrate(&BillToBill{})
	db.Session(&gorm.Session{FullSaveAssociations: true})
	return db
}

func LoadRelatedMap(jsonPath string) (*sync.Map, error) {
	relatedMap := new(sync.Map)
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		jsonPath = TitlesPath
		if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
			return relatedMap, errors.New("related bills file file not found")
		}
	}
	log.Debug().Msgf("Path to JSON file: %s", jsonPath)
	var err error
	relatedMap, err = UnmarshalJsonFile(jsonPath)
	if err != nil {
		return nil, err
	} else {
		return relatedMap, nil
	}
}

//TODO create related sync.Map from individual json files

var similarCategoryJsonFilter = func(testPath string) bool {
	matched, err := regexp.MatchString(`esSimilarCategory`, testPath)
	matchedJson, err2 := regexp.MatchString(`\.json$`, testPath)
	if err != nil || err2 != nil {
		return false
	}
	return matched && matchedJson
}

func processRelatedJson(filePath string, similarityChannel chan compareItem, sem chan bool, wg *sync.WaitGroup) error {
	defer func() {
		log.Info().Msgf("Finished processing: %s\n", filePath)
		<-sem
	}()
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Msgf("Error reading data.json: %s", err)
		return err
	}
	var dat compareItem
	_ = json.Unmarshal([]byte(file), &dat)

	// TODO: handle dat and put it into the channel
}

// jsonPath := RelatedJsonPath
// db := GetDb(BILLSRELATED_DB)
func LoadBillsRelatedToDBFromJson(db *gorm.DB, parentPath string) {
	log.Info().Msgf("Loading titles to database from json files in directory: %s", parentPath)
	defer log.Info().Msg("Done processing similar bill json")
	dataJsonFiles, error := walkDirFilter(parentPath, similarCategoryJsonFilter)
	if error != nil {
		log.Fatal().Msgf("Error getting files list: %s", error)
	}
	maxopenfiles := 100
	sem := make(chan bool, maxopenfiles)
	similarityChannel := make(chan compareItem)
	wg := &sync.WaitGroup{}
	wg.Add(len(dataJsonFiles))
	go func() {
		wg.Wait()
		close(similarityChannel)
	}()

	go func() {
		for range dataJsonFiles {
			compare := <-similarityChannel
			log.Debug().Msgf("Got compare item from Channel: %v\n", compare)
			// TODO: insert into db
		}
	}()

	for _, jpath := range dataJsonFiles {
		sem <- true
		go processRelatedJson(jpath, similarityChannel, sem, wg)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}
