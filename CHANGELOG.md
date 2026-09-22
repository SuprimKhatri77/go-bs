# Changelog

## Unreleased (accumulating toward v0.2.0)

- `Compare`, `Date.SubDays`, `DaysBetween`
- `Date.StartOfMonth`, `Date.EndOfMonth`, `Date.StartOfYear`, `Date.EndOfYear`, `Date.DayOfYear`

Still to come before v0.2.0 ships: `Format`/layout-based formatting,
Nepali-digit conversion, `MonthNameNepali`, and calendar-grid helpers.

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
