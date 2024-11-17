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

// String returns a human-readable representation of the duration
func (d Duration) String() string {
	if d.IsZero() {
		return "0s"
	}

	// Determine if the duration is negative by checking the most significant non-zero component
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

	// Make a copy to avoid modifying the original
	temp := d
	if isNegative {
		temp = Duration{
			Years:   -d.Years,
			Months:  -d.Months,
			Days:    -d.Days,
			Hours:   -d.Hours,
			Minutes: -d.Minutes,
			Seconds: -d.Seconds,
			Nanos:   -d.Nanos,
		}
	}

	parts := []string{}

	if temp.Years > 0 {
		parts = append(parts, fmt.Sprintf("%dy", temp.Years))
	}
	if temp.Months > 0 {
		parts = append(parts, fmt.Sprintf("%dmo", temp.Months))
	}
	if temp.Days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", temp.Days))
	}
	if temp.Hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", temp.Hours))
	}
	if temp.Minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", temp.Minutes))
	}
	if temp.Seconds > 0 || (temp.Nanos > 0 && temp.Nanos >= 1000000000) {
		parts = append(parts, fmt.Sprintf("%ds", temp.Seconds))
	}

	// Handle sub-second units
	if temp.Nanos > 0 {
		nanos := temp.Nanos
		if nanos >= 1000000 {
			if nanos%1000000 == 0 {
				parts = append(parts, fmt.Sprintf("%dms", nanos/1000000))
			} else {
				parts = append(parts, fmt.Sprintf("%.3fms", float64(nanos)/1000000.0))
			}
		} else if nanos >= 1000 {
			if nanos%1000 == 0 {
				parts = append(parts, fmt.Sprintf("%dµs", nanos/1000))
			} else {
				parts = append(parts, fmt.Sprintf("%.3fµs", float64(nanos)/1000.0))
			}
		} else {
			parts = append(parts, fmt.Sprintf("%dns", nanos))
		}
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
