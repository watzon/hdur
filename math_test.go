package hdur

import (
	"strings"
	"testing"
	"time"
)

func BenchmarkDurationMath(b *testing.B) {
	d := ParseMust("1y 2mo 3d")
	ops := []struct {
		name   string
		factor float64
	}{
		{"mul_2", 2},
		{"mul_0.5", 0.5},
		{"div_2", 2},
		{"div_0.5", 0.5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := ops[i%len(ops)]
		if strings.HasPrefix(op.name, "mul") {
			_ = d.Mul(op.factor)
		} else {
			_ = d.Div(op.factor)
		}
	}
}

func TestDurationComparison(t *testing.T) {
	tests := []struct {
		name     string
		d1       Duration
		d2       Duration
		expected bool
	}{
		{
			name:     "equal durations",
			d1:       Hours(24),
			d2:       Days(1),
			expected: true,
		},
		{
			name:     "different units same duration",
			d1:       Minutes(60),
			d2:       Hours(1),
			expected: true,
		},
		{
			name:     "different durations",
			d1:       Hours(1),
			d2:       Hours(2),
			expected: false,
		},
		{
			name:     "complex durations",
			d1:       Days(30),
			d2:       Months(1),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert both to nanoseconds for comparison
			t1 := time.Now()

			end1 := tt.d1.Add(t1)
			end2 := tt.d2.Add(t1)

			got := end1.Equal(end2)
			if got != tt.expected {
				t.Errorf("Duration comparison: got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDurationMathOperations(t *testing.T) {
	tests := []struct {
		name     string
		d        Duration
		op       string
		factor   float64
		expected string
	}{
		{
			name:     "multiply simple duration",
			d:        Hours(2),
			op:       "mul",
			factor:   2,
			expected: "4h",
		},
		{
			name:     "multiply complex duration",
			d:        ParseMust("1y 2mo 3d"),
			op:       "mul",
			factor:   2,
			expected: "2y 4mo 6d",
		},
		{
			name:     "divide simple duration",
			d:        Hours(4),
			op:       "div",
			factor:   2,
			expected: "2h",
		},
		{
			name:     "divide complex duration",
			d:        ParseMust("2y 4mo 6d"),
			op:       "div",
			factor:   2,
			expected: "1y 2mo 3d",
		},
		{
			name:     "multiply by zero",
			d:        Hours(2),
			op:       "mul",
			factor:   0,
			expected: "0s",
		},
		{
			name:     "multiply by negative",
			d:        Hours(2),
			op:       "mul",
			factor:   -1,
			expected: "-2h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Duration
			switch tt.op {
			case "mul":
				result = tt.d.Mul(tt.factor)
			case "div":
				result = tt.d.Div(tt.factor)
			}

			if got := result.String(); got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}

	// Test division by zero panic
	t.Run("division by zero", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic on division by zero")
			}
		}()
		Hours(1).Div(0)
	})
}

func TestDurationRoundingOperations(t *testing.T) {
	tests := []struct {
		name     string
		d        Duration
		multiple Duration
		op       string
		expected string
	}{
		{
			name:     "round to hour",
			d:        Minutes(90),
			multiple: Hour,
			op:       "round",
			expected: "2h",
		},
		{
			name:     "round to day",
			d:        Hours(36),
			multiple: Day,
			op:       "round",
			expected: "2d",
		},
		{
			name:     "truncate to hour",
			d:        Minutes(90),
			multiple: Hour,
			op:       "truncate",
			expected: "1h",
		},
		{
			name:     "truncate to day",
			d:        Hours(36),
			multiple: Day,
			op:       "truncate",
			expected: "1d",
		},
		{
			name:     "round to zero multiple",
			d:        Hours(2),
			multiple: Duration{},
			op:       "round",
			expected: "2h",
		},
		{
			name:     "round exactly half (to even)",
			d:        Duration{Minutes: 7, Seconds: 30},
			multiple: Duration{Minutes: 5},
			op:       "round",
			expected: "10m",
		},
		{
			name:     "round exactly half (to even) - negative",
			d:        Duration{Minutes: -7, Seconds: -30},
			multiple: Duration{Minutes: 5},
			op:       "round",
			expected: "-10m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Duration
			switch tt.op {
			case "round":
				result = tt.d.Round(tt.multiple)
			case "truncate":
				result = tt.d.Truncate(tt.multiple)
			}

			if got := result.String(); got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDuration_Abs(t *testing.T) {
	tests := []struct {
		name string
		d    Duration
		want Duration
	}{
		{
			name: "positive duration",
			d:    Duration{Days: 1},
			want: Duration{Days: 1},
		},
		{
			name: "negative duration",
			d:    Duration{Days: -1},
			want: Duration{Days: 1},
		},
		{
			name: "zero duration",
			d:    Duration{},
			want: Duration{},
		},
		{
			name: "complex negative duration",
			d:    Duration{Years: -1, Months: -2, Days: -3},
			want: Duration{Years: 1, Months: 2, Days: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Abs(); !got.Equal(tt.want) {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_Round_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		d        Duration
		multiple Duration
		want     Duration
	}{
		{
			name:     "round to zero multiple",
			d:        Duration{Minutes: 5},
			multiple: Duration{},
			want:     Duration{Minutes: 5},
		},
		{
			name:     "round to multiple with zero duration",
			d:        Duration{},
			multiple: Duration{Minutes: 5},
			want:     Duration{},
		},
		{
			name:     "round up minutes",
			d:        Duration{Minutes: 7},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: 5},
		},
		{
			name:     "round down minutes",
			d:        Duration{Minutes: 13},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: 15},
		},
		{
			name:     "round exactly half up",
			d:        Duration{Minutes: 2, Seconds: 30},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: 5},
		},
		{
			name:     "round exactly half down",
			d:        Duration{Minutes: 7, Seconds: 30},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: 10},
		},
		{
			name:     "round negative duration up",
			d:        Duration{Minutes: -7},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: -5},
		},
		{
			name:     "round negative duration down",
			d:        Duration{Minutes: -13},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: -15},
		},
		{
			name:     "round to larger unit",
			d:        Duration{Minutes: 55},
			multiple: Duration{Hours: 1},
			want:     Duration{Hours: 1},
		},
		{
			name:     "round to smaller unit",
			d:        Duration{Hours: 1, Minutes: 29},
			multiple: Duration{Minutes: 30},
			want:     Duration{Hours: 1, Minutes: 30},
		},
		{
			name:     "round with nanosecond precision",
			d:        Duration{Nanos: 1500},
			multiple: Duration{Nanos: 1000},
			want:     Duration{Nanos: 2000},
		},
		{
			name:     "round complex duration",
			d:        Duration{Hours: 2, Minutes: 45, Seconds: 30},
			multiple: Duration{Hours: 1, Minutes: 30},
			want:     Duration{Hours: 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.Round(tt.multiple)
			if !got.Equal(tt.want) {
				t.Errorf("Round() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_Truncate_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		d        Duration
		multiple Duration
		want     Duration
	}{
		{
			name:     "truncate to zero multiple",
			d:        Duration{Minutes: 5},
			multiple: Duration{},
			want:     Duration{Minutes: 5},
		},
		{
			name:     "truncate zero duration",
			d:        Duration{},
			multiple: Duration{Minutes: 5},
			want:     Duration{},
		},
		{
			name:     "truncate exact multiple",
			d:        Duration{Minutes: 15},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: 15},
		},
		{
			name:     "truncate non-multiple down",
			d:        Duration{Minutes: 17},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: 15},
		},
		{
			name:     "truncate negative duration",
			d:        Duration{Minutes: -17},
			multiple: Duration{Minutes: 5},
			want:     Duration{Minutes: -15},
		},
		{
			name:     "truncate to larger unit",
			d:        Duration{Minutes: 55},
			multiple: Duration{Hours: 1},
			want:     Duration{},
		},
		{
			name:     "truncate to smaller unit",
			d:        Duration{Hours: 1, Minutes: 29},
			multiple: Duration{Minutes: 30},
			want:     Duration{Hours: 1},
		},
		{
			name:     "truncate with nanosecond precision",
			d:        Duration{Nanos: 1999},
			multiple: Duration{Nanos: 1000},
			want:     Duration{Nanos: 1000},
		},
		{
			name:     "truncate complex duration",
			d:        Duration{Hours: 2, Minutes: 45, Seconds: 30},
			multiple: Duration{Hours: 1, Minutes: 30},
			want:     Duration{Hours: 1, Minutes: 30},
		},
		{
			name:     "truncate with remainder",
			d:        Duration{Hours: 5, Minutes: 45},
			multiple: Duration{Hours: 2},
			want:     Duration{Hours: 4},
		},
		{
			name:     "truncate near boundary",
			d:        Duration{Hours: 1, Minutes: 59, Seconds: 59, Nanos: 999999999},
			multiple: Duration{Hours: 1},
			want:     Duration{Hours: 1},
		},
		{
			name:     "truncate mixed units",
			d:        Duration{Days: 1, Hours: 13},
			multiple: Duration{Hours: 12},
			want:     Duration{Days: 1, Hours: 12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.Truncate(tt.multiple)
			if !got.Equal(tt.want) {
				t.Errorf("Truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDurationComparisons(t *testing.T) {
	tests := []struct {
		name    string
		d1      Duration
		d2      Duration
		less    bool
		equal   bool
		greater bool
	}{
		{
			name:    "equal durations",
			d1:      Duration{Days: 1},
			d2:      Duration{Days: 1},
			less:    false,
			equal:   true,
			greater: false,
		},
		{
			name:    "less than",
			d1:      Duration{Days: 1},
			d2:      Duration{Days: 2},
			less:    true,
			equal:   false,
			greater: false,
		},
		{
			name:    "greater than",
			d1:      Duration{Days: 2},
			d2:      Duration{Days: 1},
			less:    false,
			equal:   false,
			greater: true,
		},
		{
			name:    "complex comparison",
			d1:      Duration{Years: 1, Months: 1},
			d2:      Duration{Months: 13},
			less:    false,
			equal:   true,
			greater: false,
		},
		{
			name:    "negative vs positive",
			d1:      Duration{Days: -1},
			d2:      Duration{Days: 1},
			less:    true,
			equal:   false,
			greater: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d1.Less(tt.d2); got != tt.less {
				t.Errorf("Less() = %v, want %v", got, tt.less)
			}
			if got := tt.d1.Equal(tt.d2); got != tt.equal {
				t.Errorf("Equal() = %v, want %v", got, tt.equal)
			}
			if got := tt.d1.Greater(tt.d2); got != tt.greater {
				t.Errorf("Greater() = %v, want %v", got, tt.greater)
			}
			if got := tt.d1.LessOrEqual(tt.d2); got != (tt.less || tt.equal) {
				t.Errorf("LessOrEqual() = %v, want %v", got, tt.less || tt.equal)
			}
			if got := tt.d1.GreaterOrEqual(tt.d2); got != (tt.greater || tt.equal) {
				t.Errorf("GreaterOrEqual() = %v, want %v", got, tt.greater || tt.equal)
			}
		})
	}
}
