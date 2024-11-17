package hdur

import "time"

// ToStandard converts our Duration to a time.Duration
// Note: This is an approximation as time.Duration doesn't handle months and years
func (d Duration) ToStandard() time.Duration {
	// Convert everything to nanoseconds
	t := time.Now()
	return d.Add(t).Sub(t)
}
