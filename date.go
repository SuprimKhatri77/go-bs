// Package bs converts calendar dates between the Gregorian (AD) calendar and
// the Bikram Sambat (BS) calendar used in Nepal, for BS years MinBSYear
// through MaxBSYear inclusive.
//
// BS month lengths are not derived from a formula: they follow Nepal's
// officially published calendar and are stored as a verified, table-driven
// dataset (see docs/calendar-data.md for sources). The package has no
// runtime dependencies.
//
// Conversion operates on calendar dates (year, month, day), not timestamps.
// ADToBS considers only the Year, Month and Day components of the given
// time.Time and ignores time-of-day and location.
package bs

import (
	"fmt"
	"time"
)

// MinBSYear and MaxBSYear are the inclusive bounds of the BS years supported
// by this package.
const (
	MinBSYear = 1979
	MaxBSYear = 2100
)

// monthsPerYear is the number of months in a Bikram Sambat year.
const monthsPerYear = 12

// Month names, in order, 1-based (Month 1 is Baisakh).
var monthNames = [monthsPerYear]string{
	"Baisakh", "Jestha", "Ashadh", "Shrawan", "Bhadra", "Ashwin",
	"Kartik", "Mangsir", "Poush", "Magh", "Falgun", "Chaitra",
}

// referenceBS and referenceAD are a verified, corresponding pair of dates
// used as the anchor for all conversions: BS 1979-01-01 is AD 1922-04-13.
// See docs/calendar-data.md for how this was verified.
var (
	referenceBS = Date{Year: MinBSYear, Month: 1, Day: 1}
	referenceAD = time.Date(1922, time.April, 13, 0, 0, 0, 0, time.UTC)
)

// Date represents a Bikram Sambat calendar date. Month is 1-based, where 1 is
// Baisakh and 12 is Chaitra. A Date is not guaranteed to be valid unless it
// was constructed with NewDate or returned by this package.
type Date struct {
	Year  int
	Month int
	Day   int
}

// NewDate constructs a Date and validates it. It returns an error wrapping
// ErrInvalidYear, ErrInvalidMonth or ErrInvalidDay if the given components do
// not form a real Bikram Sambat calendar date in the supported range.
func NewDate(year, month, day int) (Date, error) {
	if year < MinBSYear || year > MaxBSYear {
		return Date{}, fmt.Errorf("%w: %d (supported range %d-%d)", ErrInvalidYear, year, MinBSYear, MaxBSYear)
	}
	if month < 1 || month > monthsPerYear {
		return Date{}, fmt.Errorf("%w: %d", ErrInvalidMonth, month)
	}
	maxDay := int(calendarData[year-MinBSYear][month-1])
	if day < 1 || day > maxDay {
		return Date{}, fmt.Errorf("%w: %d (month %d of year %d has %d days)", ErrInvalidDay, day, month, year, maxDay)
	}
	return Date{Year: year, Month: month, Day: day}, nil
}

// Valid reports whether d is a real Bikram Sambat calendar date in the
// supported range.
func (d Date) Valid() bool {
	return IsValid(d.Year, d.Month, d.Day)
}

// String returns d formatted as "YYYY-MM-DD".
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

// MonthName returns the English name of the given 1-based Bikram Sambat
// month (1 is Baisakh, 12 is Chaitra).
func MonthName(month int) (string, error) {
	if month < 1 || month > monthsPerYear {
		return "", fmt.Errorf("%w: %d", ErrInvalidMonth, month)
	}
	return monthNames[month-1], nil
}

// MonthName returns the English name of d's month (e.g. "Baisakh").
func (d Date) MonthName() (string, error) {
	return MonthName(d.Month)
}
