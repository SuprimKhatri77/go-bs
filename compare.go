package bs

// Compare compares two dates and returns -1 if a is before b, 0 if they're
// equal, and 1 if a is after b. It compares fields directly and does not
// require either date to be valid.
func Compare(a, b Date) int {
	switch {
	case dateLess(a, b):
		return -1
	case dateLess(b, a):
		return 1
	default:
		return 0
	}
}

// Before reports whether d is chronologically before other. It compares
// fields directly and does not require either date to be valid.
func (d Date) Before(other Date) bool {
	return dateLess(d, other)
}

// After reports whether d is chronologically after other. It compares
// fields directly and does not require either date to be valid.
func (d Date) After(other Date) bool {
	return dateLess(other, d)
}

// Equal reports whether d and other represent the same calendar date.
func (d Date) Equal(other Date) bool {
	return d == other
}

func dateLess(a, b Date) bool {
	if a.Year != b.Year {
		return a.Year < b.Year
	}
	if a.Month != b.Month {
		return a.Month < b.Month
	}
	return a.Day < b.Day
}
