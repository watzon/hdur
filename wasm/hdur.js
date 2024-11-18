// Initialize the WASM
let wasmInit = false;
let globalDuration = null;
let initPromise = null;

async function loadWasmExec() {
    if (typeof Go !== 'undefined') return;

    // In Node.js
    if (typeof process !== 'undefined' && process.versions && process.versions.node) {
        try {
            // Try to load from same directory first
            await import('./wasm_exec.js');
        } catch (e) {
            // If that fails, try to load from dist directory
            await import('./dist/wasm_exec.js');
        }
        return;
    }

    // In browser
    return new Promise((resolve, reject) => {
        const script = document.createElement('script');
        script.src = new URL('./dist/wasm_exec.js', import.meta.url).href;
        script.onload = () => resolve();
        script.onerror = () => reject(new Error('Failed to load wasm_exec.js'));
        document.head.appendChild(script);
    });
}

async function initWasm() {
    // Return existing initialization if it's in progress
    if (initPromise) return initPromise;
    
    // Return cached result if already initialized
    if (wasmInit) return globalDuration;
    
    initPromise = (async () => {
        try {
            await loadWasmExec();
        } catch (err) {
            throw new Error('Failed to load wasm_exec.js. Make sure the file is available in the dist directory.');
        }

        if (typeof Go === 'undefined') {
            throw new Error('wasm_exec.js was loaded but Go is not defined');
        }

        const go = new Go();
        let wasmInstance;

        // Handle Node.js environment
        if (typeof process !== 'undefined' && process.versions && process.versions.node) {
            const fs = await import('fs/promises');
            const { fileURLToPath } = await import('url');
            const { join, dirname } = await import('path');
            
            const __filename = fileURLToPath(import.meta.url);
            const __dirname = dirname(__filename);
            
            // Try current directory first, then dist
            let wasmBuffer;
            try {
                wasmBuffer = await fs.readFile(join(__dirname, 'hdur.wasm'));
            } catch (e) {
                wasmBuffer = await fs.readFile(join(__dirname, '..', 'dist', 'hdur.wasm'));
            }
            
            const result = await WebAssembly.instantiate(wasmBuffer, go.importObject);
            wasmInstance = result.instance;
        } else {
            // Browser environment
            const wasmPath = new URL('./hdur.wasm', import.meta.url);
            const result = await WebAssembly.instantiateStreaming(
                fetch(wasmPath),
                go.importObject
            );
            wasmInstance = result.instance;
        }

        go.run(wasmInstance);
        
        // Now that WASM is loaded, create the Duration class
        class Duration {
            constructor(years = 0, months = 0, days = 0, hours = 0, minutes = 0, seconds = 0, nanos = 0) {
                this.years = years;
                this.months = months;
                this.days = days;
                this.hours = hours;
                this.minutes = minutes;
                this.seconds = seconds;
                this.nanos = nanos;
            }

            toString() {
                return formatDuration(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }

            format(pattern) {
                return formatWithPattern(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    pattern
                );
            }

            add(date) {
                if (!(date instanceof Date)) {
                    throw new Error('add() requires a Date object');
                }
                const isoDate = date.toISOString();
                const result = addDuration(
                    isoDate,
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
                return new Date(result);
            }

            equals(other) {
                if (!(other instanceof Duration)) {
                    throw new Error('equals() requires a Duration object');
                }
                return equalDurations(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    other.years,
                    other.months,
                    other.days,
                    other.hours,
                    other.minutes,
                    other.seconds,
                    other.nanos
                );
            }

            lessThan(other) {
                if (!(other instanceof Duration)) {
                    throw new Error('lessThan() requires a Duration object');
                }
                return lessDurations(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    other.years,
                    other.months,
                    other.days,
                    other.hours,
                    other.minutes,
                    other.seconds,
                    other.nanos
                );
            }

            lessOrEqual(other) {
                if (!(other instanceof Duration)) {
                    throw new Error('lessOrEqual() requires a Duration object');
                }
                return lessThanOrEqual(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    other.years,
                    other.months,
                    other.days,
                    other.hours,
                    other.minutes,
                    other.seconds,
                    other.nanos
                );
            }

            greaterThan(other) {
                if (!(other instanceof Duration)) {
                    throw new Error('greaterThan() requires a Duration object');
                }
                return greaterDurations(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    other.years,
                    other.months,
                    other.days,
                    other.hours,
                    other.minutes,
                    other.seconds,
                    other.nanos
                );
            }

            greaterOrEqual(other) {
                if (!(other instanceof Duration)) {
                    throw new Error('greaterOrEqual() requires a Duration object');
                }
                return greaterThanOrEqual(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    other.years,
                    other.months,
                    other.days,
                    other.hours,
                    other.minutes,
                    other.seconds,
                    other.nanos
                );
            }

            abs() {
                const result = absDuration(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            multiply(factor) {
                if (typeof factor !== 'number') {
                    throw new Error('multiply() requires a number');
                }
                const result = mulDuration(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    factor
                );
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            divide(divisor) {
                if (typeof divisor !== 'number') {
                    throw new Error('divide() requires a number');
                }
                if (divisor === 0) {
                    throw new Error('division by zero');
                }
                const result = divDuration(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    divisor
                );
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            round(multiple) {
                if (!(multiple instanceof Duration)) {
                    throw new Error('round() requires a Duration object');
                }
                const result = roundDuration(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    multiple.years,
                    multiple.months,
                    multiple.days,
                    multiple.hours,
                    multiple.minutes,
                    multiple.seconds,
                    multiple.nanos
                );
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            truncate(multiple) {
                if (!(multiple instanceof Duration)) {
                    throw new Error('truncate() requires a Duration object');
                }
                const result = truncateDuration(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos,
                    multiple.years,
                    multiple.months,
                    multiple.days,
                    multiple.hours,
                    multiple.minutes,
                    multiple.seconds,
                    multiple.nanos
                );
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static sub(date1, date2) {
                if (!(date1 instanceof Date) || !(date2 instanceof Date)) {
                    throw new Error('sub() requires two Date objects');
                }
                const result = subDuration(date1.toISOString(), date2.toISOString());
                if (result.error) {
                    throw new Error(result.error);
                }
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            isZero() {
                return isZeroDuration(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }

            toStandard() {
                return toStandard(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }

            static parse(durationString) {
                const result = parseDuration(durationString);
                if (result.error) {
                    throw new Error(result.error);
                }
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static until(date) {
                if (!(date instanceof Date)) {
                    throw new Error('until() requires a Date object');
                }
                const result = untilTime(date.toISOString());
                if (result.error) {
                    throw new Error(result.error);
                }
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static since(date) {
                if (!(date instanceof Date)) {
                    throw new Error('since() requires a Date object');
                }
                const result = sinceTime(date.toISOString());
                if (result.error) {
                    throw new Error(result.error);
                }
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static between(date1, date2) {
                if (!(date1 instanceof Date) || !(date2 instanceof Date)) {
                    throw new Error('between() requires two Date objects');
                }
                const result = betweenTimes(date1.toISOString(), date2.toISOString());
                if (result.error) {
                    throw new Error(result.error);
                }
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            toJSON() {
                return this.toString();
            }

            static fromJSON(json) {
                return Duration.parse(json);
            }

            static fromNanoseconds(n) {
                const result = fromNanoseconds(n);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromMicroseconds(µs) {
                const result = fromMicroseconds(µs);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromMilliseconds(ms) {
                const result = fromMilliseconds(ms);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromSeconds(s) {
                const result = fromSeconds(s);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromMinutes(m) {
                const result = fromMinutes(m);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromHours(h) {
                const result = fromHours(h);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromDays(d) {
                const result = fromDays(d);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromWeeks(w) {
                const result = fromWeeks(w);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromMonths(m) {
                const result = fromMonths(m);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            static fromYears(y) {
                const result = fromYears(y);
                return new Duration(
                    result.years,
                    result.months,
                    result.days,
                    result.hours,
                    result.minutes,
                    result.seconds,
                    result.nanos
                );
            }

            inHours() {
                return inHours(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }

            inMinutes() {
                return inMinutes(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }

            inSeconds() {
                return inSeconds(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }

            inNanoseconds() {
                return inNanoseconds(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }

            inMonths() {
                return inMonths(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }

            inYears() {
                return inYears(
                    this.years,
                    this.months,
                    this.days,
                    this.hours,
                    this.minutes,
                    this.seconds,
                    this.nanos
                );
            }
        }

        // Initialize common duration constants after WASM is loaded
        Duration.Nanosecond = Duration.fromNanoseconds(1);
        Duration.Microsecond = Duration.fromMicroseconds(1);
        Duration.Millisecond = Duration.fromMilliseconds(1);
        Duration.Second = Duration.fromSeconds(1);
        Duration.Minute = Duration.fromMinutes(1);
        Duration.Hour = Duration.fromHours(1);
        Duration.Day = Duration.fromDays(1);
        Duration.Week = Duration.fromWeeks(1);
        Duration.Month = Duration.fromMonths(1);
        Duration.Year = Duration.fromYears(1);

        Duration.Seconds30 = Duration.fromSeconds(30);
        Duration.Minutes5 = Duration.fromMinutes(5);
        Duration.Minutes30 = Duration.fromMinutes(30);
        Duration.Hours2 = Duration.fromHours(2);
        Duration.Hours12 = Duration.fromHours(12);
        Duration.Hours24 = Duration.fromHours(24);
        Duration.Days7 = Duration.fromDays(7);
        Duration.Days30 = Duration.fromDays(30);
        Duration.Days90 = Duration.fromDays(90);
        Duration.Days180 = Duration.fromDays(180);
        Duration.Days365 = Duration.fromDays(365);
        Duration.Months6 = Duration.fromMonths(6);
        Duration.Years2 = Duration.fromYears(2);
        Duration.Years5 = Duration.fromYears(5);
        Duration.Years10 = Duration.fromYears(10);

        // Store the Duration class globally
        globalDuration = Duration;
        wasmInit = true;
        return Duration;
    })();

    return initPromise;
}

// Create a proxy handler that initializes WASM when needed
const handler = {
    construct: async function(target, args) {
        const RealDuration = await initWasm();
        return Reflect.construct(RealDuration, args);
    },
    get: async function(target, prop, receiver) {
        if (prop === 'then') return undefined; // Make the proxy not thenable
        const RealDuration = await initWasm();
        return Reflect.get(RealDuration, prop, receiver);
    },
    apply: async function(target, thisArg, args) {
        const RealDuration = await initWasm();
        return Reflect.apply(RealDuration, thisArg, args);
    }
};

// Create an async proxy wrapper to handle promise-based access
function createAsyncProxy(target) {
    return new Proxy(target, {
        get(target, prop) {
            if (prop === 'then') return undefined; // Make the proxy not thenable
            return async function(...args) {
                const RealDuration = await initWasm();
                return RealDuration[prop](...args);
            };
        }
    });
}

// Export a proxy that transparently handles initialization
const DurationProxy = new Proxy(function(){}, handler);
export const Duration = createAsyncProxy(DurationProxy);

// Still export initWasm for users who want to initialize explicitly
export { initWasm };
