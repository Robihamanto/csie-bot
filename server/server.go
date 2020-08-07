package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Robihamanto/csie-bot/model"
	"github.com/Robihamanto/csie-bot/service/praytime"
	"github.com/Robihamanto/csie-bot/utils/csv"
	"github.com/beevik/ntp"
	"github.com/go-vgo/robotgo"
)

// Start serving time for pray reminder
func Start(s int) {
	// Find mouse position for send message field
	// x, y := robotgo.GetMousePos()
	// Do clock sync with website
	_, _, _, gap := syncClock()
	log.Println("Gap: ", gap)
	// log.Println("Hour gap: ", hGap)
	// log.Println("Minute gap: ", mGap)
	// log.Println("Second gap: ", sGap)

	// State desc
	// 1  : Fajr Reminder 	 	2 : Fajr adzan		3 : Fajr iqomah
	// 4  : Dhuhr Reminder  	5 : Dhuhr adzan		6 : Dhuhr iqomah
	// 7  : Asr Reminder  		8 : Asr adzan		9 : Asr iqomah
	// 10 : Maghrib Reminder 	11 : Maghrib adzan	12 : Maghrib iqomah
	// 13 : Isha Reminder		14 : Isha'a adzan	15 : Isha'a iqomah

	state := s
	var iqomahTime string
	var reminderTime string
	var isShouldRecync bool

	var now string
	var hs string
	var ms string
	var ss string
	mth := 10
	mtm := 35

	fajrm := fmt.Sprintf("%d:%d", mth, mtm)
	sunrm := fmt.Sprintf("%d:%d", mth, mtm)
	dhuhrm := fmt.Sprintf("%d:%d", mth, mtm+1)
	asrrm := fmt.Sprintf("%d:%d", mth, mtm+2)
	maghrm := fmt.Sprintf("%d:%d", mth, mtm+3)
	isharm := fmt.Sprintf("%d:%d", mth, mtm+4)

	pMock := model.Praytime{
		Month:    "March",
		Day:      "29",
		Date:     "31 march",
		Fajr:     fajrm,
		Sunshine: sunrm,
		Dhuhr:    dhuhrm,
		Asr:      asrrm,
		Maghrib:  maghrm,
		Ishaa:    isharm,
	}

	log.Println("Todays praytime: ", pMock)

	body := fetchMuslimProTime()
	p := createPrayTime(body)
	// p = pMock

	today := fmt.Sprintf("🕌 Today's pray time %s\nFajr: %s\nDhuhr: %s\nAsr: %s\nMaghrib: %s\nIsya: %s\n", p.Date, p.Fajr, p.Dhuhr, p.Asr, p.Maghrib, p.Ishaa)
	sendGeneralNotification(today)

	for {
		if isShouldRecync {
			isShouldRecync = false
			_, _, _, gap := syncClock()
			log.Println("Gap: ", gap)
			csv.Write("Gap:", string(gap), p.Date)
		}

		// 2. Check is time hour == 00:01 - 01:00
		// Update today's pray time
		// Reset sent pray time and iqomah status

		t := time.Now()
		h := t.Hour()
		m := t.Minute()
		s := t.Second()

		ts := (h * 3600) + (m * 60) + s

		// adjust system time with gap time, if result (-) means tertinggal (+ means lead)
		if gap < 0 {
			ts = ts + gap
		} else {
			ts = ts + gap
		}

		h = ts / 3600
		m = (ts % 3600) / 60
		s = ts % 60

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

		now = fmt.Sprintf("%s:%s", hs, ms)
		log.Printf("%s:%s", now, ss)

		logs := fmt.Sprintf("now:%s:%s state:%d", now, ss, state)

		if h == 2 {
			isShouldRecync = true
			time.Sleep(20 * time.Minute)
			body = fetchMuslimProTime()
			p = createPrayTime(body)
			state = 0
			today := fmt.Sprintf("Today's pray time %s\nFajr: %s\nDhuhr: %s\nAsr: %s\nMaghrib: %s\nIsya: %s\n", p.Date, p.Fajr, p.Dhuhr, p.Asr, p.Maghrib, p.Ishaa)
			csv.Write(logs, "Today's Pray Notifications", p.Date)
			csv.Write(logs, today, p.Date)
		}

		if state == 1 || state == 2 || state == 3 {
			reminderTime = adzanTimeReminderBuilder(p.Fajr)
			iqomahTime = iqomahTimeReminderBuilder(p.Fajr, 30)
			log.Println("Fajr: ", now, " ", reminderTime, " : ", p.Fajr, " : ", iqomahTime, " : ", now == p.Fajr)
		}

		if state == 4 || state == 5 || state == 6 {
			reminderTime = adzanTimeReminderBuilder(p.Dhuhr)
			iqomahTime = iqomahTimeReminderBuilder(p.Dhuhr, 15)
			log.Println("Dhuhr: ", now, " ", reminderTime, " : ", p.Dhuhr, " : ", iqomahTime, " : ", now == p.Dhuhr)
		}

		if state == 7 || state == 8 || state == 9 {
			reminderTime = adzanTimeReminderBuilder(p.Asr)
			iqomahTime = iqomahTimeReminderBuilder(p.Asr, 15)
			log.Println("Asr: ", now, " ", reminderTime, " : ", p.Asr, " : ", iqomahTime, " : ", now == p.Asr)
		}

		if state == 10 || state == 11 || state == 12 {
			reminderTime = adzanTimeReminderBuilder(p.Maghrib)
			iqomahTime = iqomahTimeReminderBuilder(p.Maghrib, 10)
			log.Println("Maghrib: ", now, " ", reminderTime, " : ", p.Maghrib, " : ", iqomahTime, " : ", now == p.Maghrib)
		}

		if state == 13 || state == 14 || state == 15 {
			reminderTime = adzanTimeReminderBuilder(p.Ishaa)
			iqomahTime = iqomahTimeReminderBuilder(p.Ishaa, 10)
			log.Println("Isha'a: ", now, " ", reminderTime, " : ", p.Ishaa, " : ", iqomahTime, " : ", now == p.Ishaa)
		}

		// Check fajr time from Praytime
		if now == "03:00" && state == 0 {
			state = state + 1
			today := fmt.Sprintf("Today's pray time %s\nFajr: %s\nDhuhr: %s\nAsr: %s\nMaghrib: %s\nIsya: %s\n", p.Date, p.Fajr, p.Dhuhr, p.Asr, p.Maghrib, p.Ishaa)
			sendGeneralNotification(today)
			csv.Write(logs, "Today's Pray Notifications", p.Date)
		}

		// Check fajr time from Praytime
		if now == p.Fajr && state == 2 {
			state = state + 1
			sendAdzanReminder("Fajr", p.Fajr, iqomahTime, 30)
			log.Println("Fajr Adzan")
			csv.Write(logs, "Fajr Adzan", p.Date)
		}

		// Check dhuhr time from Praytime
		if now == p.Dhuhr && state == 5 {
			state = state + 1
			sendAdzanReminder("Dhuhr", p.Dhuhr, iqomahTime, 15)
			log.Println("Dhuhr Adzan")
			csv.Write(logs, "Dhuhr Adzan", p.Date)
		}

		// Check asr time from Asr
		if now == p.Asr && state == 8 {
			state = state + 1
			sendAdzanReminder("Asr", p.Asr, iqomahTime, 15)
			log.Println("Asr Adzan")
			csv.Write(logs, "Asr Adzan", p.Date)
		}

		// Check magrib time from Praytime
		if now == p.Maghrib && state == 11 {
			state = state + 1
			sendAdzanReminder("Maghrib", p.Maghrib, iqomahTime, 10)
			log.Println("Maghrib Adzan")
			csv.Write(logs, "Maghrib Adzan", p.Date)
		}

		// Check ishaa time from Praytime
		if now == p.Ishaa && state == 14 {
			state = state + 1
			sendAdzanReminder("Isha'a", p.Ishaa, iqomahTime, 10)
			log.Println("Isha'a Adzan")
			csv.Write(logs, "Isha'a Adzan", p.Date)
		}

		if now == iqomahTime && state == 3 {
			state = state + 1
			sendIqomahReminder()
			csv.Write(logs, "Fajr Iqomah", p.Date)
		}

		if now == iqomahTime && state == 6 {
			state = state + 1
			sendIqomahReminder()
			csv.Write(logs, "Dhuhr Iqomah", p.Date)
		}

		if now == iqomahTime && state == 9 {
			state = state + 1
			sendIqomahReminder()
			csv.Write(logs, "Asr Iqomah", p.Date)
		}

		if now == iqomahTime && state == 12 {
			state = state + 1
			sendIqomahReminder()
			csv.Write(logs, "Maghrib Iqomah", p.Date)
		}

		if now == iqomahTime && state == 15 {
			state = 0
			sendIqomahReminder()
			csv.Write(logs, "Isha'a Iqomah", p.Date)
		}

		if now == reminderTime && state == 1 {
			state = state + 1
			sendAdzanTimeReminder("Fajr", p.Fajr)
			reminderTime = adzanTimeReminderBuilder(p.Dhuhr)
			csv.Write(logs, "Fajr Adzan Reminder", p.Date)
		}

		if now == reminderTime && state == 4 {
			state = state + 1
			sendAdzanTimeReminder("Dhuhr", p.Dhuhr)
			reminderTime = adzanTimeReminderBuilder(p.Asr)
			csv.Write(logs, "Dhuhr Adzan Reminder", p.Date)
		}

		if now == reminderTime && state == 7 {
			state = state + 1
			sendAdzanTimeReminder("Asr", p.Asr)
			reminderTime = adzanTimeReminderBuilder(p.Maghrib)
			csv.Write(logs, "Asr Adzan Reminder", p.Date)
		}

		if now == reminderTime && state == 10 {
			state = state + 1
			sendAdzanTimeReminder("Maghrib", p.Maghrib)
			reminderTime = adzanTimeReminderBuilder(p.Ishaa)
			csv.Write(logs, "Maghrib Adzan Reminder", p.Date)
		}

		if now == reminderTime && state == 13 {
			state = state + 1
			sendAdzanTimeReminder("Isha'a", p.Ishaa)
			reminderTime = adzanTimeReminderBuilder(p.Fajr)
			csv.Write(logs, "Isha'a Adzan Reminder", p.Date)
		}

		schedule := fmt.Sprintf("Today's pray time %s Fajr: %s Dhuhr: %s Asr: %s Maghrib: %s Isya: %s", p.Date, p.Fajr, p.Dhuhr, p.Asr, p.Maghrib, p.Ishaa)
		// log.Println(schedule)
		csv.Write(logs, schedule, p.Date)
		time.Sleep(1 * time.Second)
	}
}

