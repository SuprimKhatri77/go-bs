package bs

import (
	"testing"
	"time"
)

// FuzzBSRoundTrip checks that no BS input (valid or invalid) causes a panic,
// and that every input NewDate accepts round-trips through BSToAD/ADToBS
// unchanged.
func FuzzBSRoundTrip(f *testing.F) {
	f.Add(1979, 1, 1)
	f.Add(2100, 12, 31)
	f.Add(2083, 6, 6)
	f.Add(0, 0, 0)
	f.Add(-1, 13, 40)

	f.Fuzz(func(t *testing.T, year, month, day int) {
		d, err := NewDate(year, month, day)
		if err != nil {
			// Invalid input must error, never panic; nothing else to check.
			return
		}
		ad, err := BSToAD(d)
		if err != nil {
			t.Fatalf("BSToAD(%v) errored after successful NewDate: %v", d, err)
		}
		back, err := ADToBS(ad)
		if err != nil {
			t.Fatalf("ADToBS(%s) errored: %v", ad.Format("2006-01-02"), err)
		}
		if back != d {
			t.Fatalf("round trip mismatch: %v -> %s -> %v", d, ad.Format("2006-01-02"), back)
		}
	})
}

// FuzzADRoundTrip checks that no AD input (in or out of the supported range)
// causes a panic, and that every in-range input round-trips through
// ADToBS/BSToAD unchanged.
func FuzzADRoundTrip(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(1 << 40))
	f.Add(int64(-1 << 40))

	minAD, _ := BSToAD(Date{MinBSYear, 1, 1})

	f.Fuzz(func(t *testing.T, offsetSeconds int64) {
		tm := minAD.Add(time.Duration(offsetSeconds) * time.Second)
		d, err := ADToBS(tm)
		if err != nil {
			// Out-of-range input must error, never panic.
			return
		}
		back, err := BSToAD(d)
		if err != nil {
			t.Fatalf("BSToAD(%v) errored after successful ADToBS: %v", d, err)
		}
		norm := normalizeADDate(tm)
		if !back.Equal(norm) {
			t.Fatalf("round trip mismatch: %s -> %v -> %s", norm.Format("2006-01-02"), d, back.Format("2006-01-02"))
		}
	})
}
