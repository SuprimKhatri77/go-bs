package bs

import (
	"fmt"
	"strconv"
	"strings"
)

// formatToken pairs a layout token with the value it expands to. Order
// matters: longer tokens must come before any shorter token that is a
// prefix of it (e.g. "YYYY" before "YY"), so Format can match greedily.
type formatToken struct {
	token string
	value string
}

// Format renders d according to layout, replacing recognized tokens:
//
//	YYYY  4-digit year, e.g. "2083"
//	YY    2-digit year, e.g. "83"
//	MMMM  full month name, e.g. "Ashwin"
//	MM    2-digit month, zero-padded
//	M     month, no leading zero
//	DD    2-digit day, zero-padded
//	D     day, no leading zero
//	dddd  full weekday name, e.g. "Thursday"
//	ddd   3-letter weekday abbreviation, e.g. "Thu"
//
// Any other character in layout (including punctuation and spaces) is
// copied through unchanged. It returns an error wrapping ErrInvalidYear,
// ErrInvalidMonth or ErrInvalidDay if d is not a valid date.
func (d Date) Format(layout string) (string, error) {
	if _, err := NewDate(d.Year, d.Month, d.Day); err != nil {
		return "", err
	}
	weekday, err := d.DayOfWeek()
	if err != nil {
		return "", err
	}
	weekdayName := weekday.String()

	tokens := []formatToken{
		{"YYYY", fmt.Sprintf("%04d", d.Year)},
		{"MMMM", monthNames[d.Month-1]},
		{"dddd", weekdayName},
		{"ddd", weekdayName[:3]},
		{"DD", fmt.Sprintf("%02d", d.Day)},
		{"MM", fmt.Sprintf("%02d", d.Month)},
		{"YY", fmt.Sprintf("%02d", d.Year%100)},
		{"D", strconv.Itoa(d.Day)},
		{"M", strconv.Itoa(d.Month)},
	}

	var b strings.Builder
	for i := 0; i < len(layout); {
		matched := false
		for _, t := range tokens {
			if strings.HasPrefix(layout[i:], t.token) {
				b.WriteString(t.value)
				i += len(t.token)
				matched = true
				break
			}
		}
		if !matched {
			b.WriteByte(layout[i])
			i++
		}
	}
	return b.String(), nil
}
