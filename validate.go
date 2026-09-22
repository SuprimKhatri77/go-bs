package bs

import "fmt"

// DaysInMonth returns the number of days in the given Bikram Sambat month.
// It returns ErrInvalidYear or ErrInvalidMonth if year or month is out of
// range.
func DaysInMonth(year, month int) (int, error) {
	if year < MinBSYear || year > MaxBSYear {
		return 0, fmt.Errorf("%w: %d (supported range %d-%d)", ErrInvalidYear, year, MinBSYear, MaxBSYear)
	}
	if month < 1 || month > monthsPerYear {
		return 0, fmt.Errorf("%w: %d", ErrInvalidMonth, month)
	}
	return int(calendarData[year-MinBSYear][month-1]), nil
}

// DaysInYear returns the total number of days in the given Bikram Sambat
// year. It returns ErrInvalidYear if year is out of range.
func DaysInYear(year int) (int, error) {
	if year < MinBSYear || year > MaxBSYear {
		return 0, fmt.Errorf("%w: %d (supported range %d-%d)", ErrInvalidYear, year, MinBSYear, MaxBSYear)
	}
	return daysInYear(year), nil
}

// IsValid reports whether year, month and day form a real Bikram Sambat
// calendar date within the supported range (MinBSYear..MaxBSYear).
func IsValid(year, month, day int) bool {
	days, err := DaysInMonth(year, month)
	if err != nil {
		return false
	}
	return day >= 1 && day <= days
}

// IsSupportedBSYear reports whether year is within MinBSYear..MaxBSYear.
func IsSupportedBSYear(year int) bool {
	return year >= MinBSYear && year <= MaxBSYear
}
