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

// new get date func to convert from new api's readable format
func GetDateFromReadable(t string) (string, string) {
	lt := len(t)
	month := t[lt-8 : lt-5]
	day := t[lt-11 : lt-9]
	day = getNumberOnly(day)

	// conv to int if pos
	_, err := strconv.Atoi(day)
	if err != nil {
		log.Fatal(err)
	}

	return day, month
}
