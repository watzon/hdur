package hdur

import (
	"encoding/json"
	"testing"
)

func BenchmarkDurationString(b *testing.B) {
	durations := []Duration{
		Hours(2),
		ParseMust("1y 2mo 3d"),
		Minutes(90),
		Milliseconds(500),
		ParseMust("2.5h"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d := durations[i%len(durations)]
		_ = d.String()
	}
}

func BenchmarkDurationFormat(b *testing.B) {
	d := ParseMust("1y 2mo 3d 4h 5m 6s")
	formats := []string{
		"%y years %M months %d days",
		"%h hours %m minutes %s seconds",
		"%y-%y-%y",
		"%s.%f seconds",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		format := formats[i%len(formats)]
		_ = d.Format(format)
	}
}

func BenchmarkDurationJSON(b *testing.B) {
	type wrapper struct {
		D Duration `json:"duration"`
	}
	w := wrapper{D: ParseMust("1y 2mo 3d")}
	data, _ := json.Marshal(w)

	b.Run("marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(w)
		}
	})

	b.Run("unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var w2 wrapper
			_ = json.Unmarshal(data, &w2)
		}
	})
}

func TestDuration_String(t *testing.T) {
	tests := []struct {
		name     string
		duration Duration
		want     string
	}{
		{
			name:     "zero duration",
			duration: Duration{},
			want:     "0s",
		},
		{
			name: "full duration",
			duration: Duration{
				Years: 1, Months: 2, Days: 3,
				Hours: 4, Minutes: 5, Seconds: 6,
			},
			want: "1y 2mo 3d 4h 5m 6s",
		},
		{
			name: "only large units",
			duration: Duration{
				Years: 1, Months: 2,
			},
			want: "1y 2mo",
		},
		{
			name: "only small units",
			duration: Duration{
				Minutes: 5, Seconds: 6,
			},
			want: "5m 6s",
		},
		{
			name: "nanoseconds",
			duration: Duration{
				Nanos: 500000000,
			},
			want: "500ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.duration.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_String_Edge_Cases(t *testing.T) {
	tests := []struct {
		name string
		d    Duration
		want string
	}{
		{
			name: "all negative components",
			d:    Duration{Years: -1, Months: -2, Days: -3, Hours: -4, Minutes: -5, Seconds: -6, Nanos: -7},
			want: "-1y 2mo 3d 4h 5m 6s 7ns",
		},
		{
			name: "only nanoseconds",
			d:    Duration{Nanos: 500},
			want: "500ns",
		},
		{
			name: "microseconds exact",
			d:    Duration{Nanos: 3000},
			want: "3µs",
		},
		{
			name: "microseconds with fraction",
			d:    Duration{Nanos: 3500},
			want: "3.500µs",
		},
		{
			name: "milliseconds exact",
			d:    Duration{Nanos: 3000000},
			want: "3ms",
		},
		{
			name: "milliseconds with fraction",
			d:    Duration{Nanos: 3500000},
			want: "3.500ms",
		},
		{
			name: "zero with nanoseconds",
			d:    Duration{Nanos: 0},
			want: "0s",
		},
		{
			name: "negative nanoseconds only",
			d:    Duration{Nanos: -500},
			want: "-500ns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_Format(t *testing.T) {
	d := ParseMust("1y 2mo 3d 4h 5m 6s")
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "full format",
			format:   "%y years %M months %d days",
			expected: "1 years 2 months 3 days",
		},
		{
			name:     "partial format",
			format:   "%y years and %d days",
			expected: "1 years and 3 days",
		},
		{
			name:     "repeated specifiers",
			format:   "%y-%y-%y",
			expected: "1-1-1",
		},
		{
			name:     "fractional seconds",
			format:   "%s.%f seconds",
			expected: "6.000000000 seconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.Format(tt.format)
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDurationJSON(t *testing.T) {
	type wrapper struct {
		D Duration `json:"duration"`
	}

	tests := []struct {
		name     string
		d        Duration
		expected string
		input    string
	}{
		{
			name:     "simple duration",
			d:        Hours(2),
			expected: `{"duration":"2h"}`,
			input:    `{"duration":"2h"}`,
		},
		{
			name:     "complex duration",
			d:        ParseMust("1y 2mo 3d"),
			expected: `{"duration":"1y 2mo 3d"}`,
			input:    `{"duration":"1y 2mo 3d"}`,
		},
		{
			name:     "zero duration",
			d:        Duration{},
			expected: `{"duration":"0s"}`,
			input:    `{"duration":"0s"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			w := wrapper{D: tt.d}
			got, err := json.Marshal(w)
			if err != nil {
				t.Errorf("marshal error: %v", err)
			}
			if string(got) != tt.expected {
				t.Errorf("marshal: got %v, want %v", string(got), tt.expected)
			}

			// Test unmarshaling
			var w2 wrapper
			err = json.Unmarshal([]byte(tt.input), &w2)
			if err != nil {
				t.Errorf("unmarshal error: %v", err)
			}
			if !w2.D.Equal(tt.d) {
				t.Errorf("unmarshal: got %v, want %v", w2.D, tt.d)
			}
		})
	}
}

func TestDurationSQL(t *testing.T) {
	tests := []struct {
		name     string
		d        Duration
		expected string
	}{
		{
			name:     "simple duration",
			d:        Hours(2),
			expected: "2h",
		},
		{
			name:     "complex duration",
			d:        ParseMust("1y 2mo 3d"),
			expected: "1y 2mo 3d",
		},
		{
			name:     "zero duration",
			d:        Duration{},
			expected: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Value
			got, err := tt.d.Value()
			if err != nil {
				t.Errorf("value error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("value: got %v, want %v", got, tt.expected)
			}

			// Test Scan
			var d2 Duration
			err = d2.Scan(tt.expected)
			if err != nil {
				t.Errorf("scan error: %v", err)
			}
			if !d2.Equal(tt.d) {
				t.Errorf("scan: got %v, want %v", d2, tt.d)
			}

			// Test Scan with []byte
			err = d2.Scan([]byte(tt.expected))
			if err != nil {
				t.Errorf("scan bytes error: %v", err)
			}
			if !d2.Equal(tt.d) {
				t.Errorf("scan bytes: got %v, want %v", d2, tt.d)
			}

			// Test Scan with nil
			err = d2.Scan(nil)
			if err != nil {
				t.Errorf("scan nil error: %v", err)
			}
			if !d2.Equal(Duration{}) {
				t.Errorf("scan nil: got %v, want zero duration", d2)
			}

			// Test Scan with invalid type
			err = d2.Scan(123)
			if err == nil {
				t.Error("expected error scanning invalid type")
			}
		})
	}
}

func TestDuration_UnmarshalJSON_Edge_Cases(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Duration
		wantErr bool
	}{
		{
			name:    "invalid JSON",
			json:    `{"duration": "1h"}`,
			wantErr: true,
		},
		{
			name:    "invalid duration string",
			json:    `"invalid"`,
			wantErr: true,
		},
		{
			name:    "empty string",
			json:    `""`,
			wantErr: true, // Changed: empty string should error
		},
		{
			name:    "null",
			json:    `null`,
			wantErr: true,
		},
		{
			name:    "number instead of string",
			json:    `123`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Duration
			err := d.UnmarshalJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !d.Equal(tt.want) {
				t.Errorf("UnmarshalJSON() = %v, want %v", d, tt.want)
			}
		})
	}
}
