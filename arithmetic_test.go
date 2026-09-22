package bs

import (
	"errors"
	"testing"
	"time"
)

func TestAddDays(t *testing.T) {
	d := Date{2083, 6, 6} // 2026-09-22, a Tuesday

	next, err := d.AddDays(1)
	if err != nil {
		t.Fatalf("AddDays(1): %v", err)
	}
	if want := (Date{2083, 6, 7}); next != want {
		t.Errorf("AddDays(1) = %v, want %v", next, want)
	}

	prev, err := d.AddDays(-1)
	if err != nil {
		t.Fatalf("AddDays(-1): %v", err)
	}
	if want := (Date{2083, 6, 5}); prev != want {
		t.Errorf("AddDays(-1) = %v, want %v", prev, want)
	}

	// Crossing a month boundary.
	lastDay, _ := DaysInMonth(2083, 6)
	last := Date{2083, 6, lastDay}
	rolled, err := last.AddDays(1)
	if err != nil {
		t.Fatalf("AddDays(1): %v", err)
	}
	if want := (Date{2083, 7, 1}); rolled != want {
		t.Errorf("AddDays(1) across month boundary = %v, want %v", rolled, want)
	}

	// AddDays(0) is a no-op.
	same, err := d.AddDays(0)
	if err != nil {
		t.Fatalf("AddDays(0): %v", err)
	}
	if same != d {
		t.Errorf("AddDays(0) = %v, want %v", same, d)
	}

	// Going past the supported range errors instead of wrapping/panicking.
	maxDay, _ := DaysInMonth(MaxBSYear, monthsPerYear)
	last2100 := Date{MaxBSYear, monthsPerYear, maxDay}
	if _, err := last2100.AddDays(1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("AddDays past MaxBSYear error = %v, want wrapping ErrOutOfRange", err)
	}

	// An invalid receiver errors rather than silently producing a result.
	invalid := Date{2083, 13, 1}
	if _, err := invalid.AddDays(1); err == nil {
		t.Errorf("AddDays on invalid date = nil error, want error")
	}
}

func TestDayOfWeek(t *testing.T) {
	d := Date{2083, 6, 6} // verified as AD 2026-09-22 elsewhere in this package
	want := time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC).Weekday()

	got, err := d.DayOfWeek()
	if err != nil {
		t.Fatalf("DayOfWeek: %v", err)
	}
	if got != want {
		t.Errorf("DayOfWeek() = %v, want %v", got, want)
	}

	invalid := Date{2083, 13, 1}
	if _, err := invalid.DayOfWeek(); err == nil {
		t.Errorf("DayOfWeek on invalid date = nil error, want error")
	}
}
