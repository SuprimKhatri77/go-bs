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

func TestSubDays(t *testing.T) {
	d := Date{2083, 6, 6}

	got, err := d.SubDays(1)
	if err != nil {
		t.Fatalf("SubDays(1): %v", err)
	}
	want, err := d.AddDays(-1)
	if err != nil {
		t.Fatalf("AddDays(-1): %v", err)
	}
	if got != want {
		t.Errorf("SubDays(1) = %v, want %v (AddDays(-1))", got, want)
	}
}

func TestNextAndPreviousDay(t *testing.T) {
	d := Date{2083, 6, 6}

	next, err := d.NextDay()
	if err != nil {
		t.Fatalf("NextDay: %v", err)
	}
	if want := (Date{2083, 6, 7}); next != want {
		t.Errorf("NextDay() = %v, want %v", next, want)
	}

	prev, err := d.PreviousDay()
	if err != nil {
		t.Fatalf("PreviousDay: %v", err)
	}
	if want := (Date{2083, 6, 5}); prev != want {
		t.Errorf("PreviousDay() = %v, want %v", prev, want)
	}
}

func TestNextAndPreviousMonth(t *testing.T) {
	d := Date{2083, 6, 15}

	next, err := d.NextMonth()
	if err != nil {
		t.Fatalf("NextMonth: %v", err)
	}
	if want := (Date{2083, 7, 15}); next != want {
		t.Errorf("NextMonth() = %v, want %v", next, want)
	}

	prev, err := d.PreviousMonth()
	if err != nil {
		t.Fatalf("PreviousMonth: %v", err)
	}
	if want := (Date{2083, 5, 15}); prev != want {
		t.Errorf("PreviousMonth() = %v, want %v", prev, want)
	}
}

func TestNextAndPreviousMonthYearRollover(t *testing.T) {
	dec := Date{2083, monthsPerYear, 5}
	next, err := dec.NextMonth()
	if err != nil {
		t.Fatalf("NextMonth: %v", err)
	}
	if want := (Date{2084, 1, 5}); next != want {
		t.Errorf("NextMonth() across year boundary = %v, want %v", next, want)
	}

	jan := Date{2083, 1, 5}
	prev, err := jan.PreviousMonth()
	if err != nil {
		t.Fatalf("PreviousMonth: %v", err)
	}
	if want := (Date{2082, monthsPerYear, 5}); prev != want {
		t.Errorf("PreviousMonth() across year boundary = %v, want %v", prev, want)
	}
}

func TestNextMonthClampsShorterMonth(t *testing.T) {
	// Ashadh (month 3) 2083 has 32 days; Shrawan (month 4) has 31, so day 32
	// must clamp to Shrawan's last day rather than rolling into month 5.
	daysMonth3, _ := DaysInMonth(2083, 3)
	daysMonth4, _ := DaysInMonth(2083, 4)
	if daysMonth3 <= daysMonth4 {
		t.Fatalf("test fixture assumption broken: month 3 has %d days, month 4 has %d", daysMonth3, daysMonth4)
	}

	last := Date{2083, 3, daysMonth3}
	next, err := last.NextMonth()
	if err != nil {
		t.Fatalf("NextMonth: %v", err)
	}
	if want := (Date{2083, 4, daysMonth4}); next != want {
		t.Errorf("NextMonth() clamp = %v, want %v", next, want)
	}
}

func TestNextAndPreviousMonthInvalidReceiver(t *testing.T) {
	invalid := Date{2083, 13, 1}
	if _, err := invalid.NextMonth(); err == nil {
		t.Errorf("NextMonth on invalid date = nil error, want error")
	}
	if _, err := invalid.PreviousMonth(); err == nil {
		t.Errorf("PreviousMonth on invalid date = nil error, want error")
	}
}

func TestDaysBetween(t *testing.T) {
	a := Date{2083, 6, 6}
	b := Date{2083, 6, 16}

	got, err := DaysBetween(a, b)
	if err != nil {
		t.Fatalf("DaysBetween: %v", err)
	}
	if got != 10 {
		t.Errorf("DaysBetween(%v, %v) = %d, want 10", a, b, got)
	}

	got, err = DaysBetween(b, a)
	if err != nil {
		t.Fatalf("DaysBetween: %v", err)
	}
	if got != -10 {
		t.Errorf("DaysBetween(%v, %v) = %d, want -10", b, a, got)
	}

	got, err = DaysBetween(a, a)
	if err != nil {
		t.Fatalf("DaysBetween: %v", err)
	}
	if got != 0 {
		t.Errorf("DaysBetween(%v, %v) = %d, want 0", a, a, got)
	}

	if _, err := DaysBetween(Date{2083, 13, 1}, a); err == nil {
		t.Errorf("DaysBetween with invalid date = nil error, want error")
	}
}

