package bs

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Value implements database/sql/driver.Valuer, so a Date can be passed
// directly as a query argument to database/sql. It encodes the same way as
// String and MarshalText ("YYYY-MM-DD").
//
// Value does not validate d and never returns an error; an invalid or
// zero-value Date still writes its (nonsensical) string form rather than
// failing silently as NULL. Use a *Date (nil for NULL) if you need nullable
// storage — Value is not called on a nil *Date.
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements database/sql.Scanner, so a Date can be read directly out
// of a database/sql row. It accepts:
//   - nil, which resets d to the zero Date
//   - string or []byte in "YYYY-MM-DD" form (as written by Value)
//   - time.Time, whose Year/Month/Day are taken literally as BS components
//     (not run through ADToBS) — this is what most drivers hand back for a
//     native DATE/DATETIME column, and taking the components as-is is what
//     makes that round-trip symmetric with Value's string encoding.
//
// Any other source type, or a string/[]byte not shaped like a date, returns
// an error.
func (d *Date) Scan(value any) error {
	switch v := value.(type) {
	case nil:
		*d = Date{}
		return nil
	case string:
		return d.UnmarshalText([]byte(v))
	case []byte:
		return d.UnmarshalText(v)
	case time.Time:
		parsed, err := NewDate(v.Year(), int(v.Month()), v.Day())
		if err != nil {
			return err
		}
		*d = parsed
		return nil
	default:
		return fmt.Errorf("bs: cannot scan %T into Date", value)
	}
}
