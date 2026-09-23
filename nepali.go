package bs

import (
	"fmt"
	"strings"
	"time"
)

// monthNamesNepali holds the Nepali (Devanagari) names of the Bikram Sambat
// months, in order, 1-based. Verified against Hamro Patro's own calendar
// page titles (the same source data.go's calendar data comes from), e.g.
// hamropatro.com/calendar/2100/1 titles its page "वैशाख २१००".
var monthNamesNepali = [monthsPerYear]string{
	"वैशाख", "जेठ", "असार", "साउन", "भदौ", "असोज",
	"कार्तिक", "मंसिर", "पुष", "माघ", "फागुन", "चैत",
}

// weekdayNamesNepali holds the Nepali (Devanagari) weekday names, indexed by
// time.Weekday (Sunday first). Spelled as Hamro Patro's calendar spells them
// (its weekday column headers, e.g. "आइतवार"), the same source as
// monthNamesNepali. Nepal Academy's standard spelling uses -बार instead
// (e.g. "आइतबार"); both are in common use.
var weekdayNamesNepali = [7]string{
	"आइतवार", "सोमवार", "मंगलवार", "बुधवार", "बिहिवार", "शुक्रवार", "शनिवार",
}

// weekdayShortNamesNepali holds the short Nepali weekday names used by
// FormatNepali's ddd token: each full name in weekdayNamesNepali without its
// "वार" suffix, as printed Nepali calendars usually head their columns.
var weekdayShortNamesNepali = [7]string{
	"आइत", "सोम", "मंगल", "बुध", "बिहि", "शुक्र", "शनि",
}

// nepaliDigits maps ASCII digits 0-9 to their Devanagari equivalents.
var nepaliDigits = [10]rune{'०', '१', '२', '३', '४', '५', '६', '७', '८', '९'}

// MonthNameNepali returns the Devanagari name of the given 1-based Bikram
// Sambat month (1 is Baisakh, 12 is Chaitra).
func MonthNameNepali(month int) (string, error) {
	if month < 1 || month > monthsPerYear {
		return "", fmt.Errorf("%w: %d", ErrInvalidMonth, month)
	}
	return monthNamesNepali[month-1], nil
}

// MonthNameNepali returns the Devanagari name of d's month.
func (d Date) MonthNameNepali() (string, error) {
	return MonthNameNepali(d.Month)
}

// WeekdayNameNepali returns the Nepali (Devanagari) name of the given
// weekday, e.g. "आइतवार" for time.Sunday. It returns an error wrapping
// ErrInvalidWeekday if weekday is not in the range time.Sunday through
// time.Saturday.
func WeekdayNameNepali(weekday time.Weekday) (string, error) {
	if weekday < time.Sunday || weekday > time.Saturday {
		return "", fmt.Errorf("%w: %d", ErrInvalidWeekday, weekday)
	}
	return weekdayNamesNepali[weekday], nil
}

// WeekdayNameNepali returns the Nepali (Devanagari) name of the weekday d
// falls on, e.g. "बुधवार". It returns an error wrapping ErrInvalidYear,
// ErrInvalidMonth or ErrInvalidDay if d is not a valid date.
func (d Date) WeekdayNameNepali() (string, error) {
	weekday, err := d.DayOfWeek()
	if err != nil {
		return "", err
	}
	return WeekdayNameNepali(weekday)
}

// ToNepaliDigits replaces every ASCII digit (0-9) in s with its Devanagari
// equivalent (e.g. "2083" becomes "२०८३"). Non-digit characters, including
// punctuation and existing Devanagari digits, are copied through unchanged.
func ToNepaliDigits(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(nepaliDigits[r-'0'])
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// FromNepaliDigits replaces every Devanagari digit (०-९) in s with its ASCII
// equivalent (e.g. "२०८३" becomes "2083"). Other characters, including
// existing ASCII digits, are copied through unchanged.
func FromNepaliDigits(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if i := nepaliDigitIndex(r); i >= 0 {
			b.WriteByte('0' + byte(i))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func nepaliDigitIndex(r rune) int {
	for i, d := range nepaliDigits {
		if d == r {
			return i
		}
	}
	return -1
}
