# Changelog

## v0.3.0 - 2026-09-22

- `Date.Format(layout string)` — token-based layouts (`YYYY`, `YY`, `MMMM`,
  `MM`, `M`, `DD`, `D`, `dddd`, `ddd`)
- `MonthNameNepali`, `Date.MonthNameNepali`
- `ToNepaliDigits`, `FromNepaliDigits`

Planned for v0.4.0: calendar-grid helpers (`MonthCalendar`, `WeeksInMonth`,
`FirstWeekdayOfMonth`).

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