func TestStartAndEndOfMonth(t *testing.T) {
	d := Date{2083, 6, 15}

	start, err := d.StartOfMonth()
	if err != nil {
		t.Fatalf("StartOfMonth: %v", err)
	}
	if want := (Date{2083, 6, 1}); start != want {
		t.Errorf("StartOfMonth() = %v, want %v", start, want)
	}

	end, err := d.EndOfMonth()
	if err != nil {
		t.Fatalf("EndOfMonth: %v", err)
	}
	days, _ := DaysInMonth(2083, 6)
	if want := (Date{2083, 6, days}); end != want {
		t.Errorf("EndOfMonth() = %v, want %v", end, want)
	}

	invalid := Date{2083, 13, 1}
	if _, err := invalid.StartOfMonth(); err == nil {
		t.Errorf("StartOfMonth on invalid date = nil error, want error")
	}
	if _, err := invalid.EndOfMonth(); err == nil {
		t.Errorf("EndOfMonth on invalid date = nil error, want error")
	}
}

func TestStartAndEndOfYear(t *testing.T) {
	d := Date{2083, 6, 15}

	start, err := d.StartOfYear()
	if err != nil {
		t.Fatalf("StartOfYear: %v", err)
	}
	if want := (Date{2083, 1, 1}); start != want {
		t.Errorf("StartOfYear() = %v, want %v", start, want)
	}

	end, err := d.EndOfYear()
	if err != nil {
		t.Fatalf("EndOfYear: %v", err)
	}
	days, _ := DaysInMonth(2083, monthsPerYear)
	if want := (Date{2083, monthsPerYear, days}); end != want {
		t.Errorf("EndOfYear() = %v, want %v", end, want)
	}

	invalid := Date{MinBSYear - 1, 1, 1}
	if _, err := invalid.StartOfYear(); err == nil {
		t.Errorf("StartOfYear on invalid date = nil error, want error")
	}
	if _, err := invalid.EndOfYear(); err == nil {
		t.Errorf("EndOfYear on invalid date = nil error, want error")
	}
}

func TestDayOfYear(t *testing.T) {
	first := Date{2083, 1, 1}
	got, err := first.DayOfYear()
	if err != nil {
		t.Fatalf("DayOfYear: %v", err)
	}
	if got != 1 {
		t.Errorf("DayOfYear() for Baisakh 1 = %d, want 1", got)
	}

	last, err := first.EndOfYear()
	if err != nil {
		t.Fatalf("EndOfYear: %v", err)
	}
	totalDays, err := DaysInYear(2083)
	if err != nil {
		t.Fatalf("DaysInYear: %v", err)
	}
	got, err = last.DayOfYear()
	if err != nil {
		t.Fatalf("DayOfYear: %v", err)
	}
	if got != totalDays {
		t.Errorf("DayOfYear() for last day of year = %d, want %d", got, totalDays)
	}

	if _, err := (Date{2083, 13, 1}).DayOfYear(); err == nil {
		t.Errorf("DayOfYear on invalid date = nil error, want error")
	}
}

// TestDayOfYearExhaustive cross-checks DayOfYear against DaysBetween(StartOfYear, d)+1
// for every supported BS date, which is small enough to check completely.
func TestDayOfYearExhaustive(t *testing.T) {
	for year := MinBSYear; year <= MaxBSYear; year++ {
		startOfYear := Date{year, 1, 1}
		for month := 1; month <= monthsPerYear; month++ {
			days, err := DaysInMonth(year, month)
			if err != nil {
				t.Fatalf("DaysInMonth(%d, %d): %v", year, month, err)
			}
			for day := 1; day <= days; day++ {
				d := Date{year, month, day}
				got, err := d.DayOfYear()
				if err != nil {
					t.Fatalf("DayOfYear(%v): %v", d, err)
				}
				diff, err := DaysBetween(startOfYear, d)
				if err != nil {
					t.Fatalf("DaysBetween: %v", err)
				}
				if want := diff + 1; got != want {
					t.Fatalf("DayOfYear(%v) = %d, want %d (DaysBetween+1)", d, got, want)
				}
			}
		}
	}
}
