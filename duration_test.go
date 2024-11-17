package hdur

import (
	"testing"
)

func TestDuration_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		input    Duration
		expected Duration
	}{
		{
			name: "normalize seconds",
			input: Duration{
				Seconds: 65,
			},
			expected: Duration{
				Minutes: 1,
				Seconds: 5,
			},
		},
		{
			name: "normalize minutes",
			input: Duration{
				Minutes: 125,
			},
			expected: Duration{
				Hours:   2,
				Minutes: 5,
			},
		},
		{
			name: "normalize hours",
			input: Duration{
				Hours: 25,
			},
			expected: Duration{
				Days:  1,
				Hours: 1,
			},
		},
		{
			name: "normalize months",
			input: Duration{
				Months: 14,
			},
			expected: Duration{
				Years:  1,
				Months: 2,
			},
		},
		{
			name: "normalize nanoseconds",
			input: Duration{
				Nanos: 1500000000, // 1.5 seconds in nanoseconds
			},
			expected: Duration{
				Seconds: 1,
				Nanos:   500000000,
			},
		},
		{
			name: "normalize complex",
			input: Duration{
				Months:  25,
				Days:    45,
				Hours:   30,
				Minutes: 70,
				Seconds: 75,
				Nanos:   2500000000,
			},
			expected: Duration{
				Years:   2,
				Months:  1,
				Days:    46,
				Hours:   7,
				Minutes: 11,
				Seconds: 17,
				Nanos:   500000000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.normalize()
			if tt.input != tt.expected {
				t.Errorf("normalize() got = %v, want %v", tt.input, tt.expected)
			}
		})
	}
}
