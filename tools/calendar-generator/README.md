# calendar-generator

Maintainer tool that produces `../../data.go`, the static BS calendar
dataset used by the `bs` package. It exists so that how `data.go` was built
is visible and reproducible — the library itself has zero runtime
dependencies and makes no network requests; only this generator does.

See [`docs/calendar-data.md`](../../docs/calendar-data.md) for the full
write-up of sources, the verification method, and known limitations.

## What it does

- BS 1979–1999: parses [amitgaru/nepali-datetime](https://github.com/amitgaru/nepali-datetime)'s
  `calendar_bs.csv` (Apache-2.0).
- BS 2000–2100: queries Hamro Patro's own calendar data API directly (the
  same data its calendar pages render from), which pairs each day's
  Gregorian date with its Bikram Sambat date. For each year it requests
  months 2, 5, 8 and 11; each request reliably returns complete day data for
  the requested month and its immediate neighbors, which together cover all
  12 months.
- Validates every year 1979–2100 is present with 12 plausible month lengths
  before writing anything.
- Writes `data.go` with a `// Code generated ... DO NOT EDIT` header.

## Usage

```sh
go run ./tools/calendar-generator -out data.go
```

Requires network access. Takes a few minutes (roughly 400 HTTP requests to
Hamro Patro, made sequentially and politely rather than in a tight
concurrent loop).

## If you change this tool

If you find and fix a real calendar-data bug by changing what this tool
fetches or how it resolves disagreements, please also:

- Regenerate `data.go` and confirm `go test ./...` still passes (in
  particular `TestExhaustiveBSRoundTrip`/`TestExhaustiveADRoundTrip`, which
  will catch most internal inconsistencies).
- Update `docs/calendar-data.md` if the sourcing story changed.
- Note the fix in your pull request per `CONTRIBUTING.md`.
