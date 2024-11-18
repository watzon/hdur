# hdur (/eɪtʃ.dɜr/)

## Human-Friendly Durations for Go

[![CI](https://github.com/watzon/hdur/actions/workflows/ci.yml/badge.svg)](https://github.com/watzon/hdur/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/watzon/hdur/branch/main/graph/badge.svg)](https://codecov.io/gh/watzon/hdur)
[![Go Report Card](https://goreportcard.com/badge/github.com/watzon/hdur)](https://goreportcard.com/report/github.com/watzon/hdur)
[![GoDoc](https://pkg.go.dev/badge/github.com/watzon/hdur)](https://pkg.go.dev/github.com/watzon/hdur)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`hdur` (Human Duration) is a flexible and intuitive duration library for Go that extends the standard `time.Duration` with calendar-aware features and human-readable parsing. It's designed to make working with durations as natural as possible.

## Features

- 📅 **Calendar-aware duration handling**
  - Properly handles months and years
  - Accounts for varying month lengths
  - Preserves day-of-month when possible
- 🔍 **Human-readable parsing**
  - Parse durations like "1 year 2 months 3 days"
  - Flexible unit names (year/yr/y, month/mo/m, etc.)
  - Handles both full and abbreviated units
- ⚡ **Rich functionality**
  - Mathematical operations (add, multiply, divide)
  - Rounding and truncation
  - Comparison operations
  - Time-based calculations
- 🔄 **Serialization support**
  - JSON marshaling/unmarshaling
  - SQL scanning/valuing
  - Custom format strings

## Installation

```bash
go get github.com/watzon/hdur
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/watzon/hdur"
)

func main() {
    // Parse a duration string
    d, err := hdur.ParseDuration("1 year 2 months 3 days")
    if err != nil {
        panic(err)
    }

    // Add to current time
    future := d.Add(time.Now())
    fmt.Printf("Future date: %v\n", future)

    // Create durations using constructors
    day := hdur.Days(1)
    week := hdur.Weeks(1)
    if week.Equal(hdur.Days(7)) {
        fmt.Println("1 week equals 7 days")
    }

    // Mathematical operations
    doubled := day.Mul(2)
    fmt.Printf("Two days: %v\n", doubled)

    // Format durations
    threeDays := hdur.Days(3)
    fmt.Println(threeDays.Format("%d days"))  // "3 days"

    // Use with time.Time
    start := time.Now()
    time.Sleep(time.Second)
    elapsed := hdur.Since(start)
    fmt.Printf("Operation took: %v\n", elapsed)
}
```

## Usage Examples

### Parsing Durations

```go
// Multiple ways to specify the same duration
d1, _ := hdur.ParseDuration("1 year 2 months")
d2, _ := hdur.ParseDuration("1y 2mo")
d3, _ := hdur.ParseDuration("1yr 2mon")

// Parse with conjunctions
d4, _ := hdur.ParseDuration("1 year and 2 months")

// Various time units
d5, _ := hdur.ParseDuration("2weeks 4days 12hours 30minutes 45seconds")
```

### Creating Durations

```go
// Using constructors
hour := hdur.Hours(1)
day := hdur.Days(1)
week := hdur.Weeks(1)
month := hdur.Months(1)
year := hdur.Years(1)

// Using common constants
thirtySeconds := hdur.Seconds30
twentyFourHours := hdur.Hours24
sixMonths := hdur.Months6
```

### Time Operations

```go
// Add duration to time
future := hdur.Months(3).Add(time.Now())

// Get duration between times
start := time.Now()
// ... do something ...
elapsed := hdur.Since(start)

// Calculate time until future date
deadline := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
remaining := hdur.Until(deadline)
```

### Mathematical Operations

```go
// Multiply duration
double := hdur.Days(1).Mul(2)
half := hdur.Hours(1).Div(2)

// Round duration
d := hdur.MustParseDuration("1h 30m")
rounded := d.Round(hdur.Hours(1)) // 2h

// Compare durations
if d1.Less(d2) {
    fmt.Println("d1 is shorter than d2")
}
```

### Formatting

```go
d := hdur.MustParseDuration("1 year 2 months 3 days 4 hours")

// Default format
fmt.Println(d) // "1y 2mo 3d 4h"

// Custom format
fmt.Println(d.Format("%y years %M months %d days %h hours"))
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE.md) for details.
