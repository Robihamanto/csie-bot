package month

var (
	// January Month name
	January = Month{
		Short: "Jan",
		Long:  "January",
	}

	// March Month name
	March = Month{
		Short: "Mar",
		Long:  "March",
	}
)

// Month represent month name
type Month struct {
	Short string
	Long  string
}
