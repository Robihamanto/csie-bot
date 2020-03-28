package model

import (
	"time"
)

// Praytime implement
type Praytime struct {
	Time     *time.Time
	Date     string
	Fajr     string
	Sunshine string
	Dhuhr    string
	Asr      string
	Maghrib  string
	Ishaa    string
}
