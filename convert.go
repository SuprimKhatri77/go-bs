package bs

import (
	"fmt"
	"time"
)

// daysInYear returns the total number of days in the given BS year. The
// caller must ensure year is within MinBSYear..MaxBSYear.
func daysInYear(year int) int {
	row := calendarData[year-MinBSYear]
	total := 0
	for _, days := range row {
		total += int(days)
	}
	return total
}

// normalizeADDate strips the time-of-day and location from t, keeping only
// its calendar (Year, Month, Day) components, fixed at UTC midnight. This
// keeps conversions calendar-safe and independent of the machine's local
// timezone or DST rules.
func normalizeADDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// bsOffset returns the number of days between referenceBS and d. d must
// already be a valid date (see NewDate/IsValid).
func bsOffset(d Date) int {
	offset := 0
	for year := MinBSYear; year < d.Year; year++ {
		offset += daysInYear(year)
	}
	row := calendarData[d.Year-MinBSYear]
	for month := 0; month < d.Month-1; month++ {
		offset += int(row[month])
	}
	offset += d.Day - 1
	return offset
}

// BSToAD converts a Bikram Sambat date to the corresponding Gregorian
// calendar date, returned as a time.Time at UTC midnight. It returns an
// error wrapping ErrInvalidYear, ErrInvalidMonth or ErrInvalidDay if d is not
// a real, supported Bikram Sambat date.
func BSToAD(d Date) (time.Time, error) {
	if _, err := NewDate(d.Year, d.Month, d.Day); err != nil {
		return time.Time{}, err
	}
	return referenceAD.AddDate(0, 0, bsOffset(d)), nil
}

// ADToBS converts a Gregorian calendar date to Bikram Sambat. Only the Year,
// Month and Day components of t are used; time-of-day and location are
// ignored. It returns an error wrapping ErrOutOfRange if t falls outside the
// Gregorian range corresponding to MinBSYear..MaxBSYear.
func ADToBS(t time.Time) (Date, error) {
	norm := normalizeADDate(t)
	remaining := int(norm.Sub(referenceAD) / (24 * time.Hour))
	if remaining < 0 {
		return Date{}, fmt.Errorf("%w: %s is before the minimum supported date", ErrOutOfRange, norm.Format("2006-01-02"))
	}

	for year := MinBSYear; year <= MaxBSYear; year++ {
		yearDays := daysInYear(year)
		if remaining < yearDays {
			row := calendarData[year-MinBSYear]
			for month := 1; month <= monthsPerYear; month++ {
				monthDays := int(row[month-1])
				if remaining < monthDays {
					return Date{Year: year, Month: month, Day: remaining + 1}, nil
				}
				remaining -= monthDays
			}
		}
		remaining -= yearDays
	}
	return Date{}, fmt.Errorf("%w: %s is after the maximum supported date", ErrOutOfRange, norm.Format("2006-01-02"))
}
