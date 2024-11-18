package hdur

import (
	"math"
	"testing"
	"time"
)

func BenchmarkDurationAdd(b *testing.B) {
	durations := []Duration{
		Hours(2),
		MustParseDuration("1y 2mo 3d"),
		Minutes(90),
		Milliseconds(500),
		MustParseDuration("2.5h"),
	}
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := durations[i%len(durations)]
		_ = d.Add(now)
	}
}

func TestGetDaysInMonth(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		want  int
	}{
		{
			name:  "January",
			year:  2024,
			month: time.January,
			want:  31,
		},
		{
			name:  "February non-leap year",
			year:  2023,
			month: time.February,
			want:  28,
		},
		{
			name:  "February leap year",
			year:  2024,
			month: time.February,
			want:  29,
		},
		{
			name:  "April",
			year:  2024,
			month: time.April,
			want:  30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDaysInMonth(tt.year, tt.month); got != tt.want {
				t.Errorf("getDaysInMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSub(t *testing.T) {
	tests := []struct {
		name     string
		time1    time.Time
		time2    time.Time
		expected Duration
	}{
		{
			name:  "same month different days",
			time1: time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC),
			time2: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: Duration{
				Days: 14,
			},
		},
		{
			name:  "different months same year",
			time1: time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
			time2: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: Duration{
				Months: 2,
			},
		},
		{
			name:  "different years",
			time1: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			time2: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: Duration{
				Years: 1,
			},
		},
		{
			name:  "complex difference",
			time1: time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC),
			time2: time.Date(2023, time.January, 1, 12, 15, 30, 0, time.UTC),
			expected: Duration{
				Years:   1,
				Months:  2,
				Days:    14,
				Hours:   2,
				Minutes: 15,
				Seconds: 15,
			},
		},
		{
			name:  "reverse order (should preserve sign)",
			time1: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			time2: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: Duration{
				Years: -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sub(tt.time1, tt.time2)
			if got != tt.expected {
				t.Errorf("Sub() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDuration_Sub_Edge_Cases(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		t1   time.Time
		t2   time.Time
		want Duration
	}{
		{
			name: "same time",
			t1:   now,
			t2:   now,
			want: Duration{},
		},
		{
			name: "one nanosecond difference",
			t1:   now,
			t2:   now.Add(1),
			want: Duration{Nanos: 1},
		},
		{
			name: "crossing month boundary forward",
			t1:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			t2:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			want: Duration{Days: 1},
		},
		{
			name: "crossing month boundary backward",
			t1:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			t2:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			want: Duration{Days: -1},
		},
		{
			name: "crossing year boundary",
			t1:   time.Date(2023, 12, 31, 23, 59, 59, 999999999, time.UTC),
			t2:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want: Duration{Nanos: 1},
		},
		{
			name: "leap year February",
			t1:   time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
			t2:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			want: Duration{Days: 2},
		},
		{
			name: "complex duration",
			t1:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			t2:   time.Date(2024, 3, 15, 12, 30, 45, 123456789, time.UTC),
			want: Duration{
				Years:   1,
				Months:  2,
				Days:    14,
				Hours:   12,
				Minutes: 30,
				Seconds: 45,
				Nanos:   123456789,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sub(tt.t2, tt.t1)
			if !got.Equal(tt.want) {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_Add(t *testing.T) {
	baseTime := time.Date(2023, time.January, 31, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		duration string
		start    time.Time
		expected time.Time
	}{
		{
			name:     "add one month to January 31",
			duration: "1 month",
			start:    baseTime,
			expected: time.Date(2023, time.February, 28, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "add two months to January 31",
			duration: "2 months",
			start:    baseTime,
			expected: time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "add one year",
			duration: "1 year",
			start:    baseTime,
			expected: time.Date(2024, time.January, 31, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "complex duration",
			duration: "1 year 1 month 1 day 1 hour",
			start:    baseTime,
			expected: time.Date(2024, time.March, 1, 13, 0, 0, 0, time.UTC),
		},
		{
			name:     "leap year February",
			duration: "1 year 1 month",
			start:    time.Date(2024, time.January, 31, 12, 0, 0, 0, time.UTC),
			expected: time.Date(2025, time.February, 28, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ParseDuration(tt.duration)
			if err != nil {
				t.Fatalf("Failed to parse duration: %v", err)
			}
			got := d.Add(tt.start)
			if !got.Equal(tt.expected) {
				t.Errorf("Add() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDurationAddition(t *testing.T) {
	baseTime := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		duration Duration
		expected time.Time
	}{
		{
			name:     "add nanoseconds",
			duration: Nanoseconds(500000000),
			expected: baseTime.Add(500 * time.Millisecond),
		},
		{
			name:     "add seconds",
			duration: Seconds(90),
			expected: baseTime.Add(90 * time.Second),
		},
		{
			name:     "add complex duration",
			duration: Hours(25),
			expected: baseTime.AddDate(0, 0, 1).Add(time.Hour),
		},
		{
			name:     "add months",
			duration: Months(14),
			expected: baseTime.AddDate(1, 2, 0),
		},
		{
			name:     "add years",
			duration: Years(2),
			expected: baseTime.AddDate(2, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.duration.Add(baseTime)
			if !got.Equal(tt.expected) {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDurationTimeOperations(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-time.Hour)
	inOneHour := now.Add(time.Hour)

	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected float64
		delta    float64
	}{
		{
			name:     "since one hour ago",
			t1:       oneHourAgo,
			expected: 1.0,
			delta:    0.01,
		},
		{
			name:     "until one hour from now",
			t1:       inOneHour,
			expected: 1.0,
			delta:    0.01,
		},
		{
			name:     "between times",
			t1:       oneHourAgo,
			t2:       inOneHour,
			expected: 2.0,
			delta:    0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got float64
			switch {
			case tt.name == "since one hour ago":
				got = Since(tt.t1).InHours()
			case tt.name == "until one hour from now":
				got = Until(tt.t1).InHours()
			default:
				got = Between(tt.t1, tt.t2).InHours()
			}

			if math.Abs(got-tt.expected) > tt.delta {
				t.Errorf("got %v, want %v (±%v)", got, tt.expected, tt.delta)
			}
		})
	}
}
