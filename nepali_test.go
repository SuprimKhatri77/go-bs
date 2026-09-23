package bs

import (
	"errors"
	"testing"
	"time"
)

func TestMonthNameNepali(t *testing.T) {
	name, err := MonthNameNepali(1)
	if err != nil {
		t.Fatalf("MonthNameNepali(1): %v", err)
	}
	if name != "वैशाख" {
		t.Errorf("MonthNameNepali(1) = %q, want %q", name, "वैशाख")
	}

	name, err = MonthNameNepali(12)
	if err != nil {
		t.Fatalf("MonthNameNepali(12): %v", err)
	}
	if name != "चैत" {
		t.Errorf("MonthNameNepali(12) = %q, want %q", name, "चैत")
	}

	if _, err := MonthNameNepali(0); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("MonthNameNepali(0) error = %v, want wrapping ErrInvalidMonth", err)
	}
	if _, err := MonthNameNepali(13); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("MonthNameNepali(13) error = %v, want wrapping ErrInvalidMonth", err)
	}
}

func TestDateMonthNameNepali(t *testing.T) {
	d := Date{2083, 6, 6}
	name, err := d.MonthNameNepali()
	if err != nil {
		t.Fatalf("MonthNameNepali: %v", err)
	}
	if name != "असोज" {
		t.Errorf("MonthNameNepali() = %q, want %q", name, "असोज")
	}
}

func TestMonthNamesNepaliCount(t *testing.T) {
	if len(monthNamesNepali) != monthsPerYear {
		t.Fatalf("len(monthNamesNepali) = %d, want %d", len(monthNamesNepali), monthsPerYear)
	}
	for i, name := range monthNamesNepali {
		if name == "" {
			t.Errorf("monthNamesNepali[%d] is empty", i)
		}
	}
}

func TestToNepaliDigits(t *testing.T) {
	cases := []struct{ in, want string }{
		{"2083-06-06", "२०८३-०६-०६"},
		{"0123456789", "०१२३४५६७८९"},
		{"no digits here", "no digits here"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := ToNepaliDigits(tc.in); got != tc.want {
			t.Errorf("ToNepaliDigits(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestFromNepaliDigits(t *testing.T) {
	cases := []struct{ in, want string }{
		{"२०८३-०६-०६", "2083-06-06"},
		{"०१२३४५६७८९", "0123456789"},
		{"no digits here", "no digits here"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := FromNepaliDigits(tc.in); got != tc.want {
			t.Errorf("FromNepaliDigits(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNepaliDigitsRoundTrip(t *testing.T) {
	original := "2083-06-06 has 10 digits: 0123456789"
	roundTripped := FromNepaliDigits(ToNepaliDigits(original))
	if roundTripped != original {
		t.Errorf("round trip = %q, want %q", roundTripped, original)
	}
}

func TestWeekdayNameNepali(t *testing.T) {
	want := []string{"आइतवार", "सोमवार", "मंगलवार", "बुधवार", "बिहिवार", "शुक्रवार", "शनिवार"}
	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		got, err := WeekdayNameNepali(weekday)
		if err != nil {
			t.Fatalf("WeekdayNameNepali(%v): %v", weekday, err)
		}
		if got != want[weekday] {
			t.Errorf("WeekdayNameNepali(%v) = %q, want %q", weekday, got, want[weekday])
		}
	}
	for _, weekday := range []time.Weekday{-1, 7} {
		if _, err := WeekdayNameNepali(weekday); !errors.Is(err, ErrInvalidWeekday) {
			t.Errorf("WeekdayNameNepali(%d) error = %v, want wrapping ErrInvalidWeekday", weekday, err)
		}
	}
}

func TestDateWeekdayNameNepali(t *testing.T) {
	d := Date{2083, 6, 6} // Tuesday, AD 2026-09-22
	name, err := d.WeekdayNameNepali()
	if err != nil {
		t.Fatalf("WeekdayNameNepali: %v", err)
	}
	if name != "मंगलवार" {
		t.Errorf("WeekdayNameNepali() = %q, want %q", name, "मंगलवार")
	}
	if _, err := (Date{2083, 13, 1}).WeekdayNameNepali(); !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("WeekdayNameNepali on invalid date error = %v, want wrapping ErrInvalidMonth", err)
	}
}

func TestWeekdayShortNamesNepaliAreStems(t *testing.T) {
	for i, short := range weekdayShortNamesNepali {
		if full := weekdayNamesNepali[i]; short+"वार" != full {
			t.Errorf("weekdayShortNamesNepali[%d] = %q, want %q without its वार suffix", i, short, full)
		}
	}
}
