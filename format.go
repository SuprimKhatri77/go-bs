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
	return d.format(layout, false)
}

// FormatNepali is like Format, but renders every token in Nepali: numbers
// in Devanagari digits, MMMM as the Nepali month name (see
// MonthNameNepali), dddd as the Nepali weekday name (see
// WeekdayNameNepali) and ddd as its short form, e.g. "बुध" for "बुधवार".
// For example, "dddd, MMMM D, YYYY" renders as "बुधवार, असोज ७, २०८३".
//
// Characters in layout that aren't part of a token, including any ASCII
// digits, are copied through unchanged. It returns the same errors as
// Format.
func (d Date) FormatNepali(layout string) (string, error) {
	return d.format(layout, true)
}

func (d Date) format(layout string, nepali bool) (string, error) {
	if _, err := NewDate(d.Year, d.Month, d.Day); err != nil {
		return "", err
	}
	weekday, err := d.DayOfWeek()
	if err != nil {
		return "", err
	}

	monthName := monthNames[d.Month-1]
	weekdayName := weekday.String()
	weekdayShort := weekdayName[:3]
	digits := func(s string) string { return s }
	if nepali {
		monthName = monthNamesNepali[d.Month-1]
		weekdayName = weekdayNamesNepali[weekday]
		weekdayShort = weekdayShortNamesNepali[weekday]
		digits = ToNepaliDigits
	}

	tokens := []formatToken{
		{"YYYY", digits(fmt.Sprintf("%04d", d.Year))},
		{"MMMM", monthName},
		{"dddd", weekdayName},
		{"ddd", weekdayShort},
		{"DD", digits(fmt.Sprintf("%02d", d.Day))},
		{"MM", digits(fmt.Sprintf("%02d", d.Month))},
		{"YY", digits(fmt.Sprintf("%02d", d.Year%100))},
		{"D", digits(strconv.Itoa(d.Day))},
		{"M", digits(strconv.Itoa(d.Month))},
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
