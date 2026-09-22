package bs

import "fmt"

// Age computes the calendar age from birthBS to todayBS (see TodayBS), as
// years, months and days such that adding that many years, months and days
// to birthBS lands on todayBS. It returns an error wrapping ErrInvalidYear,
// ErrInvalidMonth or ErrInvalidDay if either date is invalid, or
// ErrInvalidDateOrder if birthBS is after todayBS.
func Age(birthBS, todayBS Date) (years, months, days int, err error) {
	if _, err := NewDate(birthBS.Year, birthBS.Month, birthBS.Day); err != nil {
		return 0, 0, 0, err
	}
	if _, err := NewDate(todayBS.Year, todayBS.Month, todayBS.Day); err != nil {
		return 0, 0, 0, err
	}
	if birthBS.After(todayBS) {
		return 0, 0, 0, fmt.Errorf("%w: birth date %v is after reference date %v", ErrInvalidDateOrder, birthBS, todayBS)
	}

	years = todayBS.Year - birthBS.Year
	months = todayBS.Month - birthBS.Month
	days = todayBS.Day - birthBS.Day

	if days < 0 {
		months--
		borrowMonth := todayBS.Month - 1
		borrowYear := todayBS.Year
		if borrowMonth < 1 {
			borrowMonth = monthsPerYear
			borrowYear--
		}
		daysInBorrowMonth, err := DaysInMonth(borrowYear, borrowMonth)
		if err != nil {
			return 0, 0, 0, err
		}
		days += daysInBorrowMonth
	}
	if months < 0 {
		years--
		months += monthsPerYear
	}
	return years, months, days, nil
}
