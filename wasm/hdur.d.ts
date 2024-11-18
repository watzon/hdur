/**
 * A Duration represents a period of time with nanosecond precision.
 * It can represent time units from years down to nanoseconds and
 * supports both positive and negative durations.
 */
declare class Duration {
    /** Number of years in the duration */
    years: number;
    /** Number of months in the duration */
    months: number;
    /** Number of days in the duration */
    days: number;
    /** Number of hours in the duration */
    hours: number;
    /** Number of minutes in the duration */
    minutes: number;
    /** Number of seconds in the duration */
    seconds: number;
    /** Number of nanoseconds in the duration */
    nanos: number;

    /**
     * Creates a new Duration instance.
     * All parameters are optional and default to 0.
     */
    constructor(
        years?: number,
        months?: number,
        days?: number,
        hours?: number,
        minutes?: number,
        seconds?: number,
        nanos?: number
    );

    // Instance methods
    /**
     * Returns a human-readable string representation of the duration.
     * Example: "1y 2mo 3d 4h 5m 6s"
     */
    toString(): string;

    /**
     * Formats the duration using a custom pattern.
     * Format specifiers:
     * - %y: years
     * - %M: months
     * - %d: days
     * - %h: hours
     * - %m: minutes
     * - %s: seconds
     * - %f: fractional seconds (without leading dot)
     * @param pattern The format pattern to use
     */
    format(pattern: string): string;

    /**
     * Adds this duration to a date and returns the resulting date.
     * @param date The date to add this duration to
     * @returns A new Date object representing the sum
     */
    add(date: Date): Date;

    /**
     * Checks if this duration is equal to another duration.
     * @param other The duration to compare with
     */
    equals(other: Duration): boolean;

    /**
     * Checks if this duration is less than another duration.
     * @param other The duration to compare with
     */
    lessThan(other: Duration): boolean;

    /**
     * Checks if this duration is less than or equal to another duration.
     * @param other The duration to compare with
     */
    lessOrEqual(other: Duration): boolean;

    /**
     * Checks if this duration is greater than another duration.
     * @param other The duration to compare with
     */
    greaterThan(other: Duration): boolean;

    /**
     * Checks if this duration is greater than or equal to another duration.
     * @param other The duration to compare with
     */
    greaterOrEqual(other: Duration): boolean;

    /**
     * Returns the absolute value of this duration.
     * For example, -3 hours becomes 3 hours.
     */
    abs(): Duration;

    /**
     * Multiplies this duration by a factor.
     * @param factor The number to multiply by
     */
    multiply(factor: number): Duration;

    /**
     * Divides this duration by a divisor.
     * @param divisor The number to divide by
     * @throws {Error} If divisor is 0
     */
    divide(divisor: number): Duration;

    /**
     * Rounds this duration to the nearest multiple of another duration.
     * @param multiple The duration to round to
     */
    round(multiple: Duration): Duration;

    /**
     * Truncates this duration to a multiple of another duration.
     * @param multiple The duration to truncate to
     */
    truncate(multiple: Duration): Duration;

    /**
     * Checks if this duration is zero (no time).
     */
    isZero(): boolean;

    /**
     * Converts this duration to a standard format string.
     */
    toStandard(): string;

    /**
     * Converts this duration to hours.
     * @returns The total number of hours (may be fractional)
     */
    inHours(): number;

    /**
     * Converts this duration to minutes.
     * @returns The total number of minutes (may be fractional)
     */
    inMinutes(): number;

    /**
     * Converts this duration to seconds.
     * @returns The total number of seconds (may be fractional)
     */
    inSeconds(): number;

    /**
     * Converts this duration to nanoseconds.
     * @returns The total number of nanoseconds
     */
    inNanoseconds(): number;

    /**
     * Converts this duration to months.
     * @returns The approximate number of months (may be fractional)
     */
    inMonths(): number;

    /**
     * Converts this duration to years.
     * @returns The approximate number of years (may be fractional)
     */
    inYears(): number;

    /**
     * Converts this duration to a JSON string.
     */
    toJSON(): string;

