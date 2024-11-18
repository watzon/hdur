//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"
	"time"

	"github.com/watzon/hdur"
)

func main() {
	c := make(chan struct{}, 0)

	// Register functions
	js.Global().Set("parseDuration", js.FuncOf(parseDuration))
	js.Global().Set("formatDuration", js.FuncOf(formatDuration))
	js.Global().Set("formatWithPattern", js.FuncOf(formatWithPattern))
	js.Global().Set("addDuration", js.FuncOf(addDuration))
	js.Global().Set("subDuration", js.FuncOf(subDuration))
	js.Global().Set("toStandard", js.FuncOf(toStandard))
	js.Global().Set("isZeroDuration", js.FuncOf(isZeroDuration))
	js.Global().Set("untilTime", js.FuncOf(untilTime))
	js.Global().Set("sinceTime", js.FuncOf(sinceTime))
	js.Global().Set("betweenTimes", js.FuncOf(betweenTimes))
	js.Global().Set("equalDurations", js.FuncOf(equalDurations))
	js.Global().Set("lessDurations", js.FuncOf(lessDurations))
	js.Global().Set("lessThanOrEqual", js.FuncOf(lessThanOrEqual))
	js.Global().Set("greaterDurations", js.FuncOf(greaterDurations))
	js.Global().Set("greaterThanOrEqual", js.FuncOf(greaterThanOrEqual))
	js.Global().Set("absDuration", js.FuncOf(absDuration))
	js.Global().Set("mulDuration", js.FuncOf(mulDuration))
	js.Global().Set("divDuration", js.FuncOf(divDuration))
	js.Global().Set("roundDuration", js.FuncOf(roundDuration))
	js.Global().Set("truncateDuration", js.FuncOf(truncateDuration))
	js.Global().Set("fromNanoseconds", js.FuncOf(fromNanoseconds))
	js.Global().Set("fromMicroseconds", js.FuncOf(fromMicroseconds))
	js.Global().Set("fromMilliseconds", js.FuncOf(fromMilliseconds))
	js.Global().Set("fromSeconds", js.FuncOf(fromSeconds))
	js.Global().Set("fromMinutes", js.FuncOf(fromMinutes))
	js.Global().Set("fromHours", js.FuncOf(fromHours))
	js.Global().Set("fromDays", js.FuncOf(fromDays))
	js.Global().Set("fromWeeks", js.FuncOf(fromWeeks))
	js.Global().Set("fromMonths", js.FuncOf(fromMonths))
	js.Global().Set("fromYears", js.FuncOf(fromYears))
	js.Global().Set("inHours", js.FuncOf(inHours))
	js.Global().Set("inMinutes", js.FuncOf(inMinutes))
	js.Global().Set("inSeconds", js.FuncOf(inSeconds))
	js.Global().Set("inNanoseconds", js.FuncOf(inNanoseconds))
	js.Global().Set("inMonths", js.FuncOf(inMonths))
	js.Global().Set("inYears", js.FuncOf(inYears))

	<-c
}

func parseDuration(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "parse duration requires a string argument",
		}
	}

	durStr := args[0].String()
	dur, err := hdur.ParseDuration(durStr)
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"years":   dur.Years,
		"months":  dur.Months,
		"days":    dur.Days,
		"hours":   dur.Hours,
		"minutes": dur.Minutes,
		"seconds": dur.Seconds,
		"nanos":   dur.Nanos,
	}
}

func formatDuration(this js.Value, args []js.Value) interface{} {
	if len(args) < 7 {
		return map[string]interface{}{
			"error": "format duration requires years, months, days, hours, minutes, seconds, and nanos arguments",
		}
	}

	dur := hdur.Duration{
		Years:   args[0].Int(),
		Months:  args[1].Int(),
		Days:    args[2].Int(),
		Hours:   args[3].Int(),
		Minutes: args[4].Int(),
		Seconds: args[5].Int(),
		Nanos:   args[6].Int(),
	}

	return dur.String()
}

