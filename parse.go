package hdur

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var durationRegex = regexp.MustCompile(`(\d+)\s*([a-zA-Z]+)`)

var unitMap = map[string]string{
	"ns":           "nanos",
	"nano":         "nanos",
	"nanos":        "nanos",
	"nanosecond":   "nanos",
	"nanoseconds":  "nanos",
	"us":           "micros",
	"micro":        "micros",
	"micros":       "micros",
	"microsecond":  "micros",
	"microseconds": "micros",
	"ms":           "millis",
	"milli":        "millis",
	"millis":       "millis",
	"millisecond":  "millis",
	"milliseconds": "millis",
	"s":            "seconds",
	"sec":          "seconds",
	"secs":         "seconds",
	"second":       "seconds",
	"seconds":      "seconds",
	"m":            "minutes",
	"min":          "minutes",
	"mins":         "minutes",
	"minute":       "minutes",
	"minutes":      "minutes",
	"h":            "hours",
	"hr":           "hours",
	"hrs":          "hours",
	"hour":         "hours",
	"hours":        "hours",
	"d":            "days",
	"day":          "days",
	"days":         "days",
	"w":            "weeks",
	"wk":           "weeks",
	"week":         "weeks",
	"weeks":        "weeks",
	"fortnight":    "fortnights",
	"fortnights":   "fortnights",
	"y":            "years",
	"yr":           "years",
	"year":         "years",
	"years":        "years",
	"mo":           "months",
	"mon":          "months",
	"month":        "months",
	"months":       "months",
}

// normalizeInput prepares the input string for parsing
func normalizeInput(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " and ", " ")
	return strings.TrimSpace(s)
}

// parseNumber extracts and validates the numeric part of a duration component
func parseNumber(numStr string) (int64, error) {
	n, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", numStr)
	}
	return n, nil
}

// normalizeUnit validates and normalizes the time unit
func normalizeUnit(unit string) (string, error) {
	unit = strings.TrimSpace(unit)
	normalized, ok := unitMap[unit]
	if !ok {
		return "", fmt.Errorf("unknown unit: %s", unit)
	}
	return normalized, nil
}

// applyUnit adds the specified duration to the Duration struct
func applyUnit(d *Duration, n int64, unit string) {
	switch unit {
	case "nanos":
		d.Nanos += int(n)
	case "micros":
		d.Nanos += int(n * 1000)
	case "millis":
		d.Nanos += int(n * 1000000)
	case "seconds":
		d.Seconds += int(n)
	case "minutes":
		d.Minutes += int(n)
	case "hours":
		d.Hours += int(n)
	case "days":
		d.Days += int(n)
	case "weeks":
		d.Days += int(n * 7)
	case "fortnights":
		d.Days += int(n * 14)
	case "months":
		d.Months += int(n)
	case "years":
		d.Years += int(n)
	}
}

// ParseDuration parses a duration string and returns a Duration
// It supports multiple time units and ignores conjunctions like "and"
// Example: "1 day 3 hours and 5 minutes" or "2weeks 4days"
func ParseDuration(s string) (Duration, error) {
	s = normalizeInput(s)

	matches := durationRegex.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return Duration{}, fmt.Errorf("invalid duration format: %s", s)
	}

	d := Duration{}

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		n, err := parseNumber(match[1])
		if err != nil {
			return Duration{}, err
		}

		unit, err := normalizeUnit(match[2])
		if err != nil {
			return Duration{}, err
		}

		applyUnit(&d, n, unit)
	}

	d.normalize()
	return d, nil
}

// MustParseDuration is like ParseDuration but panics if the string cannot be parsed
func MustParseDuration(s string) Duration {
	d, err := ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

// ParseDurationFromNanos creates a Duration from nanoseconds
func ParseDurationFromNanos(nanos int64) Duration {
	d := Duration{Nanos: int(nanos)}
	d.normalize()
	return d
}
