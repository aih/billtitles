package billtitles

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

// Removes duplicates in a list of strings
// Returns the deduplicated list
// Trims leading and trailing space for each element
func RemoveDuplicates(elements []string) []string { // change string to int here if required
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{} // change string to int here if required
	result := []string{}             // change string to int here if required

	for v := range elements {
		currentElement := strings.TrimSpace(elements[v])
		if currentElement == "" || encountered[currentElement] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[currentElement] = true
			// Append to result slice.
			result = append(result, currentElement)
		}
	}
	// Return the new slice.
	return result
}

// Marshals a sync.Map object of the type map[string][]string
// see https://stackoverflow.com/a/46390611/628748
// and https://stackoverflow.com/a/65442862/628748
func MarshalJSONStringArray(m *sync.Map) ([]byte, error) {
	tmpMap := make(map[string][]string)
	m.Range(func(k interface{}, v interface{}) bool {
		tmpMap[k.(string)] = v.([]string)
		return true
	})
	return json.Marshal(tmpMap)
}

func UnmarshalJSON(data []byte) (*sync.Map, error) {
	var tmpMap map[interface{}]interface{}
	m := &sync.Map{}
    
	if err := json.Unmarshal(data, &tmpMap); err != nil {
	    return m, err
	}
    
	for key, value := range tmpMap {
	    m.Store(key, value)
	}
	return m, nil
}

// Gets keys of a sync.Map
func GetSyncMapKeys(m *sync.Map) (s string) {
	m.Range(func(k, v interface{}) bool {
		log.Info().Msgf("%s: %v", k, v)
		if s != "" {
			s += ", "
		}
		s += k.(string)
		return true
	})
	return
}
