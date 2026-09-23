# Changelog

## v0.7.0 - 2026-09-23

- `Date.FormatNepali(layout string)` — the same layout tokens as `Format`,
  rendered in Nepali: Devanagari digits, Nepali month names (`MMMM`),
  Nepali weekday names (`dddd`) and their short forms (`ddd`), e.g.
  `"dddd, MMMM D, YYYY"` renders as `"मंगलवार, असोज ६, २०८३"`. Literal
  characters in the layout are copied unchanged.
- `WeekdayNameNepali(time.Weekday)`, `Date.WeekdayNameNepali()` — Nepali
  weekday names, spelled as Hamro Patro's calendar spells them (आइतवार,
  सोमवार, मंगलवार, बुधवार, बिहिवार, शुक्रवार, शनिवार), the same source as
  the Nepali month names. The short forms used by `ddd` drop the "वार"
  suffix (आइत, सोम, …).
- New `ErrInvalidWeekday` sentinel, for a `time.Weekday` outside
  `time.Sunday`..`time.Saturday`.

## v0.6.1 - 2026-09-22

- Fixed: `TodayBS` now derives "today" from Nepal Standard Time
  (UTC+05:45) instead of the calling process's local timezone. On a
  server configured for UTC — a common default for cloud VMs and
  containers — `TodayBS` previously reported the wrong calendar day for
  roughly 5h45m of every day (the window after midnight has passed in
  Nepal but not yet in UTC). No API change; this is a behavior fix.

## v0.6.0 - 2026-09-22

- `Date.MarshalText` / `Date.UnmarshalText` (`encoding.TextMarshaler` /
  `TextUnmarshaler`) — a `Date` field now encodes as the plain string
  `"2083-06-06"` in `encoding/json` (and anything else that respects
  `TextMarshaler`: `encoding/gob`, `encoding/csv`, `url.Values`, map keys)
  instead of `{"Year":2083,"Month":6,"Day":6}`
- `Date.Value` / `Date.Scan` (`database/sql/driver.Valuer` /
  `database/sql.Scanner`) — a `Date` can be used directly as a
  `database/sql` query argument or scanned directly out of a row. `Scan`
  accepts a string/`[]byte` in `"YYYY-MM-DD"` form or a `time.Time` (for
  drivers that return native DATE columns as `time.Time`); `Value` always
  writes the `"YYYY-MM-DD"` string form

Neither `MarshalText`/`Value` validates its receiver (matching `String`'s
existing behavior); `UnmarshalText`/`Scan` validate the same way `Parse`
does.

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