func fetchMuslimProTime() string {
	resp, err := http.Get("https://www.muslimpro.com/muslimprowidget.js?cityid=6696918&timeformat=24&convention=EgyptBis")
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

func syncClock() (int, int, int, int) {

	ntpTime, err := ntp.Time("time.apple.com")
	if err != nil {
		fmt.Println(err)
	}

	ha := ntpTime.Hour()
	ma := ntpTime.Minute()
	sa := ntpTime.Second()

	t := time.Now()
	h := t.Hour()
	m := t.Minute()
	s := t.Second()

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

	return h, m, s, gap
}

func iqomahTimeBuilder(h, m, s, i int) string {
	t := (h * 3600) + (m * 60) + s
	t = t + (i * 60) // CHANGE TO I
	h = t / 3600
	m = (t % 3600) / 60
	s = t % 60

	var time string
	var hs string
	var ms string

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

	time = fmt.Sprintf("%s:%s", hs, ms)
	return time
}

func adzanTimeReminderBuilder(time string) string {
	h, _ := strconv.Atoi(time[0:2])
	m, _ := strconv.Atoi(time[3:5])

	t := (h * 3600) + (m * 60)
	t = t - (10 * 60)
	h = t / 3600
	m = (t % 3600) / 60

	var timeResult string
	var hs string
	var ms string

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

	timeResult = fmt.Sprintf("%s:%s", hs, ms)
	return timeResult
}

func iqomahTimeReminderBuilder(time string, i int) string {
	h, _ := strconv.Atoi(time[0:2])
	m, _ := strconv.Atoi(time[3:5])

	t := (h * 3600) + (m * 60)
	t = t + (i * 60)
	h = t / 3600
	m = (t % 3600) / 60

	var timeResult string
	var hs string
	var ms string

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

	timeResult = fmt.Sprintf("%s:%s", hs, ms)
	return timeResult
}

func sendAdzanReminder(p, t, q string, i int) {
	text := fmt.Sprintf("🕌 %s time for Zhongli District : %s\n Iqomah will be held in %d minutes (%s)", p, t, i, q)
	log.Println(text)
	doRobotJob(text)
}

func sendIqomahReminder() {
	text := fmt.Sprintf("Iqomah 🎉")
	log.Println(text)

	doRobotJob(text)
}

func sendAdzanTimeReminder(p, t string) {
	text := fmt.Sprintf("🕌 Reminder: 10 minutes before %s Adzan for Zhongli District : %s", p, t)
	log.Println(text)
	doRobotJob(text)
}

func sendGeneralNotification(t string) {
	log.Println(t)
	doRobotJob(t)
}

func doRobotJob(t string) {
	log.Println("Doing robot job for Engineering  5:) ", t)
	sendMessage(1594, 413, t)

	log.Println("Doing robot job for Engineering  LS:) ", t)
	sendMessage(1600, 738, t)
}

func sendMessage(x int, y int, m string) {
	csv.Write(m, "", "05 Apr")
	robotgo.MoveMouse(x, y)
	robotgo.Click("left", true)
	robotgo.PasteStr(m)
	robotgo.KeyTap("enter")
}
