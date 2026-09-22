package bs

import "time"

// nepalTime is Nepal Standard Time's fixed UTC+05:45 offset, in effect
// since 1986 with no daylight saving ever observed. A time.FixedZone is
// used instead of time.LoadLocation("Asia/Kathmandu") so TodayBS behaves
// identically everywhere, including minimal containers without a system
// timezone database.
var nepalTime = time.FixedZone("NPT", int((5*time.Hour + 45*time.Minute).Seconds()))

// TodayBS returns the current date in Nepal (Nepal Standard Time,
// UTC+05:45) converted to Bikram Sambat.
//
// It deliberately does not use the calling process's local timezone: on a
// server configured for UTC (a common default for cloud VMs and
// containers), "today" would otherwise be wrong for roughly 5h45m of every
// day — the window after midnight has passed in Nepal but before it's
// passed in UTC.
//
// It returns an error wrapping ErrOutOfRange if today's date in Nepal falls
// outside the supported range (i.e. this code is running before MinBSYear
// or after MaxBSYear's corresponding Gregorian date).
func TodayBS() (Date, error) {
	return ADToBS(time.Now().In(nepalTime))
}
