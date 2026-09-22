package bs

import (
	"errors"
	"testing"
	"time"
)

// knownPairs are independently verified AD/BS date pairs spanning the
// supported range. Sources and verification method are documented in
// docs/calendar-data.md.
var knownPairs = []struct {
	name string
	bs   Date
	ad   string // YYYY-MM-DD, UTC
}{
	{"min supported date (reference)", Date{1979, 1, 1}, "1922-04-13"},
	{"last day of first supported year", Date{1979, 12, 30}, "1923-04-12"},
	{"BS 2000 new year", Date{2000, 1, 1}, "1943-04-14"},
	{"BS 2010 new year", Date{2010, 1, 1}, "1953-04-13"},
	{"BS 2062 new year", Date{2062, 1, 1}, "2005-04-14"},
	{"BS 2081 new year", Date{2081, 1, 1}, "2024-04-13"},
	{"BS 2082 new year", Date{2082, 1, 1}, "2025-04-14"},
	{"BS 2083 mid-year (today's real date at time of writing)", Date{2083, 6, 6}, "2026-09-22"},
	{"max supported date", Date{2100, 12, 31}, "2044-04-13"},
}

func mustParseAD(t *testing.T, s string) time.Time {
	t.Helper()
	tm, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("bad fixture date %q: %v", s, err)
	}
	return tm.UTC()
}

func TestKnownPairsBSToAD(t *testing.T) {
	for _, tc := range knownPairs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BSToAD(tc.bs)
			if err != nil {
				t.Fatalf("BSToAD(%v) returned error: %v", tc.bs, err)
			}
			want := mustParseAD(t, tc.ad)
			if !got.Equal(want) {
				t.Errorf("BSToAD(%v) = %s, want %s", tc.bs, got.Format("2006-01-02"), tc.ad)
			}
		})
	}
}

func TestKnownPairsADToBS(t *testing.T) {
	for _, tc := range knownPairs {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ADToBS(mustParseAD(t, tc.ad))
			if err != nil {
				t.Fatalf("ADToBS(%s) returned error: %v", tc.ad, err)
			}
			if got != tc.bs {
				t.Errorf("ADToBS(%s) = %v, want %v", tc.ad, got, tc.bs)
			}
		})
	}
}

// TestExhaustiveBSRoundTrip converts every single supported BS date to AD and
// back, and checks the result matches the original. This is preferred over a
// sample because the supported range (44,562 days) is small enough to cover
// completely.
func TestExhaustiveBSRoundTrip(t *testing.T) {
	count := 0
	for year := MinBSYear; year <= MaxBSYear; year++ {
		for month := 1; month <= monthsPerYear; month++ {
			days, err := DaysInMonth(year, month)
			if err != nil {
				t.Fatalf("DaysInMonth(%d, %d): %v", year, month, err)
			}
			for day := 1; day <= days; day++ {
				original := Date{Year: year, Month: month, Day: day}
				ad, err := BSToAD(original)
				if err != nil {
					t.Fatalf("BSToAD(%v): %v", original, err)
				}
				converted, err := ADToBS(ad)
				if err != nil {
					t.Fatalf("ADToBS(%s) for original %v: %v", ad.Format("2006-01-02"), original, err)
				}
				if converted != original {
					t.Fatalf("round trip mismatch: %v -> %s -> %v", original, ad.Format("2006-01-02"), converted)
				}
				count++
			}
		}
	}
	t.Logf("exhaustively round-tripped %d BS dates", count)
}

// TestExhaustiveADRoundTrip converts every Gregorian date in the supported AD
// range to BS and back, and checks the calendar date (year/month/day)
// matches the original.
func TestExhaustiveADRoundTrip(t *testing.T) {
	minAD, err := BSToAD(Date{MinBSYear, 1, 1})
	if err != nil {
		t.Fatalf("BSToAD(min): %v", err)
	}
	maxDay, err := DaysInMonth(MaxBSYear, monthsPerYear)
	if err != nil {
		t.Fatalf("DaysInMonth(max): %v", err)
	}
	maxAD, err := BSToAD(Date{MaxBSYear, monthsPerYear, maxDay})
	if err != nil {
		t.Fatalf("BSToAD(max): %v", err)
	}

	count := 0
	for d := minAD; !d.After(maxAD); d = d.AddDate(0, 0, 1) {
		bsDate, err := ADToBS(d)
		if err != nil {
			t.Fatalf("ADToBS(%s): %v", d.Format("2006-01-02"), err)
		}
		back, err := BSToAD(bsDate)
		if err != nil {
			t.Fatalf("BSToAD(%v) for original %s: %v", bsDate, d.Format("2006-01-02"), err)
		}
		if back.Year() != d.Year() || back.Month() != d.Month() || back.Day() != d.Day() {
			t.Fatalf("AD round trip mismatch: %s -> %v -> %s", d.Format("2006-01-02"), bsDate, back.Format("2006-01-02"))
		}
		count++
	}
	t.Logf("exhaustively round-tripped %d AD dates", count)
}

