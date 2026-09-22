package bs

import (
	"errors"
	"testing"
)

func TestAgeSameDate(t *testing.T) {
	d := Date{2083, 6, 6}
	years, months, days, err := Age(d, d)
	if err != nil {
		t.Fatalf("Age: %v", err)
	}
	if years != 0 || months != 0 || days != 0 {
		t.Errorf("Age(d, d) = (%d, %d, %d), want (0, 0, 0)", years, months, days)
	}
}

func TestAgeExactYears(t *testing.T) {
	birth := Date{2060, 6, 15}
	today := Date{2083, 6, 15}
	years, months, days, err := Age(birth, today)
	if err != nil {
		t.Fatalf("Age: %v", err)
	}
	if years != 23 || months != 0 || days != 0 {
		t.Errorf("Age = (%d, %d, %d), want (23, 0, 0)", years, months, days)
	}
}

func TestAgeWithDayBorrow(t *testing.T) {
	// today's day-of-month (10) is before birth's (15), so the day count
	// must borrow from the month before today's month.
	birth := Date{2060, 6, 15}
	today := Date{2083, 6, 10}
	years, months, days, err := Age(birth, today)
	if err != nil {
		t.Fatalf("Age: %v", err)
	}
	daysInBorrowMonth, _ := DaysInMonth(2083, 5)
	wantDays := daysInBorrowMonth - 5 // 5 days short of a full 23rd year
	if years != 22 || months != 11 || days != wantDays {
		t.Errorf("Age = (%d, %d, %d), want (22, 11, %d)", years, months, days, wantDays)
	}
}

func TestAgeOneMonthNoDayBorrow(t *testing.T) {
	birth := Date{MinBSYear, 1, 1}
	today := Date{MinBSYear, 2, 1}
	years, months, days, err := Age(birth, today)
	if err != nil {
		t.Fatalf("Age: %v", err)
	}
	if years != 0 || months != 1 || days != 0 {
		t.Errorf("Age = (%d, %d, %d), want (0, 1, 0)", years, months, days)
	}
}

func TestAgeBirthAfterToday(t *testing.T) {
	birth := Date{2083, 6, 10}
	today := Date{2083, 6, 5}
	if _, _, _, err := Age(birth, today); !errors.Is(err, ErrInvalidDateOrder) {
		t.Errorf("Age error = %v, want wrapping ErrInvalidDateOrder", err)
	}
}

func TestAgeInvalidDates(t *testing.T) {
	valid := Date{2083, 6, 6}
	invalid := Date{2083, 13, 1}

	if _, _, _, err := Age(invalid, valid); err == nil {
		t.Errorf("Age with invalid birth date = nil error, want error")
	}
	if _, _, _, err := Age(valid, invalid); err == nil {
		t.Errorf("Age with invalid today date = nil error, want error")
	}
}
