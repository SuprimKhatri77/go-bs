package bs

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
