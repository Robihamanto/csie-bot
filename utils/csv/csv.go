package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

var txt = []string{}

// Write csv files with string
func Write(t, p, d string) {
	file, err := os.Create("log_" + d + ".csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	txt = append(txt, fmt.Sprintf("%s %s \n", t, p))
	err = writer.Write(txt)
	checkError("Cannot write to file", err)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
