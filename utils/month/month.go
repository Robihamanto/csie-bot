package month

var (
	January = Month{
		Short: "Jan",
		Long: "January",
	}

	March = Month{
		Short: "Mar",
		Long: "March"
	}
)


// Month represent month name
type Month struct {
	Short string
	Long  string
}
