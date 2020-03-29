package praytime

import (
	"log"
	"regexp"
	"strconv"
)

// GetDate change string to date
// return date, month, err
func GetDate(t string) (string, string) {
	lt := len(t)
	month := t[lt-3 : lt]
	day := t[lt-6 : lt-4]
	day = getNumberOnly(day)

	// conv to int if pos
	_, err := strconv.Atoi(day)
	if err != nil {
		log.Fatal(err)
	}

	return day, month
}

func getNumberOnly(s string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
		return ""
	}

	result := reg.ReplaceAllString(s, "")
	return result
}
