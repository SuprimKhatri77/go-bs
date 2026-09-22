package bs

import (
	"errors"
	"testing"
)

func TestNewDateValid(t *testing.T) {
	d, err := NewDate(2083, 6, 6)
	if err != nil {
		t.Fatalf("NewDate: %v", err)
	}
	want := Date{2083, 6, 6}
	if d != want {
		t.Errorf("NewDate = %v, want %v", d, want)
	}
}

func TestDateValid(t *testing.T) {
	valid := Date{2083, 6, 6}
	if !valid.Valid() {
		t.Errorf("%v.Valid() = false, want true", valid)
	}

	invalid := Date{2083, 13, 6}
	if invalid.Valid() {
		t.Errorf("%v.Valid() = true, want false", invalid)
	}
}

func TestDaysInMonthErrors(t *testing.T) {
	if _, err := DaysInMonth(MinBSYear-1, 1); !errors.Is(err, ErrInvalidYear) {
		t.Errorf("DaysInMonth year error = %v, want wrapping ErrInvalidYear", err)
	}
	if _, err := DaysInMonth(2080, 0); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("DaysInMonth month error = %v, want wrapping ErrInvalidMonth", err)
	}
	if _, err := DaysInMonth(2080, 13); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("DaysInMonth month error = %v, want wrapping ErrInvalidMonth", err)
	}
}

func TestDateString(t *testing.T) {
	d := Date{2083, 6, 6}
	if got, want := d.String(), "2083-06-06"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestDateMonthName(t *testing.T) {
	d := Date{2083, 1, 1}
	name, err := d.MonthName()
	if err != nil {
		t.Fatalf("MonthName: %v", err)
	}
	if name != "Baisakh" {
		t.Errorf("MonthName() = %q, want %q", name, "Baisakh")
	}
}

func TestMonthName(t *testing.T) {
	name, err := MonthName(1)
	if err != nil {
		t.Fatalf("MonthName(1): %v", err)
	}
	if name != "Baisakh" {
		t.Errorf("MonthName(1) = %q, want %q", name, "Baisakh")
	}

	name, err = MonthName(12)
	if err != nil {
		t.Fatalf("MonthName(12): %v", err)
	}
	if name != "Chaitra" {
		t.Errorf("MonthName(12) = %q, want %q", name, "Chaitra")
	}

	if _, err := MonthName(0); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("MonthName(0) error = %v, want wrapping ErrInvalidMonth", err)
	}
	if _, err := MonthName(13); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("MonthName(13) error = %v, want wrapping ErrInvalidMonth", err)
	}
}
