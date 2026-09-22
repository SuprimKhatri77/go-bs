package bs

import (
	"errors"
	"testing"
)

func TestFirstWeekdayOfMonth(t *testing.T) {
	got, err := FirstWeekdayOfMonth(2083, 6)
	if err != nil {
		t.Fatalf("FirstWeekdayOfMonth: %v", err)
	}
	want, err := (Date{2083, 6, 1}).DayOfWeek()
	if err != nil {
		t.Fatalf("DayOfWeek: %v", err)
	}
	if got != want {
		t.Errorf("FirstWeekdayOfMonth(2083, 6) = %v, want %v", got, want)
	}

	if _, err := FirstWeekdayOfMonth(MinBSYear-1, 1); !errors.Is(err, ErrInvalidYear) {
		t.Errorf("FirstWeekdayOfMonth error = %v, want wrapping ErrInvalidYear", err)
	}
	if _, err := FirstWeekdayOfMonth(2080, 13); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("FirstWeekdayOfMonth error = %v, want wrapping ErrInvalidMonth", err)
	}
}

func TestWeeksInMonthMatchesMonthCalendar(t *testing.T) {
	got, err := WeeksInMonth(2083, 6)
	if err != nil {
		t.Fatalf("WeeksInMonth: %v", err)
	}
	weeks, err := MonthCalendar(2083, 6)
	if err != nil {
		t.Fatalf("MonthCalendar: %v", err)
	}
	if got != len(weeks) {
		t.Errorf("WeeksInMonth(2083, 6) = %d, want %d (len(MonthCalendar))", got, len(weeks))
	}

	if _, err := WeeksInMonth(MinBSYear-1, 1); !errors.Is(err, ErrInvalidYear) {
		t.Errorf("WeeksInMonth error = %v, want wrapping ErrInvalidYear", err)
	}
}

func TestMonthCalendarStructure(t *testing.T) {
	weeks, err := MonthCalendar(2083, 6)
	if err != nil {
		t.Fatalf("MonthCalendar: %v", err)
	}

	firstWeekday, _ := FirstWeekdayOfMonth(2083, 6)
	days, _ := DaysInMonth(2083, 6)

	var flat []*Date
	for _, week := range weeks {
		if len(week) != 7 {
			t.Fatalf("week has %d cells, want 7: %v", len(week), week)
		}
		flat = append(flat, week...)
	}

	// Leading nils, in count, equal the first weekday's index.
	for i := 0; i < int(firstWeekday); i++ {
		if flat[i] != nil {
			t.Fatalf("cell %d = %v, want nil (before day 1)", i, flat[i])
		}
	}

	// The non-nil run is exactly days 1..N in order, immediately after the
	// leading nils, with only trailing nils (if any) after it.
	for i := range days {
		cell := flat[int(firstWeekday)+i]
		if cell == nil {
			t.Fatalf("cell for day %d is nil", i+1)
		}
		want := Date{2083, 6, i + 1}
		if *cell != want {
			t.Fatalf("cell for day %d = %v, want %v", i+1, *cell, want)
		}
	}
	for i := int(firstWeekday) + days; i < len(flat); i++ {
		if flat[i] != nil {
			t.Fatalf("cell %d = %v, want nil (after last day)", i, flat[i])
		}
	}

	if _, err := MonthCalendar(2080, 0); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("MonthCalendar error = %v, want wrapping ErrInvalidMonth", err)
	}
}

// TestMonthCalendarExhaustive checks the same invariants as
// TestMonthCalendarStructure, plus the WeeksInMonth/len(MonthCalendar)
// agreement, for every supported (year, month).
func TestMonthCalendarExhaustive(t *testing.T) {
	for year := MinBSYear; year <= MaxBSYear; year++ {
		for month := 1; month <= monthsPerYear; month++ {
			weeks, err := MonthCalendar(year, month)
			if err != nil {
				t.Fatalf("MonthCalendar(%d, %d): %v", year, month, err)
			}
			weeksCount, err := WeeksInMonth(year, month)
			if err != nil {
				t.Fatalf("WeeksInMonth(%d, %d): %v", year, month, err)
			}
			if weeksCount != len(weeks) {
				t.Fatalf("WeeksInMonth(%d, %d) = %d, want %d (len(MonthCalendar))", year, month, weeksCount, len(weeks))
			}

			firstWeekday, err := FirstWeekdayOfMonth(year, month)
			if err != nil {
				t.Fatalf("FirstWeekdayOfMonth(%d, %d): %v", year, month, err)
			}
			days, err := DaysInMonth(year, month)
			if err != nil {
				t.Fatalf("DaysInMonth(%d, %d): %v", year, month, err)
			}

			var flat []*Date
			for _, week := range weeks {
				if len(week) != 7 {
					t.Fatalf("%d-%02d: week has %d cells, want 7", year, month, len(week))
				}
				flat = append(flat, week...)
			}

			nonNil := 0
			day := 1
			for i, cell := range flat {
				switch {
				case i < int(firstWeekday):
					if cell != nil {
						t.Fatalf("%d-%02d: cell %d = %v, want nil (before day 1)", year, month, i, cell)
					}
				case day <= days:
					if cell == nil {
						t.Fatalf("%d-%02d: cell %d is nil, want day %d", year, month, i, day)
					}
					want := Date{year, month, day}
					if *cell != want {
						t.Fatalf("%d-%02d: cell %d = %v, want %v", year, month, i, *cell, want)
					}
					nonNil++
					day++
				default:
					if cell != nil {
						t.Fatalf("%d-%02d: cell %d = %v, want nil (after last day)", year, month, i, cell)
					}
				}
			}
			if nonNil != days {
				t.Fatalf("%d-%02d: %d non-nil cells, want %d", year, month, nonNil, days)
			}
		}
	}
}
