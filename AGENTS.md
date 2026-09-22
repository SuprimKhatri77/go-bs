# AGENTS.md

Guidance for AI coding agents working in this repository. Project facts and
conventions live here; see `CONTRIBUTING.md` for the human-facing version
(they should stay consistent).

## What this is

`go-bs` (`github.com/suprimkhatri77/go-bs`) is a small, zero-dependency Go
library converting between the Gregorian (AD) and Bikram Sambat (BS)
calendars, for BS years `MinBSYear`–`MaxBSYear` (1979–2100). Single package
`bs` at the repo root, one file per concern:

| File | Concern |
| --- | --- |
| `date.go` | `Date` type, `NewDate`, `MonthName` |
| `convert.go` | `ADToBS`, `BSToAD` |
| `validate.go` | `IsValid`, `DaysInMonth`, `DaysInYear`, `IsSupportedBSYear` |
| `compare.go` | `Compare`, `Before`/`After`/`Equal` |
| `arithmetic.go` | `AddDays`, `DaysBetween`, `StartOfMonth`/`EndOfMonth`, ... |
| `calendar.go` | `MonthCalendar`, `WeeksInMonth`, `FirstWeekdayOfMonth` |
| `format.go` | `Date.Format` (token-based layouts) |
| `nepali.go` | `MonthNameNepali`, `ToNepaliDigits`/`FromNepaliDigits` |
| `parse.go` | `Parse`, `MustParse` |
| `today.go` | `TodayBS` |
| `age.go` | `Age` |
| `text.go` | `encoding.TextMarshaler`/`TextUnmarshaler` |
| `sql.go` | `database/sql/driver.Valuer`/`database/sql.Scanner` |
| `errors.go` | sentinel errors, checked via `errors.Is` |
| `data.go` | generated — see below |

## Before considering any change done

```sh
gofmt -l .
go vet ./...
go test ./...
go test -race ./...
```

All four must be clean. CI (`.github/workflows/test.yml`) runs the same
checks. Any `.md` change should also pass:

```sh
npx --yes markdownlint-cli2 "LICENSE" "**/*.md" "#node_modules" "#CLAUDE.md"
```

New exported symbols need a doc comment starting with their own name (what
pkg.go.dev renders).

## `data.go` is generated — never hand-edit it

It carries a `// Code generated ... DO NOT EDIT` header and is produced by
`tools/calendar-generator` (see that directory's `README.md` and
`docs/calendar-data.md` for the full sourcing/verification write-up).
Regenerate via that tool. If hand-patching a single verified value is truly
unavoidable, say so explicitly and cite a source in the commit message.

## Calendar data changes need a cited source

Any change to BS month lengths needs, in the PR/commit: the affected BS
year(s)/month(s), and an independently-checkable source (a live calendar, an
official publication) — not just another library's table. `docs/calendar-data.md`
documents why two libraries agreeing with each other turned out not to be
reliable evidence of correctness. Where possible, add a verified AD/BS pair
to `knownPairs` in `convert_test.go`.

Fiscal-year-style facts (government policy, not calendar astronomy) need
the same rigor but a different kind of source — see the note in
`CHANGELOG.md`'s v0.6.0 entry about why fiscal-year helpers were deliberately
left unimplemented rather than guessed.

## Scope

Keep changes to what this library is: BS/AD conversion, validation, date
arithmetic/formatting, and their tests/docs. Frontend/UI code and other
calendar systems belong in a different project (see `go-bs-docs` for the
documentation site, a separate repo).
