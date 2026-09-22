package bs

import (
	"errors"
	"testing"
	"time"
)

func TestValue(t *testing.T) {
	d := Date{2083, 6, 6}
	got, err := d.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if want := "2083-06-06"; got != want {
		t.Errorf("Value = %v, want %q", got, want)
	}
}

func TestScanString(t *testing.T) {
	var d Date
	if err := d.Scan("2083-06-06"); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if want := (Date{2083, 6, 6}); d != want {
		t.Errorf("Scan -> %v, want %v", d, want)
	}
}

func TestScanBytes(t *testing.T) {
	var d Date
	if err := d.Scan([]byte("2083-06-06")); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if want := (Date{2083, 6, 6}); d != want {
		t.Errorf("Scan -> %v, want %v", d, want)
	}
}

func TestScanTime(t *testing.T) {
	var d Date
	in := time.Date(2083, time.June, 6, 15, 30, 0, 0, time.UTC)
	if err := d.Scan(in); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if want := (Date{2083, 6, 6}); d != want {
		t.Errorf("Scan(time.Time) -> %v, want %v (time-of-day must be ignored)", d, want)
	}
}

func TestScanNil(t *testing.T) {
	d := Date{2083, 6, 6}
	if err := d.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if want := (Date{}); d != want {
		t.Errorf("Scan(nil) -> %v, want zero value %v", d, want)
	}
}

func TestScanInvalidString(t *testing.T) {
	var d Date
	err := d.Scan("2083-13-01")
	if !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("Scan error = %v, want wrapping ErrInvalidMonth", err)
	}
}

func TestScanInvalidTime(t *testing.T) {
	var d Date
	// A real Gregorian date, but 500 isn't a supported BS year.
	err := d.Scan(time.Date(500, time.June, 6, 0, 0, 0, 0, time.UTC))
	if !errors.Is(err, ErrInvalidYear) {
		t.Errorf("Scan error = %v, want wrapping ErrInvalidYear", err)
	}
}

func TestScanUnsupportedType(t *testing.T) {
	var d Date
	err := d.Scan(42)
	if err == nil {
		t.Error("Scan(int) should error")
	}
}

func TestSQLRoundTrip(t *testing.T) {
	original := Date{2083, 6, 6}
	v, err := original.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	var d Date
	if err := d.Scan(v); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if d != original {
		t.Errorf("round trip = %v, want %v", d, original)
	}
}
