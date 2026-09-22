package bs

import "errors"

// Sentinel errors returned by this package. Use errors.Is to check for them;
// wrapped errors include additional context via fmt.Errorf's %w verb.
var (
	// ErrInvalidYear is returned when a BS year is outside MinBSYear..MaxBSYear.
	ErrInvalidYear = errors.New("bs: invalid year")

	// ErrInvalidMonth is returned when a BS month is not in the range 1..12.
	ErrInvalidMonth = errors.New("bs: invalid month")

	// ErrInvalidDay is returned when a BS day is not a valid day of the given
	// month (either out of the generic 1..32 bound, or greater than the
	// actual number of days in that month/year).
	ErrInvalidDay = errors.New("bs: invalid day")

	// ErrOutOfRange is returned when an AD date falls outside the Gregorian
	// range corresponding to MinBSYear..MaxBSYear.
	ErrOutOfRange = errors.New("bs: date outside supported range")

	// ErrInvalidFormat is returned when a string passed to Parse is not
	// shaped like "YYYY-MM-DD".
	ErrInvalidFormat = errors.New("bs: invalid date format")
)
