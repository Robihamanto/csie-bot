package loader

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/Robihamanto/csie-bot/model"
)

// Read csv files
func Read(f string) ([]model.Praytime, error) {
	csvfile, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	r := csv.NewReader(csvfile)

	var praytimes []model.Praytime

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		praytime := model.Praytime{
			Date:     record[0],
			Fajr:     record[1],
			Sunshine: record[2],
			Dhuhr:    record[3],
			Asr:      record[4],
			Maghrib:  record[5],
			Ishaa:    record[6],
		}

		praytimes = append(praytimes, praytime)
	}

	return praytimes, nil
}
