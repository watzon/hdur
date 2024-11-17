package hdur

import "time"

// Nanoseconds creates a Duration from a number of nanoseconds
func Nanoseconds(n float64) Duration {
	d := Duration{}
	total := int(n)
	d.Nanos = total
	d.normalize()
	return d
}

// Microseconds creates a Duration from a number of microseconds
func Microseconds(µs float64) Duration {
	return Nanoseconds(µs * 1000)
}

// Milliseconds creates a Duration from a number of milliseconds
func Milliseconds(ms float64) Duration {
	return Nanoseconds(ms * 1000000)
}

// Seconds creates a Duration from a number of seconds
func Seconds(s float64) Duration {
	d := Duration{}
	total := int(s)
	d.Seconds = total
	d.normalize()
	return d
}

// Minutes creates a Duration from a number of minutes
func Minutes(m float64) Duration {
	d := Duration{}
	total := int(m)
	d.Minutes = total
	d.normalize()
	return d
}

// Hours creates a Duration from a number of hours
func Hours(h float64) Duration {
	d := Duration{}
	total := int(h)
	d.Hours = total
	d.normalize()
	return d
}

// Days creates a Duration from a number of days
func Days(days float64) Duration {
	d := Duration{}
	total := int(days)
	isNegative := days < 0
	if isNegative {
		total = -total
	}

	// Convert days to months when appropriate
	if total >= 30 {
		d.Months = total / 30
		total = total % 30
	}

	d.Days = total

	if isNegative {
		d.Months = -d.Months
		d.Days = -d.Days
	}

	d.normalize()
	return d
}

// Weeks creates a Duration from a number of weeks
func Weeks(weeks float64) Duration {
	return Days(weeks * 7)
}

// Months creates a Duration from a number of months
func Months(months float64) Duration {
	d := Duration{}
	total := int(months)
	d.Months = total
	d.normalize()
	return d
}

// Years creates a Duration from a number of years
func Years(years float64) Duration {
	return Months(years * 12)
}

// FromStandard converts a time.Duration to our Duration type
func FromStandard(d time.Duration) Duration {
	return Duration{
		Seconds: int(d.Seconds()),
		Nanos:   int(d.Nanoseconds() % 1000000000),
	}
}

// InHours returns the duration as a floating-point number of hours
func (d Duration) InHours() float64 {
	t := time.Now()
	return d.Add(t).Sub(t).Hours()
}

// InMinutes returns the duration as a floating-point number of minutes
func (d Duration) InMinutes() float64 {
	t := time.Now()
	return d.Add(t).Sub(t).Minutes()
}

// InSeconds returns the duration as a floating-point number of seconds
func (d Duration) InSeconds() float64 {
	t := time.Now()
	return d.Add(t).Sub(t).Seconds()
}

// InNanoseconds returns the duration as an integer number of nanoseconds
func (d Duration) InNanoseconds() int64 {
	t := time.Now()
	return d.Add(t).Sub(t).Nanoseconds()
}

// InMonths returns the approximate number of months in the duration
func (d Duration) InMonths() float64 {
	return float64(d.Years*12+d.Months) + float64(d.Days)/30.44
}

// InYears returns the approximate number of years in the duration
func (d Duration) InYears() float64 {
	return float64(d.Years) + float64(d.Months)/12.0 + float64(d.Days)/365.25
}

// Common durations
var (
	Nanosecond  = Nanoseconds(1)
	Microsecond = Microseconds(1)
	Millisecond = Milliseconds(1)
	Second      = Seconds(1)
	Minute      = Minutes(1)
	Hour        = Hours(1)
	Day         = Days(1)
	Week        = Weeks(1)
	Month       = Months(1)
	Year        = Years(1)

	Seconds30 = Seconds(30)
	Minutes5  = Minutes(5)
	Minutes30 = Minutes(30)
	Hours2    = Hours(2)
	Hours12   = Hours(12)
	Hours24   = Hours(24)
	Days7     = Days(7)
	Days30    = Days(30)
	Days90    = Days(90)
	Days180   = Days(180)
	Days365   = Days(365)
	Months6   = Months(6)
	Years2    = Years(2)
	Years5    = Years(5)
	Years10   = Years(10)
)
