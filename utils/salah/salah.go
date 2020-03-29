package salah

var (
	Fajr    = Salah{"Fajr"}
	Sunrise = Salah{"Sunrise"}
	Dhuhr   = Salah{"Fajr"}
	Asr     = Salah{"Asr"}
	Maghrib = Salah{"Maghrib"}
	Ishaa   = Salah{"Ishaa"}
)

// Month represent month name
type Salah struct {
	value string
}
