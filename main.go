package main

import (
	"fmt"
	"time"

	"github.com/beevik/ntp"
)

func main() {
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
