package hdur

import (
	"math"
	"time"
)

// Equal returns true if the durations are equal
func (d Duration) Equal(other Duration) bool {
	// Normalize both durations before comparison
	d.normalize()
	other.normalize()

	return d.Years == other.Years &&
		d.Months == other.Months &&
		d.Days == other.Days &&
		d.Hours == other.Hours &&
		d.Minutes == other.Minutes &&
		d.Seconds == other.Seconds &&
		d.Nanos == other.Nanos
}

// Less returns true if d is less than other
func (d Duration) Less(other Duration) bool {
	// Compare using a fixed point in time
	t := time.Now()
	return d.Add(t).Before(other.Add(t))
}

// LessOrEqual returns true if d is less than or equal to other
func (d Duration) LessOrEqual(other Duration) bool {
	return d.Less(other) || d.Equal(other)
}

// Greater returns true if d is greater than other
func (d Duration) Greater(other Duration) bool {
	return !d.LessOrEqual(other)
}

// GreaterOrEqual returns true if d is greater than or equal to other
func (d Duration) GreaterOrEqual(other Duration) bool {
	return !d.Less(other)
}

// Abs returns the absolute value of the duration
func (d Duration) Abs() Duration {
	if d.Less(Duration{}) {
		return d.Mul(-1)
	}
	return d
}

// Mul returns the duration multiplied by the given factor
func (d Duration) Mul(factor float64) Duration {
	if factor == 0 {
		return Duration{}
	}

	// For negative factors, multiply by the absolute value and then negate
	if factor < 0 {
		factor = -factor
		result := Duration{
			Years:   -int(float64(d.Years) * factor),
			Months:  -int(float64(d.Months) * factor),
			Days:    -int(float64(d.Days) * factor),
			Hours:   -int(float64(d.Hours) * factor),
			Minutes: -int(float64(d.Minutes) * factor),
			Seconds: -int(float64(d.Seconds) * factor),
			Nanos:   -int(float64(d.Nanos) * factor),
		}
		result.normalize()
		return result
	}

	result := Duration{
		Years:   int(float64(d.Years) * factor),
		Months:  int(float64(d.Months) * factor),
		Days:    int(float64(d.Days) * factor),
		Hours:   int(float64(d.Hours) * factor),
		Minutes: int(float64(d.Minutes) * factor),
		Seconds: int(float64(d.Seconds) * factor),
		Nanos:   int(float64(d.Nanos) * factor),
	}

	result.normalize()
	return result
}

// Div returns the duration divided by the given divisor
func (d Duration) Div(divisor float64) Duration {
	if divisor == 0 {
		panic("division by zero")
	}
	return d.Mul(1 / divisor)
}

// Round rounds the duration to the nearest multiple of the given duration
func (d Duration) Round(multiple Duration) Duration {
	if multiple.IsZero() {
		return d
	}

	// Convert to nanoseconds for comparison
	t := time.Now()
	dNanos := d.Add(t).Sub(t)
	mNanos := multiple.Add(t).Sub(t)

	if mNanos == 0 {
		return d
	}

	rounded := time.Duration(math.Round(float64(dNanos)/float64(mNanos))) * time.Duration(mNanos)
	result := Duration{
		Seconds: int(rounded.Seconds()),
		Nanos:   int(rounded.Nanoseconds() % 1000000000),
	}
	result.normalize()
	return result
}

// Truncate truncates the duration to a multiple of the given duration
func (d Duration) Truncate(multiple Duration) Duration {
	if multiple.IsZero() {
		return d
	}

	// Convert to nanoseconds for comparison
	t := time.Now()
	dNanos := d.Add(t).Sub(t)
	mNanos := multiple.Add(t).Sub(t)

	if mNanos == 0 {
		return d
	}

	truncated := time.Duration(dNanos/mNanos) * time.Duration(mNanos)
	result := Duration{
		Seconds: int(truncated.Seconds()),
		Nanos:   int(truncated.Nanoseconds() % 1000000000),
	}
	result.normalize()
	return result
}

func abs(x int) int {
	return int(math.Abs(float64(x)))
}
