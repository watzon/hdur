package hdur

import "testing"

func TestDurationConstructors(t *testing.T) {
	tests := []struct {
		name     string
		input    Duration
		expected string
	}{
		{
			name:     "nanoseconds below second",
			input:    Nanoseconds(500000000),
			expected: "500ms",
		},
		{
			name:     "nanoseconds above second",
			input:    Nanoseconds(1500000000),
			expected: "1s 500ms",
		},
		{
			name:     "microseconds",
			input:    Microseconds(1500),
			expected: "1.500ms",
		},
		{
			name:     "milliseconds",
			input:    Milliseconds(1500),
			expected: "1s 500ms",
		},
		{
			name:     "seconds below minute",
			input:    Seconds(45),
			expected: "45s",
		},
		{
			name:     "seconds above minute",
			input:    Seconds(90),
			expected: "1m 30s",
		},
		{
			name:     "minutes below hour",
			input:    Minutes(45),
			expected: "45m",
		},
		{
			name:     "minutes above hour",
			input:    Minutes(90),
			expected: "1h 30m",
		},
		{
			name:     "hours below day",
			input:    Hours(23),
			expected: "23h",
		},
		{
			name:     "hours above day",
			input:    Hours(25),
			expected: "1d 1h",
		},
		{
			name:     "days below month",
			input:    Days(25),
			expected: "25d",
		},
		{
			name:     "days above month",
			input:    Days(45),
			expected: "1mo 15d",
		},
		{
			name:     "weeks",
			input:    Weeks(2),
			expected: "14d",
		},
		{
			name:     "months below year",
			input:    Months(11),
			expected: "11mo",
		},
		{
			name:     "months above year",
			input:    Months(14),
			expected: "1y 2mo",
		},
		{
			name:     "years",
			input:    Years(2),
			expected: "2y",
		},
		{
			name:     "fractional values",
			input:    Hours(1.5),
			expected: "1h", // We currently truncate fractional values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.String()
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDuration_Constructors_Comprehensive(t *testing.T) {
	tests := []struct {
		name  string
		fn    func(float64) Duration
		input float64
		want  Duration
	}{
		// Nanoseconds
		{
			name:  "nanoseconds zero",
			fn:    Nanoseconds,
			input: 0,
			want:  Duration{},
		},
		{
			name:  "nanoseconds positive",
			fn:    Nanoseconds,
			input: 1500000000,
			want:  Duration{Seconds: 1, Nanos: 500000000},
		},
		{
			name:  "nanoseconds negative",
			fn:    Nanoseconds,
			input: -1500000000,
			want:  Duration{Seconds: -1, Nanos: -500000000},
		},
		{
			name:  "nanoseconds fractional",
			fn:    Nanoseconds,
			input: 1.5,
			want:  Duration{Nanos: 1}, // Should truncate
		},

		// Seconds
		{
			name:  "seconds zero",
			fn:    Seconds,
			input: 0,
			want:  Duration{},
		},
		{
			name:  "seconds positive",
			fn:    Seconds,
			input: 65,
			want:  Duration{Minutes: 1, Seconds: 5},
		},
		{
			name:  "seconds negative",
			fn:    Seconds,
			input: -65,
			want:  Duration{Minutes: -1, Seconds: -5},
		},
		{
			name:  "seconds fractional",
			fn:    Seconds,
			input: 1.5,
			want:  Duration{Seconds: 1}, // Should truncate
		},

		// Minutes
		{
			name:  "minutes zero",
			fn:    Minutes,
			input: 0,
			want:  Duration{},
		},
		{
			name:  "minutes positive",
			fn:    Minutes,
			input: 65,
			want:  Duration{Hours: 1, Minutes: 5},
		},
		{
			name:  "minutes negative",
			fn:    Minutes,
			input: -65,
			want:  Duration{Hours: -1, Minutes: -5},
		},
		{
			name:  "minutes fractional",
			fn:    Minutes,
			input: 1.5,
			want:  Duration{Minutes: 1}, // Should truncate
		},

		// Hours
		{
			name:  "hours zero",
			fn:    Hours,
			input: 0,
			want:  Duration{},
		},
		{
			name:  "hours positive",
			fn:    Hours,
			input: 25,
			want:  Duration{Days: 1, Hours: 1},
		},
		{
			name:  "hours negative",
			fn:    Hours,
			input: -25,
			want:  Duration{Days: -1, Hours: -1},
		},
		{
			name:  "hours fractional",
			fn:    Hours,
			input: 1.5,
			want:  Duration{Hours: 1}, // Should truncate
		},

		// Months
		{
			name:  "months zero",
			fn:    Months,
			input: 0,
			want:  Duration{},
		},
		{
			name:  "months positive",
			fn:    Months,
			input: 14,
			want:  Duration{Years: 1, Months: 2},
		},
		{
			name:  "months negative",
			fn:    Months,
			input: -14,
			want:  Duration{Years: -1, Months: -2},
		},
		{
			name:  "months fractional",
			fn:    Months,
			input: 1.5,
			want:  Duration{Months: 1}, // Should truncate
		},
		{
			name:  "months large value",
			fn:    Months,
			input: 25,
			want:  Duration{Years: 2, Months: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input)
			if !got.Equal(tt.want) {
				t.Errorf("%T(%v) = %v, want %v", tt.fn, tt.input, got, tt.want)
			}
		})
	}
}

func TestDays_Edge_Cases(t *testing.T) {
	tests := []struct {
		name string
		days float64
		want Duration
	}{
		{
			name: "exactly 30 days",
			days: 30,
			want: Duration{Months: 1},
		},
		{
			name: "31 days",
			days: 31,
			want: Duration{Months: 1, Days: 1},
		},
		{
			name: "negative 30 days",
			days: -30,
			want: Duration{Months: -1},
		},
		{
			name: "negative 31 days",
			days: -31,
			want: Duration{Months: -1, Days: -1},
		},
		{
			name: "60 days",
			days: 60,
			want: Duration{Months: 2},
		},
		{
			name: "365 days",
			days: 365,
			want: Duration{Months: 12, Days: 5},
		},
		{
			name: "zero days",
			days: 0,
			want: Duration{},
		},
		{
			name: "fractional days",
			days: 30.5,
			want: Duration{Months: 1, Days: 0}, // Should truncate
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Days(tt.days)
			if !got.Equal(tt.want) {
				t.Errorf("Days() = %v, want %v", got, tt.want)
			}
		})
	}
}
