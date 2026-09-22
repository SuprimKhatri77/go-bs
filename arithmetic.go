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
