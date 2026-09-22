package bs

import (
	"errors"
	"testing"
)

func TestParseValid(t *testing.T) {
	d, err := Parse("2083-06-06")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	want := Date{2083, 6, 6}
	if d != want {
		t.Errorf("Parse = %v, want %v", d, want)
	}
}

func TestParseRoundTripsWithString(t *testing.T) {
	original := Date{2083, 6, 6}
	d, err := Parse(original.String())
	if err != nil {
		t.Fatalf("Parse(%q): %v", original.String(), err)
	}
	if d != original {
		t.Errorf("Parse(String()) = %v, want %v", d, original)
	}
}

func TestParseInvalid(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"empty string", "", ErrInvalidFormat},
		{"wrong separators", "2083/06/06", ErrInvalidFormat},
		{"missing zero padding", "2083-6-6", ErrInvalidFormat},
		{"trailing garbage", "2083-06-06x", ErrInvalidFormat},
		{"month out of range", "2083-13-01", ErrInvalidMonth},
		{"year out of range", "1978-01-01", ErrInvalidYear},
		{"day out of range", "2083-06-32", ErrInvalidDay},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.input); !errors.Is(err, tc.wantErr) {
				t.Errorf("Parse(%q) error = %v, want wrapping %v", tc.input, err, tc.wantErr)
			}
		})
	}
}