func formatWithPattern(this js.Value, args []js.Value) interface{} {
	if len(args) < 8 {
		return map[string]interface{}{
			"error": "format with pattern requires years, months, days, hours, minutes, seconds, nanos, and pattern arguments",
		}
	}

	dur := hdur.Duration{
		Years:   args[0].Int(),
		Months:  args[1].Int(),
		Days:    args[2].Int(),
		Hours:   args[3].Int(),
		Minutes: args[4].Int(),
		Seconds: args[5].Int(),
		Nanos:   args[6].Int(),
	}
	pattern := args[7].String()

	return dur.Format(pattern)
}

func addDuration(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "add duration requires a time argument in ISO format",
		}
	}

	timeStr := args[0].String()
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return map[string]interface{}{
			"error": "invalid time format, use ISO format",
		}
	}

	dur := getDurationFromArgs(args, 1)
	result := dur.Add(t)

	return result.Format(time.RFC3339)
}

func subDuration(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{
			"error": "subtract duration requires two time arguments in ISO format",
		}
	}

	t1Str := args[0].String()
	t2Str := args[1].String()

	t1, err := time.Parse(time.RFC3339, t1Str)
	if err != nil {
		return map[string]interface{}{
			"error": "invalid time format for first argument, use ISO format",
		}
	}

	t2, err := time.Parse(time.RFC3339, t2Str)
	if err != nil {
		return map[string]interface{}{
			"error": "invalid time format for second argument, use ISO format",
		}
	}

	result := hdur.Sub(t1, t2)

	return map[string]interface{}{
		"years":   result.Years,
		"months":  result.Months,
		"days":    result.Days,
		"hours":   result.Hours,
		"minutes": result.Minutes,
		"seconds": result.Seconds,
		"nanos":   result.Nanos,
	}
}

func toStandard(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	std := dur.ToStandard()
	return std.Nanoseconds()
}

func getDurationFromArgs(args []js.Value, startIdx int) hdur.Duration {
	return hdur.Duration{
		Years:   args[startIdx].Int(),
		Months:  args[startIdx+1].Int(),
		Days:    args[startIdx+2].Int(),
		Hours:   args[startIdx+3].Int(),
		Minutes: args[startIdx+4].Int(),
		Seconds: args[startIdx+5].Int(),
		Nanos:   args[startIdx+6].Int(),
	}
}

func isZeroDuration(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	return dur.IsZero()
}

func untilTime(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "until requires a time argument in ISO format",
		}
	}

	timeStr := args[0].String()
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return map[string]interface{}{
			"error": "invalid time format, use ISO format",
		}
	}

	result := hdur.Until(t)
	return durationToMap(result)
}

func sinceTime(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{
			"error": "since requires a time argument in ISO format",
		}
	}

	timeStr := args[0].String()
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return map[string]interface{}{
			"error": "invalid time format, use ISO format",
		}
	}

	result := hdur.Since(t)
	return durationToMap(result)
}

func betweenTimes(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{
			"error": "between requires two time arguments in ISO format",
		}
	}

	t1Str := args[0].String()
	t2Str := args[1].String()

	t1, err := time.Parse(time.RFC3339, t1Str)
	if err != nil {
		return map[string]interface{}{
			"error": "invalid time format for first argument, use ISO format",
		}
	}

	t2, err := time.Parse(time.RFC3339, t2Str)
	if err != nil {
		return map[string]interface{}{
			"error": "invalid time format for second argument, use ISO format",
		}
	}

	result := hdur.Between(t1, t2)
	return durationToMap(result)
}

func durationToMap(d hdur.Duration) map[string]interface{} {
	return map[string]interface{}{
		"years":   d.Years,
		"months":  d.Months,
		"days":    d.Days,
		"hours":   d.Hours,
		"minutes": d.Minutes,
		"seconds": d.Seconds,
		"nanos":   d.Nanos,
	}
}

func equalDurations(this js.Value, args []js.Value) interface{} {
	if len(args) < 14 {
		return false
	}

	d1 := getDurationFromArgs(args, 0)
	d2 := getDurationFromArgs(args, 7)

	return d1.Equal(d2)
}