func TestInvalidInputs(t *testing.T) {
	cases := []struct {
		name             string
		year, month, day int
		wantErr          error
	}{
		{"year below minimum", MinBSYear - 1, 1, 1, ErrInvalidYear},
		{"year above maximum", MaxBSYear + 1, 1, 1, ErrInvalidYear},
		{"month zero", 2080, 0, 10, ErrInvalidMonth},
		{"month thirteen", 2080, 13, 10, ErrInvalidMonth},
		{"negative month", 2080, -1, 10, ErrInvalidMonth},
		{"day zero", 2080, 1, 0, ErrInvalidDay},
		{"negative day", 2080, 1, -5, ErrInvalidDay},
		{"day beyond month length", 2080, 1, 32, ErrInvalidDay},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewDate(tc.year, tc.month, tc.day)
			if err == nil {
				t.Fatalf("NewDate(%d, %d, %d) = nil error, want error", tc.year, tc.month, tc.day)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("NewDate(%d, %d, %d) error = %v, want wrapping %v", tc.year, tc.month, tc.day, err, tc.wantErr)
			}
			if IsValid(tc.year, tc.month, tc.day) {
				t.Errorf("IsValid(%d, %d, %d) = true, want false", tc.year, tc.month, tc.day)
			}
		})
	}
}

func TestDayBeyondActualMonthLength(t *testing.T) {
	// A day that is within the generic 1..32 bound but beyond this specific
	// month's actual length must still be rejected.
	days, err := DaysInMonth(2080, 8)
	if err != nil {
		t.Fatalf("DaysInMonth: %v", err)
	}
	if IsValid(2080, 8, days+1) {
		t.Errorf("IsValid(2080, 8, %d) = true, want false (month only has %d days)", days+1, days)
	}
	if _, err := NewDate(2080, 8, days+1); !errors.Is(err, ErrInvalidDay) {
		t.Errorf("NewDate error = %v, want wrapping ErrInvalidDay", err)
	}
}

func TestOutOfRangeADDates(t *testing.T) {
	minAD, _ := BSToAD(Date{MinBSYear, 1, 1})
	maxDay, _ := DaysInMonth(MaxBSYear, monthsPerYear)
	maxAD, _ := BSToAD(Date{MaxBSYear, monthsPerYear, maxDay})

	before := minAD.AddDate(0, 0, -1)
	after := maxAD.AddDate(0, 0, 1)

	if _, err := ADToBS(before); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("ADToBS(%s) error = %v, want wrapping ErrOutOfRange", before.Format("2006-01-02"), err)
	}
	if _, err := ADToBS(after); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("ADToBS(%s) error = %v, want wrapping ErrOutOfRange", after.Format("2006-01-02"), err)
	}
}

func TestBoundaryTransitions(t *testing.T) {
	// Last day of Chaitra (month 12) rolls into Baisakh 1 of the next year.
	lastDay, _ := DaysInMonth(2079, 12)
	lastOfYear := Date{2079, 12, lastDay}
	ad, err := BSToAD(lastOfYear)
	if err != nil {
		t.Fatalf("BSToAD(%v): %v", lastOfYear, err)
	}
	nextDayBS, err := ADToBS(ad.AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("ADToBS: %v", err)
	}
	want := Date{2080, 1, 1}
	if nextDayBS != want {
		t.Errorf("day after %v = %v, want %v", lastOfYear, nextDayBS, want)
	}

	// Last day of a month rolls into day 1 of the next month.
	daysInMonth6, _ := DaysInMonth(2080, 6)
	lastOfMonth := Date{2080, 6, daysInMonth6}
	ad2, err := BSToAD(lastOfMonth)
	if err != nil {
		t.Fatalf("BSToAD(%v): %v", lastOfMonth, err)
	}
	nextDayBS2, err := ADToBS(ad2.AddDate(0, 0, 1))
	if err != nil {
		t.Fatalf("ADToBS: %v", err)
	}
	want2 := Date{2080, 7, 1}
	if nextDayBS2 != want2 {
		t.Errorf("day after %v = %v, want %v", lastOfMonth, nextDayBS2, want2)
	}
}

func TestTimezoneIndependence(t *testing.T) {
	locs := []string{"UTC", "Asia/Kathmandu", "America/New_York"}
	for _, name := range locs {
		t.Run(name, func(t *testing.T) {
			loc, err := time.LoadLocation(name)
			if err != nil {
				t.Skipf("timezone data for %s not available: %v", name, err)
			}
			// Same calendar date (2026-09-22), different location and
			// time-of-day; ADToBS must ignore both and see the same BS date.
			t1 := time.Date(2026, 9, 22, 0, 0, 0, 0, loc)
			t2 := time.Date(2026, 9, 22, 23, 59, 59, 0, loc)
			d1, err := ADToBS(t1)
			if err != nil {
				t.Fatalf("ADToBS: %v", err)
			}
			d2, err := ADToBS(t2)
			if err != nil {
				t.Fatalf("ADToBS: %v", err)
			}
			want := Date{2083, 6, 6}
			if d1 != want || d2 != want {
				t.Errorf("ADToBS in %s = %v / %v, want %v for both", name, d1, d2, want)
			}
		})
	}
}
