package hdur

import "testing"

func BenchmarkParseDuration(b *testing.B) {
	inputs := []string{
		"2h",
		"1y 2mo 3d",
		"1h 30m 45s",
		"500ms",
		"2.5h",
		"1 year 2 months 3 days and 4 hours",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		_, err := ParseDuration(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMustParseDuration_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParseDuration() did not panic with invalid input")
		}
	}()
	MustParseDuration("invalid duration")
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Duration
		wantErr  bool
	}{
		{
			name:     "simple seconds",
			input:    "30s",
			expected: Duration{Seconds: 30},
		},
		{
			name:     "complex duration",
			input:    "1 year 2 months 3 days 4 hours 5 minutes 6 seconds",
			expected: Duration{Years: 1, Months: 2, Days: 3, Hours: 4, Minutes: 5, Seconds: 6},
		},
		{
			name:     "multiple units with and",
			input:    "2 days and 4 hours",
			expected: Duration{Days: 2, Hours: 4},
		},
		{
			name:     "mixed case and plural",
			input:    "1 YEAR 1 Month 1 DAY",
			expected: Duration{Years: 1, Months: 1, Days: 1},
		},
		{
			name:     "alternative unit names",
			input:    "1 yr 1 mo 1 d",
			expected: Duration{Years: 1, Months: 1, Days: 1},
		},
		{
			name:     "weeks and fortnights",
			input:    "2 weeks 1 fortnight",
			expected: Duration{Days: 28}, // 14 days + 14 days
		},
		{
			name:     "nanoseconds",
			input:    "500ns",
			expected: Duration{Nanos: 500},
		},
		{
			name:     "microseconds",
			input:    "500us",
			expected: Duration{Nanos: 500000},
		},
		{
			name:     "milliseconds",
			input:    "500ms",
			expected: Duration{Nanos: 500000000},
		},
		{
			name:    "invalid unit",
			input:   "5 invalid",
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   "invalid days",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "just invalid",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got != tt.expected {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseDurationFromNanos(t *testing.T) {
	tests := []struct {
		name  string
		nanos int64
		want  Duration
	}{
		{
			name:  "zero",
			nanos: 0,
			want:  Duration{},
		},
		{
			name:  "one second",
			nanos: 1_000_000_000,
			want:  Duration{Seconds: 1},
		},
		{
			name:  "negative one second",
			nanos: -1_000_000_000,
			want:  Duration{Seconds: -1},
		},
		{
			name:  "complex duration",
			nanos: 90_061_000_000_000, // 25 hours 1 minute 1 second
			want:  Duration{Days: 1, Hours: 1, Minutes: 1, Seconds: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseDurationFromNanos(tt.nanos)
			if !got.Equal(tt.want) {
				t.Errorf("ParseDurationFromNanos() = %v, want %v", got, tt.want)
			}
		})
	}
}
