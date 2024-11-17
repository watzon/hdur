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

// isNegativeDuration checks if the duration is negative by examining
// the first non-zero component in order of significance
func (d *Duration) isNegativeDuration() bool {
	switch {
	case d.Years != 0:
		return d.Years < 0
	case d.Months != 0:
		return d.Months < 0
	case d.Days != 0:
		return d.Days < 0
	case d.Hours != 0:
		return d.Hours < 0
	case d.Minutes != 0:
		return d.Minutes < 0
	case d.Seconds != 0:
		return d.Seconds < 0
	case d.Nanos != 0:
		return d.Nanos < 0
	default:
		return false
	}
}

// makePositive makes all components of the duration positive
func (d *Duration) makePositive() {
	d.Years = abs(d.Years)
	d.Months = abs(d.Months)
	d.Days = abs(d.Days)
	d.Hours = abs(d.Hours)
	d.Minutes = abs(d.Minutes)
	d.Seconds = abs(d.Seconds)
	d.Nanos = abs(d.Nanos)
}

// makeNegative makes all components of the duration negative
func (d *Duration) makeNegative() {
	d.Years = -abs(d.Years)
	d.Months = -abs(d.Months)
	d.Days = -abs(d.Days)
	d.Hours = -abs(d.Hours)
	d.Minutes = -abs(d.Minutes)
	d.Seconds = -abs(d.Seconds)
	d.Nanos = -abs(d.Nanos)
}

// normalizeTimeUnits normalizes time units from smallest to largest
func (d *Duration) normalizeTimeUnits() {
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
}

// normalizeMonthsToYears normalizes months to years
func (d *Duration) normalizeMonthsToYears() {
	if d.Months >= 12 {
		d.Years += d.Months / 12
		d.Months = d.Months % 12
	}
}

// normalize ensures all duration components are within their natural ranges
func (d *Duration) normalize() {
	// First, determine if the duration is negative
	isNegative := d.isNegativeDuration()

	// Make all components positive for normalization
	if isNegative {
		d.makePositive()
	}

	// Normalize all time units
	d.normalizeTimeUnits()

	// Normalize months to years
	d.normalizeMonthsToYears()

	// Restore negative sign if needed
	if isNegative {
		d.makeNegative()
	}
}
