/*
package hdur provides a flexible and human-friendly duration type that extends Go's time.Duration.

Key features:
  - Calendar-aware duration handling (months, years)
  - Human-readable parsing ("1 year 2 months 3 days")
  - Rich set of duration manipulation methods
  - SQL and JSON serialization support
  - Comprehensive time unit support from nanoseconds to years

Example usage:

	// Parse a duration string
	d, err := ParseDuration("1 year 2 months 3 days")
	if err != nil {
		log.Fatal(err)
	}

	// Add duration to a time
	future := d.Add(time.Now())

	// Create durations using constructors
	d1 := Hours(24)
	d2 := Days(1)
	if d1.Equal(d2) {
		fmt.Println("24 hours equals 1 day")
	}

	// Use duration constants
	timeout := Seconds30
	time.Sleep(timeout.ToStandard())

	// Format durations
	d3 := Months(3).Add(Days(15))
	fmt.Println(d3.Format("%M months and %d days"))

	// Compare durations
	if d1.Less(d2) {
		fmt.Println("d1 is shorter than d2")
	}

	// Mathematical operations
	doubled := d1.Mul(2)
	halved := d1.Div(2)
*/

package hdur

// Duration represents a time duration with extended functionality
type Duration struct {
	Years   int
	Months  int
	Days    int
	Hours   int
	Minutes int
	Seconds int
	Nanos   int
}

// normalize ensures all duration components are within their natural ranges
func (d *Duration) normalize() {
	// First, determine if the duration is negative by checking the most significant non-zero component
	isNegative := false
	if d.Years < 0 {
		isNegative = true
	} else if d.Years == 0 && d.Months < 0 {
		isNegative = true
	} else if d.Years == 0 && d.Months == 0 && d.Days < 0 {
		isNegative = true
	} else if d.Years == 0 && d.Months == 0 && d.Days == 0 && d.Hours < 0 {
		isNegative = true
	} else if d.Years == 0 && d.Months == 0 && d.Days == 0 && d.Hours == 0 && d.Minutes < 0 {
		isNegative = true
	} else if d.Years == 0 && d.Months == 0 && d.Days == 0 && d.Hours == 0 && d.Minutes == 0 && d.Seconds < 0 {
		isNegative = true
	} else if d.Years == 0 && d.Months == 0 && d.Days == 0 && d.Hours == 0 && d.Minutes == 0 && d.Seconds == 0 && d.Nanos < 0 {
		isNegative = true
	}

	// Make all components positive for normalization
	if isNegative {
		d.Years = -d.Years
		d.Months = -d.Months
		d.Days = -d.Days
		d.Hours = -d.Hours
		d.Minutes = -d.Minutes
		d.Seconds = -d.Seconds
		d.Nanos = -d.Nanos
	}

	// Handle nanoseconds overflow
	if d.Nanos >= 1000000000 {
		d.Seconds += d.Nanos / 1000000000
		d.Nanos = d.Nanos % 1000000000
	}

	// Handle seconds overflow
	if d.Seconds >= 60 {
		d.Minutes += d.Seconds / 60
		d.Seconds = d.Seconds % 60
	}

	// Handle minutes overflow
	if d.Minutes >= 60 {
		d.Hours += d.Minutes / 60
		d.Minutes = d.Minutes % 60
	}

	// Handle hours overflow
	if d.Hours >= 24 {
		d.Days += d.Hours / 24
		d.Hours = d.Hours % 24
	}

	// We don't normalize days to months by default, as months have variable lengths
	// Only normalize months to years
	if d.Months >= 12 {
		d.Years += d.Months / 12
		d.Months = d.Months % 12
	}

	// Restore negative sign if needed
	if isNegative {
		d.Years = -d.Years
		d.Months = -d.Months
		d.Days = -d.Days
		d.Hours = -d.Hours
		d.Minutes = -d.Minutes
		d.Seconds = -d.Seconds
		d.Nanos = -d.Nanos
	}
}
