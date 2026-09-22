package bs

// MarshalText implements encoding.TextMarshaler, rendering d the same way as
// String ("YYYY-MM-DD"). Because encoding/json prefers TextMarshaler over a
// struct's default field-by-field encoding, this is what makes a Date field
// serialize as the string "2083-06-06" instead of {"Year":2083,"Month":6,
// "Day":6} — and it also makes Date usable as a map key, and with
// encoding/gob, encoding/csv and url.Values, all for free.
//
// Like String, MarshalText does not validate d; an invalid or zero-value
// Date still marshals, matching time.Time's behavior.
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler, parsing the same
// "YYYY-MM-DD" shape Parse accepts. It returns the same errors Parse does.
func (d *Date) UnmarshalText(data []byte) error {
	parsed, err := Parse(string(data))
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}
