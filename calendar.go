package bs

import "time"

// FirstWeekdayOfMonth returns the day of the week that day 1 of the given
// Bikram Sambat month falls on. It returns an error wrapping ErrInvalidYear
// or ErrInvalidMonth if year or month is out of range.
func FirstWeekdayOfMonth(year, month int) (time.Weekday, error) {
	d, err := NewDate(year, month, 1)
	if err != nil {
		return 0, err
	}
	return d.DayOfWeek()
}

// WeeksInMonth returns the number of calendar-grid rows needed to display
// the given Bikram Sambat month, with weeks running Sunday through
// Saturday and the month's own days aligned under their weekday (i.e. the
// same row count MonthCalendar returns). It returns an error wrapping
// ErrInvalidYear or ErrInvalidMonth if year or month is out of range.
func WeeksInMonth(year, month int) (int, error) {
	firstWeekday, err := FirstWeekdayOfMonth(year, month)
	if err != nil {
		return 0, err
	}
	days, err := DaysInMonth(year, month)
	if err != nil {
		return 0, err
	}
	cells := int(firstWeekday) + days
	return (cells + 6) / 7, nil
}

// MonthCalendar returns a week-by-week grid of the given Bikram Sambat
// month, for building calendar UIs. Each returned week is exactly 7 cells
// (Sunday through Saturday); a nil cell means no day of this month falls in
// that slot (padding at the start of the first week and/or the end of the
// last week). Every date from day 1 to the month's last day appears exactly
// once, in order. It returns an error wrapping ErrInvalidYear or
// ErrInvalidMonth if year or month is out of range.
func MonthCalendar(year, month int) ([][]*Date, error) {
	firstWeekday, err := FirstWeekdayOfMonth(year, month)
	if err != nil {
		return nil, err
	}
	days, err := DaysInMonth(year, month)
	if err != nil {
		return nil, err
	}

	var weeks [][]*Date
	week := make([]*Date, 7)
	col := int(firstWeekday)
	for day := 1; day <= days; day++ {
		d := Date{Year: year, Month: month, Day: day}
		week[col] = &d
		col++
		if col == 7 {
			weeks = append(weeks, week)
			week = make([]*Date, 7)
			col = 0
		}
	}
	if col != 0 {
		weeks = append(weeks, week)
	}
	return weeks, nil
}
