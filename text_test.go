package bs

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestMarshalText(t *testing.T) {
	d := Date{2083, 6, 6}
	got, err := d.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	if want := "2083-06-06"; string(got) != want {
		t.Errorf("MarshalText = %q, want %q", got, want)
	}
}

func TestUnmarshalText(t *testing.T) {
	var d Date
	if err := d.UnmarshalText([]byte("2083-06-06")); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}
	if want := (Date{2083, 6, 6}); d != want {
		t.Errorf("UnmarshalText -> %v, want %v", d, want)
	}
}

func TestUnmarshalTextInvalid(t *testing.T) {
	var d Date
	err := d.UnmarshalText([]byte("2083-13-01"))
	if !errors.Is(err, ErrInvalidMonth) {
		t.Errorf("UnmarshalText error = %v, want wrapping ErrInvalidMonth", err)
	}
}

func TestTextRoundTrip(t *testing.T) {
	original := Date{2083, 6, 6}
	text, err := original.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	var d Date
	if err := d.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}
	if d != original {
		t.Errorf("round trip = %v, want %v", d, original)
	}
}

// encoding/json prefers TextMarshaler/TextUnmarshaler over default struct
// encoding, so a Date field should serialize as a plain "YYYY-MM-DD" string.
func TestJSONUsesTextMarshaler(t *testing.T) {
	type payload struct {
		Date Date `json:"date"`
	}
	b, err := json.Marshal(payload{Date: Date{2083, 6, 6}})
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if want := `{"date":"2083-06-06"}`; string(b) != want {
		t.Errorf("json.Marshal = %s, want %s", b, want)
	}

	var got payload
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if want := (Date{2083, 6, 6}); got.Date != want {
		t.Errorf("json.Unmarshal -> %v, want %v", got.Date, want)
	}
}
