package hdur

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// Format returns a string representation of the duration using the given format
// Format specifiers:
// %y - years
// %M - months
// %d - days
// %h - hours
// %m - minutes
// %s - seconds
// %f - fractional seconds (without leading dot)
func (d Duration) Format(format string) string {
	d.normalize()

	// Handle fractional seconds specially to avoid double dots
	format = strings.ReplaceAll(format, "%s.%f", fmt.Sprintf("%d.%09d", d.Seconds, d.Nanos))

	replacer := strings.NewReplacer(
		"%y", fmt.Sprintf("%d", d.Years),
		"%M", fmt.Sprintf("%d", d.Months),
		"%d", fmt.Sprintf("%d", d.Days),
		"%h", fmt.Sprintf("%d", d.Hours),
		"%m", fmt.Sprintf("%d", d.Minutes),
		"%s", fmt.Sprintf("%d", d.Seconds),
		"%f", fmt.Sprintf("%09d", d.Nanos),
	)

	return replacer.Replace(format)
}

// abs returns the absolute value of the Duration
func (d Duration) abs() Duration {
	return Duration{
		Years:   abs(int(d.Years)),
		Months:  abs(int(d.Months)),
		Days:    abs(int(d.Days)),
		Hours:   abs(int(d.Hours)),
		Minutes: abs(int(d.Minutes)),
		Seconds: abs(int(d.Seconds)),
		Nanos:   abs(int(d.Nanos)),
	}
}

// formatMainUnits formats years through seconds
func (d Duration) formatMainUnits() []string {
	parts := []string{}
	if d.Years > 0 {
		parts = append(parts, fmt.Sprintf("%dy", d.Years))
	}
	if d.Months > 0 {
		parts = append(parts, fmt.Sprintf("%dmo", d.Months))
	}
	if d.Days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", d.Days))
	}
	if d.Hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", d.Hours))
	}
	if d.Minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", d.Minutes))
	}
	if d.Seconds > 0 || (d.Nanos > 0 && d.Nanos >= 1000000000) {
		parts = append(parts, fmt.Sprintf("%ds", d.Seconds))
	}
	return parts
}

// formatNanos formats sub-second units
func (d Duration) formatNanos() string {
	nanos := d.Nanos
	if nanos == 0 {
		return ""
	}

	switch {
	case nanos >= 1000000:
		if nanos%1000000 == 0 {
			return fmt.Sprintf("%dms", nanos/1000000)
		}
		return fmt.Sprintf("%.3fms", float64(nanos)/1000000.0)
	case nanos >= 1000:
		if nanos%1000 == 0 {
			return fmt.Sprintf("%dµs", nanos/1000)
		}
		return fmt.Sprintf("%.3fµs", float64(nanos)/1000.0)
	default:
		return fmt.Sprintf("%dns", nanos)
	}
}

// String returns a human-readable representation of the duration
func (d Duration) String() string {
	if d.IsZero() {
		return "0s"
	}

	isNegative := d.isNegativeDuration()
	temp := d
	if isNegative {
		temp = temp.abs()
	}

	parts := temp.formatMainUnits()
	if nanos := temp.formatNanos(); nanos != "" {
		parts = append(parts, nanos)
	}

	result := strings.Join(parts, " ")
	if isNegative {
		return "-" + result
	}
	return result
}

// MarshalJSON implements json.Marshaler
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsed, err := ParseDuration(s)
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}

// Value implements the driver.Valuer interface
func (d Duration) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements the sql.Scanner interface
func (d *Duration) Scan(value interface{}) error {
	var err error

	switch v := value.(type) {
	case []byte:
		*d, err = ParseDuration(string(v))
	case string:
		*d, err = ParseDuration(v)
	case nil:
		*d = Duration{}
	default:
		err = fmt.Errorf("cannot scan type %T into Duration", value)
	}

	return err
}
