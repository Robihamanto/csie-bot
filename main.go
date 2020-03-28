package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Robihamanto/csie-bot/service"
	"github.com/Robihamanto/csie-bot/utils/config"
	"github.com/Robihamanto/csie-bot/utils/loader"
	"github.com/beevik/ntp"
)

func main() {
	cfg, err := config.Load("local")
	checkErr(err)
	log.Println(cfg)

	ntpTime, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Year:", ntpTime.Year)
	log.Println("Month:", ntpTime.Month)
	log.Println("Day:", ntpTime.Year)
	log.Println("Hour:", ntpTime.Year)
	log.Println("Minute:", ntpTime.Minute)
	log.Println("Second:", ntpTime.Second)

	ntpTimeFormatted := ntpTime.Format(time.UnixDate)

	fmt.Printf("Network time: %v\n", ntpTime)
	fmt.Printf("Unix Date Network time: %v\n", ntpTimeFormatted)

	praytimes, err := loader.Read("praytime.march.csv")

	log.Println("Praytimes :", praytimes[30])
	d, m, err := service.GetDate(&praytimes[30])

	log.Println("Praytimes :", d)
	log.Println("Praytimes :", m)

	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func localMachinetime() {
	ntpTime, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}

	ntpTimeFormatted := ntpTime.Format(time.UnixDate)

	fmt.Printf("Network time: %v\n", ntpTime)
	fmt.Printf("Unix Date Network time: %v\n", ntpTimeFormatted)
	fmt.Println("+++++++++++++++++++++++++++++++")
	timeFormatted := time.Now().Local().Format(time.UnixDate)
	fmt.Printf("System time: %v\n", time.Now())
	fmt.Printf("Unix Date System time: %v\n", timeFormatted)
}
