# Calendar data: sources and verification

This package supports Bikram Sambat (BS) years `1979` through `2100`
inclusive (`MinBSYear`/`MaxBSYear` in `date.go`). BS month lengths are not
computed from a formula — they follow Nepal's officially published calendar
and are recorded as a static, table-driven dataset in `data.go`.

## Why this needed real verification

Every BS/AD converter is only as good as its calendar data. While researching
this package, three actively used, independent open-source calendar
libraries were compared against each other and against a live authoritative
source, and **all three were found to contain real errors** in at least a
few individual years. Two libraries agreeing with each other was not, by
itself, reliable evidence of correctness — two projects can share a common
(and equally wrong) upstream data table. This is why the process below
leans on a live, independently queryable source wherever one was available,
rather than trusting library consensus alone.

## Sources consulted

- **[amitgaru/nepali-datetime](https://github.com/amitgaru/nepali-datetime)**
  (Python, Apache-2.0). Ships a plain CSV of BS month lengths for years
  1975–2100, with a documented reference pair (BS 1975-01-01 = AD
  1918-04-13).
- **[askbuddie/bikram-sambat](https://github.com/askbuddie/bikram-sambat)**
  (TypeScript, MIT). Ships an explicit `DaysInMonthsMappingData` table for BS
  1975–2100.
- **[puncoz-official/bikram-sambat-js](https://github.com/puncoz-official/bikram-sambat-js)**
  (TypeScript, MIT, published as `bikram-sambat-js` on npm). This is the
  library the originating project's frontend already uses to convert BS to
  AD before sending both values to the backend, so consistency with it
  matters operationally. Its data is a run-length-encoded table for BS
  1970–2100.
- **[Hamro Patro](https://www.hamropatro.com/calendar)**, Nepal's most
  widely used calendar site/app. Not open source and not used as copied
  code, but used as a live, authoritative reference to adjudicate
  disagreements between the library sources above, by directly reading each
  day's structured calendar record it serves (year/month/day in both AD and
  BS for every day of a given month). Hamro Patro's own data only goes back
  to BS 2000; it has no record of years before that.

## What was found

Cross-comparing amitgaru and askbuddie's tables across the full 1979–2100
range, they agreed on 123 of 126 years exactly, disagreeing only on BS 2004,
2082 and 2083. That level of agreement looks like strong corroboration, but
checking specific disputed and non-disputed months directly against Hamro
Patro's live data showed that agreement between the two libraries did not
reliably predict correctness — for example, both libraries agreed with each
other on BS 2010 Baisakh (31 days), and both were wrong (Hamro Patro's own
day-by-day records show 31 days is right — an earlier, less careful
read of Hamro Patro's *rendered page text* misread this as 30, before the
structured data endpoint was used instead; see "Method" below for why the
rendered-text approach was abandoned). `puncoz-official/bikram-sambat-js`
also disagreed with the amitgaru/askbuddie pairing for a cluster of years in
the 1979–1998 range, and its own AD-range validation has an off-by-one bug
(it derives its maximum supported AD year as `maxBSYear - 57`, a fixed
offset, which is one year too low for BS dates that fall late in the BS year
— e.g. it rejects `BSToAD("2100-12-01")` even though its own data covers BS
2100).

## Method used to build `data.go`

1. **BS 2000–2100** (101 years): fetched directly from Hamro Patro's
   internal calendar data (the JSON payload its calendar pages are rendered
   from, which records `year_ad`/`month_ad`/`day_ad` alongside
   `year_bs`/`month_bs`/`day_bs` for every calendar day). For each
   (year, month), the number of days was taken as the highest `day_bs` seen
   for that (year_bs, month_bs) pair, only trusting months that appeared as
   a complete, gap-free `1..N` sequence in the response (a given request
   returns the requested month plus one month on each side in full, plus a
   partial trailing preview of a third month, which was discarded). This
   was fetched for every year 2000–2100, cross-checked for internal
   consistency (every month resolved to a clean `1..N` sequence, zero
   parse errors across 1,212 month-fetches), and used as-is. This is why
   `data.go`'s BS 2000–2100 rows sometimes differ from all three library
   sources above.
2. **BS 1979–1999** (21 years): Hamro Patro has no data this far back, so
   these years use the amitgaru table (which agrees with askbuddie for all
   of them). As a sanity check, the running total of days across BS
   1979–1999 from this table was added to the verified BS 1979-01-01
   anchor and confirmed to land exactly on Hamro Patro's own BS 2000-01-01
   date (AD 1943-04-14) — i.e. the *year-boundary* dates for this span are
   corroborated even though individual month lengths within it are not
   independently verified against a live source. **This is a known
   limitation**: unlike 2000–2100, the 1979–1999 rows have not been checked
   against a live calendar day-by-day, only cross-validated between two
   library sources plus this boundary check. Anyone with access to physical
   Patro almanacs for this period, or another live source that covers it,
   is encouraged to verify and send corrections.

A note on the abandoned first attempt: the initial approach to querying
Hamro Patro was to render its calendar page and count day numbers in the
visible text. This was dropped after it produced at least one confirmed
wrong reading (BS 2010 Baisakh, misread as 30 instead of 31) due to the
page including duplicate/adjacent-month content in its rendered text that
was easy to miscount by hand or with a naive scan. The structured JSON
approach in step 1 above does not have this problem and is what `data.go`
is actually built from for 2000–2100.

## Reference date

The package's internal anchor is:

```text
BS 1979-01-01 = AD 1922-04-13
```

This is `MinBSYear-01-01`, chosen so every supported date has a
non-negative offset from the reference, which keeps the conversion
algorithm in `convert.go` simple (no negative-offset case to handle). It was
derived by summing amitgaru's day counts for BS 1975–1978 starting from its
documented BS 1975-01-01 = AD 1918-04-13 anchor, and cross-checked
independently by asking `bikram-sambat-js` (the frontend's own library) for
`BSToAD("1979-01-01")`, which returned the same date, `1922-04-13`.

The maximum supported date, BS 2100-12-31, corresponds to **AD 2044-04-13**
(verified directly from Hamro Patro's data, as described above).

Both ends of the range, and the current date at the time this package was
built (BS 2083-06-06 = AD 2026-09-22), are covered by table-driven tests in
`convert_test.go`.

## Reproducing this data

The extraction described above is implemented as a small maintainer tool in
[`tools/calendar-generator`](../tools/calendar-generator), which regenerates
`data.go` from the two sources above. It is not part of the library — the
library itself has no runtime dependencies and makes no network calls; only
this generator does, and only when a maintainer runs it deliberately. See
that directory's README for usage.

## For contributors

If you're changing `calendarData`:

- Cite a source for the new values in your pull request description — ideally
  one that's independently checkable (a live calendar, an official
  publication), not just another library's table.
- Add or update the relevant entry in `knownPairs` (`convert_test.go`) if you
  can verify an AD/BS pair for the year you're touching.
- Run `go test ./...` — `TestExhaustiveBSRoundTrip` and
  `TestExhaustiveADRoundTrip` will catch most internal inconsistencies
  introduced by a bad edit, and `TestCalendarDataIntegrity` checks basic
  shape (12 months per year, plausible day counts, plausible year totals).
- If you find and fix a specific error, please also open an issue or PR
  against whichever of the upstream libraries above you can, so the wider
  ecosystem benefits too.

The plan is to extend `MaxBSYear` beyond 2100 once Nepal's calendar
authorities have published data far enough in advance to responsibly do so
(realistically not before around BS 2095).
