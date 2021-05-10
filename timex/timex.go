package timex

import "time"

// NowFunc is a type of SetTimeNow argument
type NowFunc func() time.Time

var defaultNow NowFunc = func() time.Time {
	return time.Now()
}

var now = defaultNow

// Now returns current time
//
// Returns same as time.Now() at default
// but if set custom function via SetTimeNow, uses it.
// This is useful for testing
func Now() time.Time {
	return now()
}

// SetTimeNow sets custom function to change behavior of Now()
func SetTimeNow(f NowFunc) {
	now = f
}
