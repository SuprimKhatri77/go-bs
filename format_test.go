package bs

import (
	"errors"
	"testing"
)

func TestFormat(t *testing.T) {
	d := Date{2083, 6, 6} // verified elsewhere as a Tuesday (AD 2026-09-22)

	cases := []struct {
		layout string
		want   string
	}{
		{"YYYY-MM-DD", "2083-06-06"},
		{"YY-M-D", "83-6-6"},
		{"MMMM D, YYYY", "Ashwin 6, 2083"},
		{"dddd, MMMM D, YYYY", "Tuesday, Ashwin 6, 2083"},
		{"YYYY/MM/DD", "2083/06/06"},
		{"no tokens here", "no tokens here"},
		{"", ""},
	}
	for _, tc := range cases {
		t.Run(tc.layout, func(t *testing.T) {
			got, err := d.Format(tc.layout)
			if err != nil {
				t.Fatalf("Format(%q): %v", tc.layout, err)
			}
			if got != tc.want {
				t.Errorf("Format(%q) = %q, want %q", tc.layout, got, tc.want)
			}
		})
	}
}

func TestFormatSingleDigitDay(t *testing.T) {
	d := Date{2083, 1, 1}
	got, err := d.Format("D-M-YYYY")
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if want := "1-1-2083"; got != want {
		t.Errorf("Format(D-M-YYYY) = %q, want %q", got, want)
	}
	gotPadded, err := d.Format("DD-MM-YYYY")
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if want := "01-01-2083"; gotPadded != want {
		t.Errorf("Format(DD-MM-YYYY) = %q, want %q", gotPadded, want)
	}
}

func TestFormatWeekdayAbbreviation(t *testing.T) {
	d := Date{2083, 6, 6}
	got, err := d.Format("ddd")
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	full, err := d.Format("dddd")
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if want := full[:3]; got != want {
		t.Errorf("Format(ddd) = %q, want first 3 letters of %q = %q", got, full, want)
	}
}

func TestFormatRoundTripsWithParse(t *testing.T) {
	original := Date{2083, 6, 6}
	s, err := original.Format("YYYY-MM-DD")
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	back, err := Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q): %v", s, err)
	}
	if back != original {
		t.Errorf("Parse(Format(YYYY-MM-DD)) = %v, want %v", back, original)
	}
}

// TestFormatMatchesStringExhaustive checks that Format("YYYY-MM-DD") agrees
// with String() for every supported date, since the two are meant to
// produce (and Parse to accept) the same shape.
func TestFormatMatchesStringExhaustive(t *testing.T) {
	for year := MinBSYear; year <= MaxBSYear; year++ {
		for month := 1; month <= monthsPerYear; month++ {
			days, err := DaysInMonth(year, month)
			if err != nil {
				t.Fatalf("DaysInMonth(%d, %d): %v", year, month, err)
			}
			for day := 1; day <= days; day++ {
				d := Date{year, month, day}
				got, err := d.Format("YYYY-MM-DD")
				if err != nil {
					t.Fatalf("Format(%v): %v", d, err)
				}
				if want := d.String(); got != want {
					t.Fatalf("Format(YYYY-MM-DD) for %v = %q, want %q (String())", d, got, want)
				}
			}
		}
	}
}

func TestFormatInvalidDate(t *testing.T) {
	invalid := Date{2083, 13, 1}
	if _, err := invalid.Format("YYYY-MM-DD"); err == nil {
		t.Errorf("Format on invalid date = nil error, want error")
	}

	invalidYear := Date{MinBSYear - 1, 1, 1}
	if _, err := invalidYear.Format("YYYY"); !errors.Is(err, ErrInvalidYear) {
		t.Errorf("Format error = %v, want wrapping ErrInvalidYear", err)
	}
}
