# go-bs

[![test](https://github.com/suprimkhatri77/go-bs/actions/workflows/test.yml/badge.svg)](https://github.com/suprimkhatri77/go-bs/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/suprimkhatri77/go-bs.svg)](https://pkg.go.dev/github.com/suprimkhatri77/go-bs)
[![Go Report Card](https://goreportcard.com/badge/github.com/suprimkhatri77/go-bs)](https://goreportcard.com/report/github.com/suprimkhatri77/go-bs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A small, dependency-free Go library for converting dates between Gregorian
(AD) and Bikram Sambat (BS), the calendar used in Nepal.

```text
✓ 122 BS years covered (1979-2100)
✓ 1,464 BS months represented
✓ Every supported BS date tested
✓ AD → BS → AD round-trip tested
✓ BS → AD → BS round-trip tested
✓ Zero runtime dependencies
✓ Zero network requests
```

## Features

- `ADToBS` / `BSToAD` conversion
- Supports BS years 1979–2100 inclusive
- Strict date validation against real BS month lengths (not just shape)
- Date arithmetic and comparison (`AddDays`, `DaysBetween`, `Before`/`After`, ...)
- Layout-based formatting and Nepali-digit conversion
- Calendar-grid helpers for building calendar UIs (`MonthCalendar`, ...)
- Drop-in JSON encoding and `database/sql` support (`Date` implements
  `encoding.TextMarshaler`/`TextUnmarshaler` and
  `driver.Valuer`/`sql.Scanner`)
- Zero runtime dependencies
- Timezone-safe: conversion is based on calendar date, not time-of-day
- Table-driven calendar data, verified against a live source where possible
  (see [docs/calendar-data.md](docs/calendar-data.md))
- Exhaustively tested: every one of the 44,562 supported days round-trips
  in both directions

## Installation

```sh
go get github.com/suprimkhatri77/go-bs
```

## Usage

Convert AD to BS:

```go
package main

import (
	"fmt"
	"log"
	"time"

	bs "github.com/suprimkhatri77/go-bs"
)

func main() {
	ad := time.Date(2026, time.September, 22, 0, 0, 0, 0, time.UTC)

	d, err := bs.ADToBS(ad)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(d) // 2083-06-06
}
```

Convert BS to AD:

```go
d, err := bs.NewDate(2083, 6, 6)
if err != nil {
	log.Fatal(err)
}

ad, err := bs.BSToAD(d)
if err != nil {
	log.Fatal(err)
}

fmt.Println(ad.Format("2006-01-02")) // 2026-09-22
```

Validate a BS date without converting it:

```go
if !bs.IsValid(2083, 2, 32) { // Jestha 2083 has 31 days
	fmt.Println("not a real Bikram Sambat date")
}
```

Date arithmetic and comparison:

```go
d, _ := bs.NewDate(2083, 6, 15)

next, _ := d.AddDays(10)
fmt.Println(next) // 2083-06-25

fmt.Println(d.Before(next)) // true

weekday, _ := d.DayOfWeek()
fmt.Println(weekday) // Thursday

name, _ := d.MonthName()
fmt.Println(name) // Ashwin

diff, _ := bs.DaysBetween(d, next)
fmt.Println(diff) // 10

start, _ := d.StartOfMonth()
end, _ := d.EndOfMonth()
fmt.Println(start, end) // 2083-06-01 2083-06-31
```

Parsing:

```go
d, err := bs.Parse("2083-06-15")
if err != nil {
	log.Fatal(err)
}
```

Formatting and Nepali digits:

```go
d, _ := bs.NewDate(2083, 6, 6)

s, _ := d.Format("dddd, MMMM D, YYYY")
fmt.Println(s) // Tuesday, Ashwin 6, 2083

fmt.Println(bs.ToNepaliDigits(d.String())) // २०८३-०६-०६

name, _ := d.MonthNameNepali()
fmt.Println(name) // असोज
```

Building a calendar UI:

```go
weeks, _ := bs.MonthCalendar(2083, 6) // [][]*bs.Date, Sunday-first, nil-padded

for _, week := range weeks {
	for _, day := range week {
		if day == nil {
			fmt.Print("   ") // no day of this month in this cell
		} else {
			fmt.Printf("%2d ", day.Day)
		}
	}
	fmt.Println()
}
```

`TodayBS`, `MustParse`, `NextMonth`/`PreviousMonth`, and `Age`:

```go
today, _ := bs.TodayBS()
fmt.Println(today) // 2083-06-06 (whatever "today" is when this runs)

birth := bs.MustParse("2060-06-15")

next, _ := birth.NextMonth()
fmt.Println(next) // 2060-07-15

years, months, days, _ := bs.Age(birth, today)
fmt.Println(years, months, days) // 22 11 22
```

## API overview

```go
const (
	MinBSYear = 1979
	MaxBSYear = 2100
)

type Date struct {
	Year, Month, Day int
}

func NewDate(year, month, day int) (Date, error)
func Parse(s string) (Date, error)    // "YYYY-MM-DD"
func MustParse(s string) Date         // panics instead of erroring
func TodayBS() (Date, error)

func (d Date) Valid() bool
func (d Date) String() string // "YYYY-MM-DD"
func (d Date) MonthName() (string, error)
func (d Date) MonthNameNepali() (string, error)
func (d Date) Format(layout string) (string, error) // e.g. "YYYY-MM-DD"

func (d Date) Before(other Date) bool
func (d Date) After(other Date) bool
func (d Date) Equal(other Date) bool
func Compare(a, b Date) int // -1, 0, 1

func (d Date) AddDays(n int) (Date, error)
func (d Date) SubDays(n int) (Date, error)
func (d Date) NextDay() (Date, error)
func (d Date) PreviousDay() (Date, error)
func (d Date) NextMonth() (Date, error)     // clamps to target month's last day
func (d Date) PreviousMonth() (Date, error) // clamps to target month's last day
func (d Date) DayOfWeek() (time.Weekday, error)
func DaysBetween(a, b Date) (int, error)
func Age(birthBS, todayBS Date) (years, months, days int, err error)

func (d Date) StartOfMonth() (Date, error)
func (d Date) EndOfMonth() (Date, error)
func (d Date) StartOfYear() (Date, error)
func (d Date) EndOfYear() (Date, error)
func (d Date) DayOfYear() (int, error)

func ADToBS(t time.Time) (Date, error)
func BSToAD(d Date) (time.Time, error)

func IsValid(year, month, day int) bool
func IsSupportedBSYear(year int) bool
func DaysInMonth(year, month int) (int, error)
func DaysInYear(year int) (int, error)
func MonthName(month int) (string, error)
func MonthNameNepali(month int) (string, error)
func ToNepaliDigits(s string) string
func FromNepaliDigits(s string) string

func FirstWeekdayOfMonth(year, month int) (time.Weekday, error)
func WeeksInMonth(year, month int) (int, error)
func MonthCalendar(year, month int) ([][]*Date, error) // Sunday-first, nil-padded

func (d Date) MarshalText() ([]byte, error)  // encoding.TextMarshaler; same as String
func (d *Date) UnmarshalText(data []byte) error // encoding.TextUnmarshaler; same as Parse
func (d Date) Value() (driver.Value, error)  // database/sql/driver.Valuer
func (d *Date) Scan(value any) error         // database/sql.Scanner

var (
	ErrInvalidYear      error
	ErrInvalidMonth     error
	ErrInvalidDay       error
	ErrInvalidFormat    error
	ErrOutOfRange       error
	ErrInvalidDateOrder error
)
```

All errors support `errors.Is`, e.g. `errors.Is(err, bs.ErrInvalidDay)`.

`Date.Month` is 1-based: 1 is Baisakh, 12 is Chaitra.

`MonthCalendar`'s weeks run Sunday through Saturday, matching how Nepali
calendars are conventionally laid out (Hamro Patro included).

## Supported range

- Bikram Sambat: **1979–2100** (`MinBSYear`–`MaxBSYear`), inclusive.
- The corresponding Gregorian range is **1922-04-13 to 2044-04-13**, derived
  from the verified calendar data rather than assumed.

Most BS calendar libraries stop around this range because that's roughly the
edge of what Nepal's calendar authorities have officially published in
advance. `MaxBSYear` is expected to move out further (realistically not
before around BS 2095) once more official data exists.

## Accuracy

BS month lengths are not computed from a formula — they follow the
officially published Nepali calendar and are stored as a static,
table-driven dataset. That dataset was cross-checked against multiple
existing open-source implementations and, where possible, against a live
calendar source, rather than copied from a single upstream project as-is.
See [docs/calendar-data.md](docs/calendar-data.md) for the sources, the
verification method, and a documented known limitation (BS years
1979–1999 could not be checked against a live source and rely on two
library sources agreeing with each other).

## Timezone behavior

`ADToBS` looks only at the `Year`, `Month` and `Day` of the given
`time.Time` — its time-of-day and location are ignored, and the same
Gregorian calendar date always converts to the same BS date, regardless of
which timezone the `time.Time` is expressed in.

## Testing

```sh
go test ./...
go test -race ./...
go vet ./...
```

`TestExhaustiveBSRoundTrip` and `TestExhaustiveADRoundTrip` convert every one
of the 44,562 supported days in both directions and confirm the round trip
is exact, rather than relying on a sample.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).
