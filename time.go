package hdur

import "time"

// getDaysInMonth returns the number of days in the given year and month
func getDaysInMonth(year int, month time.Month) int {
	// This trick gets the last day of the month by going to the first day of the next month
	// and subtracting one day
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfNextMonth.Add(-time.Second)
	return lastOfMonth.Day()
}

// IsZero returns true if the duration is zero
func (d Duration) IsZero() bool {
	return d.Years == 0 && d.Months == 0 && d.Days == 0 &&
		d.Hours == 0 && d.Minutes == 0 && d.Seconds == 0 && d.Nanos == 0
}

// Until returns the duration until the given time
func Until(t time.Time) Duration {
	return Sub(t, time.Now())
}

// Since returns the duration since the given time
func Since(t time.Time) Duration {
	return Sub(time.Now(), t)
}

// Between returns the duration between two times
func Between(t1, t2 time.Time) Duration {
	return Sub(t2, t1)
}

// Add adds the duration to a time and returns the resulting time
func (d Duration) Add(t time.Time) time.Time {
	// First add years, which is straightforward
	t = t.AddDate(d.Years, 0, 0)

	// Then add months while preserving the original day of month when possible
	if d.Months > 0 {
		year, month, day := t.Date()
		hour, min, sec := t.Clock()
		nsec := t.Nanosecond()

		// Calculate target month and year
		totalMonths := int(month) + d.Months
		targetYear := year + (totalMonths-1)/12
		targetMonth := time.Month((totalMonths-1)%12 + 1)

		// Get the last day of the target month
		lastDay := time.Date(targetYear, targetMonth, 1, 0, 0, 0, 0, t.Location()).
			AddDate(0, 1, -1).Day()

		// Use the original day unless it would be invalid in the target month
		if day > lastDay {
			day = lastDay
		}

		t = time.Date(targetYear, targetMonth, day, hour, min, sec, nsec, t.Location())
	}

	// Finally add the remaining components
	return t.AddDate(0, 0, d.Days).
		Add(time.Duration(d.Hours)*time.Hour +
			time.Duration(d.Minutes)*time.Minute +
			time.Duration(d.Seconds)*time.Second +
			time.Duration(d.Nanos)*time.Nanosecond)
}

// Sub returns the duration between two times, attempting to preserve month and year units.
// The duration will be negative if t1 is before t2, following the behavior of time.Sub.
func Sub(t1, t2 time.Time) Duration {
	// Track if we need to negate the result
	needsNegation := t1.Before(t2)

	// Ensure a is the later time for calculation
	a, b := t1, t2
	if needsNegation {
		a, b = t2, t1
	}

	years := a.Year() - b.Year()
	months := int(a.Month() - b.Month())
	days := a.Day() - b.Day()
	hours := a.Hour() - b.Hour()
	minutes := a.Minute() - b.Minute()
	seconds := a.Second() - b.Second()
	nanos := a.Nanosecond() - b.Nanosecond()

	// Normalize negative values
	if nanos < 0 {
		nanos += 1000000000
		seconds--
	}
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	if days < 0 {
		// Use the previous month's length
		prevMonth := a.AddDate(0, -1, 0)
		days += getDaysInMonth(prevMonth.Year(), prevMonth.Month())
		months--
	}
	if months < 0 {
		months += 12
		years--
	}

	result := Duration{
		Years:   years,
		Months:  months,
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
		Nanos:   nanos,
	}

	if needsNegation {
		result = Duration{
			Years:   -years,
			Months:  -months,
			Days:    -days,
			Hours:   -hours,
			Minutes: -minutes,
			Seconds: -seconds,
			Nanos:   -nanos,
		}
	}

	return result
}
