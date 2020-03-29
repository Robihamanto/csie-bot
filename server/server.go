package server

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

// Start serving time for pray reminder
func Start() {
	// Do clock sync with website
	hGap, mGap, sGap := syncClock()
	// log.Println("Hour gap: ", hGap)
	// log.Println("Minute gap: ", mGap)
	// log.Println("Second gap: ", sGap)

	// State desc
	// 1 : Fajr adzan
	// 2 : Dhuhr adzan
	// 3 : Asr adzan
	// 4 : Maghrib adzan
	// 5 : Isha'a adzan

	var iqomah bool
	var iqomahTime string

	var now string
	var hs string
	var ms string
	var ss string
	mth := 15
	mtm := 46
	mts := 00

	fajrm := fmt.Sprintf("%d:%d:%d", mth, mtm, mts)
	sunrm := fmt.Sprintf("%d:%d:%d", mth, mtm, mts+30)
	dhuhrm := fmt.Sprintf("%d:%d:%d", mth, mtm, mts)
	asrrm := fmt.Sprintf("%d:%d:%d", mth, mtm, mts)
	maghrm := fmt.Sprintf("%d:%d:%d", mth, mtm, mts)
	isharm := fmt.Sprintf("%d:%d:%d", mth, mtm, mts)

	pMock := model.Praytime{
		Month:    "March",
		Day:      "29",
		Date:     "29 march",
		Fajr:     fajrm,
		Sunshine: sunrm,
		Dhuhr:    dhuhrm,
		Asr:      asrrm,
		Maghrib:  maghrm,
		Ishaa:    isharm,
	}

	resetTime := "00:00:01"
	body := fetchMuslimProTime()
	p := createPrayTime(body)
	p = pMock

	for {
		time.Sleep(1 * time.Second)
		// 2. Check is time hour == 00:01 - 01:00
		// Update today's pray time
		// Reset sent pray time and iqomah status

		t := time.Now()
		h := t.Hour()
		m := t.Minute()
		s := t.Second()

		// adjust system time with gap time (-) means tertinggal (+ means lead)
		if hGap < 0 || mGap < 0 || sGap < 0 {
			h = h - hGap
			m = m - mGap
			s = s - sGap
		} else {
			h = h + hGap
			m = m + mGap
			s = s + sGap
		}

		if h < 10 {
			hs = fmt.Sprintf("0%d", h)
		} else {
			hs = fmt.Sprintf("%d", h)
		}

		if m < 10 {
			ms = fmt.Sprintf("0%d", m)
		} else {
			ms = fmt.Sprintf("%d", m)
		}

		if s < 10 {
			ss = fmt.Sprintf("0%d", s)
		} else {
			ss = fmt.Sprintf("%d", s)
		}

		now = fmt.Sprintf("%s:%s:%s", hs, ms, ss)
		log.Println(now)

		if now == resetTime {
			body = fetchMuslimProTime()
			p = createPrayTime(body)
			log.Println("Praytime ", p)

			log.Println("Now:\t", now)
			log.Println("Fajr:\t", p.Fajr)
			log.Println("Dhuhr:\t", p.Dhuhr)
			log.Println("Asr:\t", p.Asr)
			log.Println("Maghrib:\t", p.Maghrib)
			log.Println("Ishaa:\t", p.Ishaa)
		}

		// Check fajr time from Praytime
		if now == p.Fajr && !iqomah {
			// Send adzan reminder
			log.Println("Fajr Adzan")
			iqomah = true
			iqomahTime = iqomahTimeBuilder(h, m, s, 30)
			log.Println("Iqomah :", iqomahTime)
		}

		// Check dhuhr time from Praytime
		if now == p.Dhuhr && !iqomah {
			// Send adzan reminder
			log.Println("Dhuhr Adzan")
			iqomah = true
			iqomahTime = iqomahTimeBuilder(h, m, s, 20)
		}

		// Check asr time from Asr
		if now == p.Asr && !iqomah {
			// Send adzan reminder
			log.Println("Asr Adzan")
			iqomah = true
			iqomahTime = iqomahTimeBuilder(h, m, s, 20)
		}

		// Check magrib time from Praytime
		if now == p.Maghrib && !iqomah {
			// Send adzan reminder
			log.Println("Maghrib Adzan")
			iqomah = true
			iqomahTime = iqomahTimeBuilder(h, m, s, 15)
		}

		// Check ishaa time from Praytime
		if now == p.Ishaa && !iqomah {
			// Send adzan reminder
			log.Println("Ishaa Adzan")
			iqomah = true
			iqomahTime = iqomahTimeBuilder(h, m, s, 15)
		}

		if iqomah && iqomahTime == now {
			iqomah = false
			// Send iqomah reminder
			log.Println("Iqomah")
		}
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
	v = fmt.Sprintf("%s:00", v)
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

	// year, month, day := ntpTime.Date()
	// log.Printf("Internet Today %d %s %d", year, month, day)

	ha := ntpTime.Hour()
	ma := ntpTime.Minute()
	sa := ntpTime.Second()
	// log.Printf("Clock %d:%d:%d", ha, ma, sa)

	t := time.Now()
	// y := t.Year()
	// mon := t.Month()
	// d := t.Day()
	// log.Printf("System Today %d %s %d", y, mon, d)

	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	// log.Printf("Clock %d:%d:%d", h, m, s)

	// Case 1
	// 14:05:30 world
	// 14:05:15 local

	// Case 2
	// 14:06:05 world
	// 14:05:50 local

	gap := ((ha * 3600) - (h * 3600)) + ((ma * 60) - (m * 60)) + (sa - s)

	h = gap / 3600
	m = (gap % 3600) / 60
	s = gap % 60

	return h, m, s
}

func iqomahTimeBuilder(h, m, s, i int) string {
	t := (h * 3600) + (m * 60) + s
	t = t + 10 // CHANGE TO I
	h = t / 3600
	m = (t % 3600) / 60
	s = t % 60

	var time string
	var hs string
	var ms string
	var ss string
	if h < 10 {
		hs = fmt.Sprintf("0%d", h)
	} else {
		hs = fmt.Sprintf("%d", h)
	}

	if m < 10 {
		ms = fmt.Sprintf("0%d", m)
	} else {
		ms = fmt.Sprintf("%d", m)
	}

	if s < 10 {
		ss = fmt.Sprintf("0%d", s)
	} else {
		ss = fmt.Sprintf("%d", s)
	}

	time = fmt.Sprintf("%s:%s:%s", hs, ms, ss)
	return time
}
