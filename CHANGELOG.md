# Changelog

## v0.5.0 - 2026-09-22

- `TodayBS() (Date, error)` — current system date as a BS `Date`
- `MustParse(s string) Date` — panics instead of erroring; documented as such
- `Date.NextDay`/`Date.PreviousDay`
- `Date.NextMonth`/`Date.PreviousMonth` — clamp to the target month's last
  day rather than rolling over (e.g. day 31 in a 30-day month becomes day 30)
- `Age(birthBS, todayBS Date) (years, months, days int, err error)` —
  calendar age, not just a day count; new `ErrInvalidDateOrder` sentinel for
  when birthBS is after todayBS

These were the remaining items from the original API wishlist that didn't
make the v0.1.0–v0.4.0 phased plan.

## v0.4.0 - 2026-09-22

- `MonthCalendar(year, month int) ([][]*Date, error)` — Sunday-first,
  nil-padded week grid for building calendar UIs
- `WeeksInMonth`, `FirstWeekdayOfMonth`

## v0.3.0 - 2026-09-22

- `Date.Format(layout string)` — token-based layouts (`YYYY`, `YY`, `MMMM`,
  `MM`, `M`, `DD`, `D`, `dddd`, `ddd`)
- `MonthNameNepali`, `Date.MonthNameNepali`
- `ToNepaliDigits`, `FromNepaliDigits`

## v0.2.0 - 2026-09-22

- `Compare`, `Date.SubDays`, `DaysBetween`
- `Date.StartOfMonth`, `Date.EndOfMonth`, `Date.StartOfYear`, `Date.EndOfYear`, `Date.DayOfYear`

## v0.1.0 - 2026-09-22

- `ADToBS` / `BSToAD` conversion
- `NewDate`, `Parse`, `Date.String`, `Date.MonthName`
- `Date.Before`, `Date.After`, `Date.Equal`
- `Date.AddDays`, `Date.DayOfWeek`
- `IsValid`, `IsSupportedBSYear`, `DaysInMonth`, `DaysInYear`, `MonthName`
- Supports BS 1979–2100
- Static, verified calendar dataset (see `docs/calendar-data.md`)
- Exhaustive round-trip tests over every supported day
- Zero runtime dependencies
