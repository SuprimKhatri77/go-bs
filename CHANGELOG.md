# Changelog

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

## Planned for v0.2.0

- `Compare`, `SubDays`, `DaysBetween`
- `StartOfMonth`, `EndOfMonth`, `StartOfYear`, `EndOfYear`, `DayOfYear`
- `Format`/layout-based formatting, Nepali-digit conversion
  (`ToNepaliDigits`/`FromNepaliDigits`), `MonthNameNepali`
- Calendar-grid helpers (`MonthCalendar`, `WeeksInMonth`,
  `FirstWeekdayOfMonth`) for building calendar UIs
