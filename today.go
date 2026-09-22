package bs

import "time"

// TodayBS returns the current system date converted to Bikram Sambat. It
// returns an error wrapping ErrOutOfRange if today's date falls outside the
// supported range (i.e. this code is running before MinBSYear or after
// MaxBSYear's corresponding Gregorian date).
func TodayBS() (Date, error) {
	return ADToBS(time.Now())
}