    // Static methods
    /**
     * Parses a duration string into a Duration object.
     * Supports formats like: "1y 2mo 3d 4h 5m 6s"
     * @param durationString The string to parse
     * @throws {Error} If the string cannot be parsed
     */
    static parse(durationString: string): Duration;

    /**
     * Creates a Duration from a JSON string.
     * @param json The JSON string to parse
     */
    static fromJSON(json: string): Duration;

    /**
     * Calculates the duration between two dates.
     * @param date1 The first date
     * @param date2 The second date
     * @returns The duration between the dates
     */
    static sub(date1: Date, date2: Date): Duration;

    /**
     * Calculates the duration until a future date.
     * @param date The future date
     * @returns The duration from now until the date
     */
    static until(date: Date): Duration;

    /**
     * Calculates the duration since a past date.
     * @param date The past date
     * @returns The duration from the date until now
     */
    static since(date: Date): Duration;

    /**
     * Calculates the duration between two dates.
     * @param date1 The first date
     * @param date2 The second date
     * @returns The duration between the dates
     */
    static between(date1: Date, date2: Date): Duration;

    // Constructor functions
    /** Creates a Duration from a number of nanoseconds */
    static fromNanoseconds(n: number): Duration;
    /** Creates a Duration from a number of microseconds */
    static fromMicroseconds(µs: number): Duration;
    /** Creates a Duration from a number of milliseconds */
    static fromMilliseconds(ms: number): Duration;
    /** Creates a Duration from a number of seconds */
    static fromSeconds(s: number): Duration;
    /** Creates a Duration from a number of minutes */
    static fromMinutes(m: number): Duration;
    /** Creates a Duration from a number of hours */
    static fromHours(h: number): Duration;
    /** Creates a Duration from a number of days */
    static fromDays(d: number): Duration;
    /** Creates a Duration from a number of weeks */
    static fromWeeks(w: number): Duration;
    /** Creates a Duration from a number of months */
    static fromMonths(m: number): Duration;
    /** Creates a Duration from a number of years */
    static fromYears(y: number): Duration;

    // Common duration constants
    /** Duration representing 1 nanosecond */
    static readonly Nanosecond: Duration;
    /** Duration representing 1 microsecond */
    static readonly Microsecond: Duration;
    /** Duration representing 1 millisecond */
    static readonly Millisecond: Duration;
    /** Duration representing 1 second */
    static readonly Second: Duration;
    /** Duration representing 1 minute */
    static readonly Minute: Duration;
    /** Duration representing 1 hour */
    static readonly Hour: Duration;
    /** Duration representing 1 day */
    static readonly Day: Duration;
    /** Duration representing 1 week */
    static readonly Week: Duration;
    /** Duration representing 1 month (approximately 30.44 days) */
    static readonly Month: Duration;
    /** Duration representing 1 year (approximately 365.25 days) */
    static readonly Year: Duration;

    // Additional common durations
    /** Duration representing 30 seconds */
    static readonly Seconds30: Duration;
    /** Duration representing 5 minutes */
    static readonly Minutes5: Duration;
    /** Duration representing 30 minutes */
    static readonly Minutes30: Duration;
    /** Duration representing 2 hours */
    static readonly Hours2: Duration;
    /** Duration representing 12 hours */
    static readonly Hours12: Duration;
    /** Duration representing 24 hours */
    static readonly Hours24: Duration;
    /** Duration representing 7 days */
    static readonly Days7: Duration;
    /** Duration representing 30 days */
    static readonly Days30: Duration;
    /** Duration representing 90 days */
    static readonly Days90: Duration;
    /** Duration representing 180 days */
    static readonly Days180: Duration;
    /** Duration representing 365 days */
    static readonly Days365: Duration;
    /** Duration representing 6 months */
    static readonly Months6: Duration;
    /** Duration representing 2 years */
    static readonly Years2: Duration;
    /** Duration representing 5 years */
    static readonly Years5: Duration;
    /** Duration representing 10 years */
    static readonly Years10: Duration;
}

/**
 * Initializes the WebAssembly module.
 * This must be called before using any Duration functionality.
 * @throws {Error} If the WASM module fails to initialize
 */
declare function initWasm(): Promise<void>;

export { Duration, initWasm };
