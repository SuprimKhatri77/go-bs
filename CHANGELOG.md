# Changelog

## v0.4.0 - 2026-09-22

- `MonthCalendar(year, month int) ([][]*Date, error)` — Sunday-first,
  nil-padded week grid for building calendar UIs
- `WeeksInMonth`, `FirstWeekdayOfMonth`

This completes the API surface originally planned across v0.1.0–v0.4.0.

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
