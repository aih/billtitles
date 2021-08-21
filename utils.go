package billtitles

import "strings"

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
