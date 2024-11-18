package hdur

import (
	"reflect"
	"time"
)

// ToStandard converts our Duration to a time.Duration
// Note: This is an approximation as time.Duration doesn't handle months and years
func (d Duration) ToStandard() time.Duration {
	// Convert everything to nanoseconds
	t := time.Now()
	return d.Add(t).Sub(t)
}

// ToValue converts a duration string to a reflect.Value, perfect for use
// with Fiber
func ToValue(value string) reflect.Value {
	if v, err := ParseDuration(value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{}
}
