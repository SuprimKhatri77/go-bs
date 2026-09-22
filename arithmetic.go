package bs

import "time"

// AddDays returns the date n calendar days after d (or before, if n is
// negative). It returns an error wrapping ErrInvalidYear, ErrInvalidMonth or
// ErrInvalidDay if d is not itself a valid date, or ErrOutOfRange if the
// result falls outside the supported range.
func (d Date) AddDays(n int) (Date, error) {
	ad, err := BSToAD(d)
	if err != nil {
		return Date{}, err
	}
	return ADToBS(ad.AddDate(0, 0, n))
}

// DayOfWeek returns the day of the week d falls on, derived through its
// equivalent Gregorian date. It returns an error wrapping ErrInvalidYear,
// ErrInvalidMonth or ErrInvalidDay if d is not a valid date.
func (d Date) DayOfWeek() (time.Weekday, error) {
	ad, err := BSToAD(d)
	if err != nil {
		return 0, err
	}
	return ad.Weekday(), nil
}

// SubDays returns the date n calendar days before d (or after, if n is
// negative). It returns the same errors as AddDays.
func (d Date) SubDays(n int) (Date, error) {
	return d.AddDays(-n)
}

// DaysBetween returns the number of calendar days from a to b: positive if b
// is after a, negative if b is before a, zero if they're equal. It returns
// an error wrapping ErrInvalidYear, ErrInvalidMonth or ErrInvalidDay if
// either date is invalid.
func DaysBetween(a, b Date) (int, error) {
	adA, err := BSToAD(a)
	if err != nil {
		return 0, err
	}
	adB, err := BSToAD(b)
	if err != nil {
		return 0, err
	}
	return int(adB.Sub(adA) / (24 * time.Hour)), nil
}

// StartOfMonth returns the first day (day 1) of d's month. It returns an
// error wrapping ErrInvalidYear or ErrInvalidMonth if d's year or month is
// invalid.
func (d Date) StartOfMonth() (Date, error) {
	return NewDate(d.Year, d.Month, 1)
}

// EndOfMonth returns the last day of d's month. It returns an error wrapping
// ErrInvalidYear or ErrInvalidMonth if d's year or month is invalid.
func (d Date) EndOfMonth() (Date, error) {
	days, err := DaysInMonth(d.Year, d.Month)
	if err != nil {
		return Date{}, err
	}
	return NewDate(d.Year, d.Month, days)
}

// StartOfYear returns Baisakh 1 of d's year. It returns an error wrapping
// ErrInvalidYear if d's year is invalid.
func (d Date) StartOfYear() (Date, error) {
	return NewDate(d.Year, 1, 1)
}

// EndOfYear returns the last day of Chaitra of d's year. It returns an error
// wrapping ErrInvalidYear if d's year is invalid.
func (d Date) EndOfYear() (Date, error) {
	days, err := DaysInMonth(d.Year, monthsPerYear)
	if err != nil {
		return Date{}, err
	}
	return NewDate(d.Year, monthsPerYear, days)
}

// DayOfYear returns d's 1-based ordinal day within its year (1 for Baisakh
// 1, up to 365 or 366 for the last day of Chaitra). It returns an error
// wrapping ErrInvalidYear, ErrInvalidMonth or ErrInvalidDay if d is not a
// valid date.
func (d Date) DayOfYear() (int, error) {
	if _, err := NewDate(d.Year, d.Month, d.Day); err != nil {
		return 0, err
	}
	row := calendarData[d.Year-MinBSYear]
	day := d.Day
	for month := 0; month < d.Month-1; month++ {
		day += int(row[month])
	}
	return day, nil
}