func lessDurations(this js.Value, args []js.Value) interface{} {
	if len(args) < 14 {
		return false
	}

	d1 := getDurationFromArgs(args, 0)
	d2 := getDurationFromArgs(args, 7)

	return d1.Less(d2)
}

func lessThanOrEqual(this js.Value, args []js.Value) interface{} {
	if len(args) < 14 {
		return false
	}

	d1 := getDurationFromArgs(args, 0)
	d2 := getDurationFromArgs(args, 7)

	return d1.LessOrEqual(d2)
}

func greaterDurations(this js.Value, args []js.Value) interface{} {
	if len(args) < 14 {
		return false
	}

	d1 := getDurationFromArgs(args, 0)
	d2 := getDurationFromArgs(args, 7)

	return d1.Greater(d2)
}

func greaterThanOrEqual(this js.Value, args []js.Value) interface{} {
	if len(args) < 14 {
		return false
	}

	d1 := getDurationFromArgs(args, 0)
	d2 := getDurationFromArgs(args, 7)

	return d1.GreaterOrEqual(d2)
}

func absDuration(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	result := dur.Abs()
	return durationToMap(result)
}

func mulDuration(this js.Value, args []js.Value) interface{} {
	if len(args) < 8 {
		return map[string]interface{}{
			"error": "multiply duration requires a duration and a factor",
		}
	}

	dur := getDurationFromArgs(args, 0)
	factor := args[7].Float()

	result := dur.Mul(factor)
	return durationToMap(result)
}

func divDuration(this js.Value, args []js.Value) interface{} {
	if len(args) < 8 {
		return map[string]interface{}{
			"error": "divide duration requires a duration and a divisor",
		}
	}

	dur := getDurationFromArgs(args, 0)
	divisor := args[7].Float()

	if divisor == 0 {
		return map[string]interface{}{
			"error": "division by zero",
		}
	}

	result := dur.Div(divisor)
	return durationToMap(result)
}

func roundDuration(this js.Value, args []js.Value) interface{} {
	if len(args) < 14 {
		return map[string]interface{}{
			"error": "round duration requires two durations",
		}
	}

	d1 := getDurationFromArgs(args, 0)
	d2 := getDurationFromArgs(args, 7)

	result := d1.Round(d2)
	return durationToMap(result)
}

func truncateDuration(this js.Value, args []js.Value) interface{} {
	if len(args) < 14 {
		return map[string]interface{}{
			"error": "truncate duration requires two durations",
		}
	}

	d1 := getDurationFromArgs(args, 0)
	d2 := getDurationFromArgs(args, 7)

	result := d1.Truncate(d2)
	return durationToMap(result)
}

func fromNanoseconds(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	n := args[0].Float()
	return durationToMap(hdur.Nanoseconds(n))
}

func fromMicroseconds(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	µs := args[0].Float()
	return durationToMap(hdur.Microseconds(µs))
}

func fromMilliseconds(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	ms := args[0].Float()
	return durationToMap(hdur.Milliseconds(ms))
}

func fromSeconds(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	s := args[0].Float()
	return durationToMap(hdur.Seconds(s))
}

func fromMinutes(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	m := args[0].Float()
	return durationToMap(hdur.Minutes(m))
}

func fromHours(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	h := args[0].Float()
	return durationToMap(hdur.Hours(h))
}

func fromDays(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	d := args[0].Float()
	return durationToMap(hdur.Days(d))
}

func fromWeeks(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	w := args[0].Float()
	return durationToMap(hdur.Weeks(w))
}

func fromMonths(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	m := args[0].Float()
	return durationToMap(hdur.Months(m))
}

func fromYears(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return durationToMap(hdur.Duration{})
	}
	y := args[0].Float()
	return durationToMap(hdur.Years(y))
}

func inHours(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	return dur.InHours()
}

func inMinutes(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	return dur.InMinutes()
}

func inSeconds(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	return dur.InSeconds()
}

func inNanoseconds(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	return dur.InNanoseconds()
}

func inMonths(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	return dur.InMonths()
}

func inYears(this js.Value, args []js.Value) interface{} {
	dur := getDurationFromArgs(args, 0)
	return dur.InYears()
}
