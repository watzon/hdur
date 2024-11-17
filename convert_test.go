package hdur

import (
	"math"
	"testing"
)

func TestDuration_Conversion(t *testing.T) {
	tests := []struct {
		name string
		d    Duration
	}{
		{
			name: "simple duration",
			d:    Duration{Days: 1, Hours: 2, Minutes: 30},
		},
		{
			name: "complex duration",
			d:    Duration{Years: 1, Months: 2, Days: 3, Hours: 4},
		},
		{
			name: "negative duration",
			d:    Duration{Days: -1, Hours: -12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			std := tt.d.ToStandard()
			got := FromStandard(std)

			// Compare using nanoseconds since direct comparison might fail due to month/year approximation
			if got.InNanoseconds() != tt.d.InNanoseconds() {
				t.Errorf("ToStandard/FromStandard roundtrip failed: got %v, want %v", got, tt.d)
			}
		})
	}
}

func TestDurationConversions(t *testing.T) {
	tests := []struct {
		name     string
		d        Duration
		unit     string
		expected float64
		delta    float64
	}{
		{
			name:     "hours simple",
			d:        Hours(2),
			unit:     "hours",
			expected: 2,
			delta:    0.001,
		},
		{
			name:     "minutes simple",
			d:        Minutes(120),
			unit:     "minutes",
			expected: 120,
			delta:    0.001,
		},
		{
			name:     "seconds simple",
			d:        Seconds(3600),
			unit:     "seconds",
			expected: 3600,
			delta:    0.001,
		},
		{
			name:     "months simple",
			d:        Months(2),
			unit:     "months",
			expected: 2,
			delta:    0.001,
		},
		{
			name:     "years simple",
			d:        Years(1),
			unit:     "years",
			expected: 1,
			delta:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got float64
			switch tt.unit {
			case "hours":
				got = tt.d.InHours()
			case "minutes":
				got = tt.d.InMinutes()
			case "seconds":
				got = tt.d.InSeconds()
			case "months":
				got = tt.d.InMonths()
			case "years":
				got = tt.d.InYears()
			}

			if math.Abs(got-tt.expected) > tt.delta {
				t.Errorf("got %v, want %v (±%v)", got, tt.expected, tt.delta)
			}
		})
	}
}
