package bs

import (
	"fmt"
	"regexp"
	"strconv"
)

// dateFormat matches the "YYYY-MM-DD" shape produced by Date.String.
var dateFormat = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)

// Parse parses a Bikram Sambat date in "YYYY-MM-DD" format (the same format
// Date.String produces). It returns an error wrapping ErrInvalidFormat if s
// is not shaped like a date, or ErrInvalidYear, ErrInvalidMonth or
// ErrInvalidDay if it's shaped correctly but not a real supported date.
func Parse(s string) (Date, error) {
	m := dateFormat.FindStringSubmatch(s)
	if m == nil {
		return Date{}, fmt.Errorf("%w: %q (want YYYY-MM-DD)", ErrInvalidFormat, s)
	}
	year, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	day, _ := strconv.Atoi(m[3])
	return NewDate(year, month, day)
}
