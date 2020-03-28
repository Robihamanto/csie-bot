package service

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/Robihamanto/csie-bot/model"
)

// GetDate change string to date
// return date, month, err
func GetDate(t *model.Praytime) (int, string, error) {
	lt := len(t.Date)
	month := t.Date[lt-3 : lt]
	date := t.Date[lt-6 : lt-4]
	date = getNumberOnly(date)

	log.Println("Month:", month)
	log.Println("Date:", date)

	dateint, err := strconv.Atoi(date)
	if err == nil {
		fmt.Println(dateint)
	}

	return dateint, month, nil
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
