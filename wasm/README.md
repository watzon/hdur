# @watzon/hdur

A powerful duration manipulation library with nanosecond precision, built with Go and WebAssembly.

## Features

- 🎯 Precise duration calculations from nanoseconds to years
- 🔄 Convert between different time units
- 📅 Add durations to dates
- 🧮 Mathematical operations (add, subtract, multiply, divide)
- 🎨 Custom format strings
- 📊 Compare durations
- 💪 Type-safe with full TypeScript support

## Installation

```bash
npm install @watzon/hdur
```

## Usage

```typescript
import { Duration, initWasm } from '@watzon/hdur';

async function example() {
    // Initialize WASM (required before using any functionality)
    await initWasm();

    // Parse a duration string
    const dur = Duration.parse('1 year 2 months 3 days');

    // Create durations using constructor functions
    const twoHours = Duration.fromHours(2);
    const threeWeeks = Duration.fromWeeks(3);

    // Add duration to a date
    const future = dur.add(new Date());

    // Compare durations
    if (twoHours.lessThan(threeWeeks)) {
        console.log('2 hours is less than 3 weeks');
    }

    // Mathematical operations
    const doubled = dur.multiply(2);
    const halved = dur.divide(2);

    // Convert to different units
    console.log(threeWeeks.inDays());  // ~21
    console.log(threeWeeks.inHours()); // ~504

    // Format with custom pattern
    console.log(dur.format('%y years %M months %d days'));
    // Output: "1 years 2 months 3 days"

    // Use common duration constants
    const minute = Duration.Minute;
    const day = Duration.Day;
    const year = Duration.Year;
}
```

## API Reference

### Constructor Functions

- `Duration.fromNanoseconds(n: number)`
- `Duration.fromMicroseconds(µs: number)`
- `Duration.fromMilliseconds(ms: number)`
- `Duration.fromSeconds(s: number)`
- `Duration.fromMinutes(m: number)`
- `Duration.fromHours(h: number)`
- `Duration.fromDays(d: number)`
- `Duration.fromWeeks(w: number)`
- `Duration.fromMonths(m: number)`
- `Duration.fromYears(y: number)`

### Instance Methods

- `toString()`: Human-readable string
- `format(pattern: string)`: Custom format
- `add(date: Date)`: Add to date
- `equals(other: Duration)`: Compare equality
- `lessThan(other: Duration)`: Compare less than
- `greaterThan(other: Duration)`: Compare greater than
- `abs()`: Absolute value
- `multiply(factor: number)`: Multiply duration
- `divide(divisor: number)`: Divide duration
- `round(multiple: Duration)`: Round to multiple
- `truncate(multiple: Duration)`: Truncate to multiple

### Conversion Methods

- `inHours()`: Convert to hours
- `inMinutes()`: Convert to minutes
- `inSeconds()`: Convert to seconds
- `inNanoseconds()`: Convert to nanoseconds
- `inMonths()`: Convert to months
- `inYears()`: Convert to years

### Common Duration Constants

- `Duration.Nanosecond`
- `Duration.Microsecond`
- `Duration.Millisecond`
- `Duration.Second`
- `Duration.Minute`
- `Duration.Hour`
- `Duration.Day`
- `Duration.Week`
- `Duration.Month`
- `Duration.Year`

## Format Patterns

The `format()` method supports the following specifiers:

- `%y`: years
- `%M`: months
- `%d`: days
- `%h`: hours
- `%m`: minutes
- `%s`: seconds
- `%f`: fractional seconds (without leading dot)

## License

MIT
