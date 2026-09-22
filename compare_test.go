package bs

import "testing"

func TestBeforeAfterEqual(t *testing.T) {
	a := Date{2083, 6, 6}
	b := Date{2083, 6, 7}
	c := Date{2083, 6, 6}

	if !a.Before(b) {
		t.Errorf("%v.Before(%v) = false, want true", a, b)
	}
	if a.After(b) {
		t.Errorf("%v.After(%v) = true, want false", a, b)
	}
	if !b.After(a) {
		t.Errorf("%v.After(%v) = false, want true", b, a)
	}
	if b.Before(a) {
		t.Errorf("%v.Before(%v) = true, want false", b, a)
	}
	if !a.Equal(c) {
		t.Errorf("%v.Equal(%v) = false, want true", a, c)
	}
	if a.Before(c) || a.After(c) {
		t.Errorf("%v and %v are Equal but also Before/After", a, c)
	}

	// Cross-month and cross-year ordering.
	monthBoundary := Date{2083, 7, 1}
	if !b.Before(monthBoundary) {
		t.Errorf("%v.Before(%v) = false, want true", b, monthBoundary)
	}
	yearBoundary := Date{2084, 1, 1}
	if !monthBoundary.Before(yearBoundary) {
		t.Errorf("%v.Before(%v) = false, want true", monthBoundary, yearBoundary)
	}
}

func TestCompare(t *testing.T) {
	a := Date{2083, 6, 6}
	b := Date{2083, 6, 7}
	c := Date{2083, 6, 6}

	if got := Compare(a, b); got != -1 {
		t.Errorf("Compare(%v, %v) = %d, want -1", a, b, got)
	}
	if got := Compare(b, a); got != 1 {
		t.Errorf("Compare(%v, %v) = %d, want 1", b, a, got)
	}
	if got := Compare(a, c); got != 0 {
		t.Errorf("Compare(%v, %v) = %d, want 0", a, c, got)
	}
}

func TestIsSupportedBSYear(t *testing.T) {
	if !IsSupportedBSYear(MinBSYear) || !IsSupportedBSYear(MaxBSYear) {
		t.Errorf("IsSupportedBSYear should be true at both boundaries")
	}
	if IsSupportedBSYear(MinBSYear-1) || IsSupportedBSYear(MaxBSYear+1) {
		t.Errorf("IsSupportedBSYear should be false just outside the boundaries")
	}
}

func TestDaysInYear(t *testing.T) {
	got, err := DaysInYear(2083)
	if err != nil {
		t.Fatalf("DaysInYear: %v", err)
	}
	want := 0
	for m := 1; m <= monthsPerYear; m++ {
		d, err := DaysInMonth(2083, m)
		if err != nil {
			t.Fatalf("DaysInMonth: %v", err)
		}
		want += d
	}
	if got != want {
		t.Errorf("DaysInYear(2083) = %d, want %d (sum of DaysInMonth)", got, want)
	}

	if _, err := DaysInYear(MinBSYear - 1); err == nil {
		t.Errorf("DaysInYear(%d) = nil error, want error", MinBSYear-1)
	}
}
