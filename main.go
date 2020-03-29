package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Robihamanto/csie-bot/model"
	"github.com/Robihamanto/csie-bot/service/praytime"
	"github.com/beevik/ntp"
)

func main() {
	// Do clock sync with website
	hGap, mGap, sGap := syncClock()
	log.Println("Hour gap: ", hGap)
	log.Println("Minute gap: ", mGap)
	log.Println("Second gap: ", sGap)

	// State desc
	// 1 : Fajr adzan, 1 : Fajr iqomah
	// 2 : Dhuhr adzan, 3 : Dhuhr iqomah
	// 3 : Asr adzan, 5 : Asr iqomah
	// 4 : Maghrib adzan, 7 : Maghrib iqomah
	// 5 : Isha'a adzan, 9 : Isha'a iqomah

	var iqomah bool
	var state int

	// 2. Check is time hour == 00:01 - 01:00
	// Update today's pray time
	// Reset sent pray time and iqomah status
	body := fetchMuslimProTime()
	p := createPrayTime(body)
	state = 0
	log.Println("Praytime := ", p)

	var now string
	// if h < 10 {
	// 	now = fmt.Sprintf("0%d:%d", h, m)
	// } else {
	// 	now = fmt.Sprintf("%d:%d", h, m)
	// }

	log.Println("Now: ", now)
	log.Println("Fajr: ", p.Fajr)

	// Check fajr time from Praytime
	if now == p.Fajr && state == 1 {
		// Send adzan reminder
	}

	// Check dhuhr time from Praytime
	if now == p.Dhuhr && state == 2 {
		// Send adzan reminder
	}

	// Check asr time from Asr
	if now == p.Asr && state == 3 {
		// Send adzan reminder
	}

	// Check magrib time from Praytime
	if now == p.Maghrib && state == 4 {
		// Send adzan reminder
	}

	// Check ishaa time from Praytime
	if now == p.Ishaa && state == 5 {
		// Send adzan reminder
	}

	if iqomah {

	}

}

func fetchMuslimProTime() string {
	resp, err := http.Get("http://www.muslimpro.com/muslimprowidget.js?cityid=6696918&timeformat=24&convention=EgyptBis")
	checkErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	return string(body)
}

func createPrayTime(source string) model.Praytime {

	date := findString(source, "<div class=\"daterow\">", "</div></td>")
	fajr := findPraytimeString(source, "<td>Fajr</td><td>")
	sunshine := findPraytimeString(source, "<td>Sunrise</td><td>")
	dhuhr := findPraytimeString(source, "<td>Dhuhr</td><td>")
	asr := findPraytimeString(source, "<td>Asr</td><td>")
	maghrib := findPraytimeString(source, "<td>Maghrib</td><td>")
	ishaa := findPraytimeString(source, "<td>Isha&#39;a</td><td>")
	day, month := praytime.GetDate(date)

	p := model.Praytime{
		Month:    month,
		Day:      day,
		Date:     date,
		Fajr:     fajr,
		Sunshine: sunshine,
		Dhuhr:    dhuhr,
		Asr:      asr,
		Maghrib:  maghrib,
		Ishaa:    ishaa,
	}
	return p
}

func findString(source, fs string, fe string) string {
	start := strings.Index(source, fs) + len(fs)
	end := strings.Index(source, fe)
	v := source[start:end]
	return v
}

func findPraytimeString(source, fs string) string {
	start := strings.Index(source, fs) + len(fs)
	end := start + 5
	v := source[start:end]
	return v
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func syncClock() (int, int, int) {

	ntpTime, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}

	year, month, day := ntpTime.Date()
	log.Printf("Internet Today %d %s %d", year, month, day)

	ha := ntpTime.Hour()
	ma := ntpTime.Minute()
	sa := ntpTime.Second()
	log.Printf("Clock %d:%d:%d", ha, ma, sa)

	t := time.Now()
	y := t.Year()
	mon := t.Month()
	d := t.Day()
	log.Printf("System Today %d %s %d", y, mon, d)

	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	log.Printf("Clock %d:%d:%d", h, m, s)

	return (ha - h), (ma - m), (sa - s)
}
