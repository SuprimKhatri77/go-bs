// Command calendar-generator produces the ../../data.go calendar dataset for
// package bs (github.com/suprimkhatri77/go-bs).
//
// It is a one-time/occasional maintainer tool, not part of the library: the
// library itself has no runtime dependencies and makes no network requests.
// Running this tool does require network access. See ../../docs/calendar-data.md
// for the full sourcing and verification write-up; this file only implements
// the mechanical extraction described there.
//
// Sourcing:
//   - BS 1979-1999: amitgaru/nepali-datetime's calendar_bs.csv (Apache-2.0),
//     cross-validated in docs/calendar-data.md against askbuddie/bikram-sambat
//     and against Hamro Patro's BS 2000-01-01 anchor date.
//   - BS 2000-2100: extracted directly from Hamro Patro's own calendar data
//     API (the same JSON its calendar pages render from), which records each
//     day's AD and BS date together. This is the live-authoritative source
//     used to catch and correct errors present in every library dataset
//     checked (see docs/calendar-data.md).
//
// Usage:
//
//	go run ./tools/calendar-generator -out data.go
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	minBSYear     = 1979
	maxBSYear     = 2100
	monthsPerYear = 12

	// amitgaruMaxYear is the last year (inclusive) sourced from amitgaru's
	// table; everything from amitgaruMaxYear+1 onward is sourced from Hamro
	// Patro instead.
	amitgaruMaxYear = 1999

	amitgaruCSVURL = "https://raw.githubusercontent.com/amitgaru/nepali-datetime/master/nepali_datetime/data/calendar_bs.csv"
	hamroURLFormat = "https://www.hamropatro.com/calendar/%d/%d"
)

func main() {
	out := flag.String("out", "data.go", "output path for the generated data.go")
	pkg := flag.String("package", "bs", "package name for the generated file")
	flag.Parse()

	rows := map[int][monthsPerYear]int{}

	amitRows, err := fetchAmitgaru(minBSYear, amitgaruMaxYear)
	if err != nil {
		log.Fatalf("fetching amitgaru data: %v", err)
	}
	for y, row := range amitRows {
		rows[y] = row
	}

	hamroRows, err := fetchHamro(amitgaruMaxYear+1, maxBSYear)
	if err != nil {
		log.Fatalf("fetching Hamro Patro data: %v", err)
	}
	for y, row := range hamroRows {
		rows[y] = row
	}

	if err := validate(rows); err != nil {
		log.Fatalf("validation failed: %v", err)
	}

	src, err := render(*pkg, rows)
	if err != nil {
		log.Fatalf("rendering source: %v", err)
	}

	if err := os.WriteFile(*out, []byte(src), 0o644); err != nil {
		log.Fatalf("writing %s: %v", *out, err)
	}
	fmt.Printf("wrote %s (%d years)\n", *out, len(rows))
}

// fetchAmitgaru downloads and parses amitgaru/nepali-datetime's calendar CSV,
// returning rows for years in [minYear, maxYear].
func fetchAmitgaru(minYear, maxYear int) (map[int][monthsPerYear]int, error) {
	resp, err := http.Get(amitgaruCSVURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %s fetching %s", resp.Status, amitgaruCSVURL)
	}

	reader := csv.NewReader(bufio.NewReader(resp.Body))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("empty CSV from %s", amitgaruCSVURL)
	}

	out := map[int][monthsPerYear]int{}
	for _, rec := range records[1:] { // skip header
		if len(rec) != monthsPerYear+1 {
			continue
		}
		year, err := strconv.Atoi(rec[0])
		if err != nil {
			continue
		}
		if year < minYear || year > maxYear {
			continue
		}
		var row [monthsPerYear]int
		for i := range monthsPerYear {
			v, err := strconv.Atoi(rec[i+1])
			if err != nil {
				return nil, fmt.Errorf("year %d: bad value %q: %w", year, rec[i+1], err)
			}
			row[i] = v
		}
		out[year] = row
	}
	for y := minYear; y <= maxYear; y++ {
		if _, ok := out[y]; !ok {
			return nil, fmt.Errorf("amitgaru data missing year %d", y)
		}
	}
	return out, nil
}

// dayEntryPattern matches one calendar-day record inside Hamro Patro's
// calendar page data, which pairs a Gregorian date with the corresponding
// Bikram Sambat date for every day it renders.
var dayEntryPattern = regexp.MustCompile(
	`"year_ad":(\d+),"month_ad":(\d+),"day_ad":(\d+),"year_bs":(\d+),"month_bs":(\d+),"day_bs":(\d+)`,
)

// fetchHamro extracts BS month lengths for years [minYear, maxYear] from
// Hamro Patro. For each year it requests months 2, 5, 8 and 11; each request
// reliably returns complete (gap-free) day data for the requested month and
// its immediate neighbors, which together cover all 12 months of the year.
func fetchHamro(minYear, maxYear int) (map[int][monthsPerYear]int, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	maxDay := map[[2]int]int{} // [year, month] -> highest day_bs seen in a complete run

	for year := minYear; year <= maxYear; year++ {
		for _, anchorMonth := range []int{2, 5, 8, 11} {
			url := fmt.Sprintf(hamroURLFormat, year, anchorMonth)
			body, err := getWithRetry(client, url, 3)
			if err != nil {
				return nil, fmt.Errorf("fetching %s: %w", url, err)
			}

			byMonth := map[[2]int]map[int]bool{}
			for _, m := range dayEntryPattern.FindAllStringSubmatch(body, -1) {
				yb, _ := strconv.Atoi(m[4])
				mb, _ := strconv.Atoi(m[5])
				db, _ := strconv.Atoi(m[6])
				key := [2]int{yb, mb}
				if byMonth[key] == nil {
					byMonth[key] = map[int]bool{}
				}
				byMonth[key][db] = true
			}

			// Only trust the anchor month and its immediate neighbors: a
			// single request also includes a partial preview of a third,
			// more distant month, which is not reliable for a max-day count.
			for _, delta := range []int{-1, 0, 1} {
				y, m := addMonths(year, anchorMonth, delta)
				days, ok := byMonth[[2]int{y, m}]
				if !ok || !isCompleteRun(days) {
					continue
				}
				n := len(days)
				if cur, ok := maxDay[[2]int{y, m}]; !ok || n > cur {
					maxDay[[2]int{y, m}] = n
				}
			}
		}
	}

	out := map[int][monthsPerYear]int{}
	for year := minYear; year <= maxYear; year++ {
		var row [monthsPerYear]int
		for month := 1; month <= monthsPerYear; month++ {
			n, ok := maxDay[[2]int{year, month}]
			if !ok {
				return nil, fmt.Errorf("no data for BS %d-%02d", year, month)
			}
			row[month-1] = n
		}
		out[year] = row
	}
	return out, nil
}

func addMonths(year, month, delta int) (int, int) {
	m := month + delta
	for m < 1 {
		m += monthsPerYear
		year--
	}
	for m > monthsPerYear {
		m -= monthsPerYear
		year++
	}
	return year, m
}

// isCompleteRun reports whether days is exactly {1, 2, ..., len(days)}, i.e.
// a full month was observed with no gaps.
func isCompleteRun(days map[int]bool) bool {
	for i := 1; i <= len(days); i++ {
		if !days[i] {
			return false
		}
	}
	return len(days) > 0
}

func getWithRetry(client *http.Client, url string, attempts int) (string, error) {
	var lastErr error
	for range attempts {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("RSC", "1") // request Hamro Patro's React Server Component data payload
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Second)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			time.Sleep(time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status %s", resp.Status)
			time.Sleep(time.Second)
			continue
		}
		return string(body), nil
	}
	return "", lastErr
}

func validate(rows map[int][monthsPerYear]int) error {
	for y := minBSYear; y <= maxBSYear; y++ {
		row, ok := rows[y]
		if !ok {
			return fmt.Errorf("missing year %d", y)
		}
		total := 0
		for _, d := range row {
			if d < 28 || d > 32 {
				return fmt.Errorf("year %d: implausible month length %d", y, d)
			}
			total += d
		}
		if total < 355 || total > 375 {
			return fmt.Errorf("year %d: implausible year length %d", y, total)
		}
	}
	return nil
}

const fileTemplate = `// Code generated by tools/calendar-generator. DO NOT EDIT.
//
// Regenerate with: go run ./tools/calendar-generator -out data.go
// See docs/calendar-data.md for sources and verification method.
package %s

// calendarData holds the number of days in each Bikram Sambat month for every
// year from MinBSYear to MaxBSYear, inclusive. Row index 0 is MinBSYear.
var calendarData = [MaxBSYear - MinBSYear + 1][monthsPerYear]uint8{
%s}
`

func render(pkg string, rows map[int][monthsPerYear]int) (string, error) {
	years := make([]int, 0, len(rows))
	for y := range rows {
		years = append(years, y)
	}
	sort.Ints(years)

	var b strings.Builder
	for _, y := range years {
		row := rows[y]
		fmt.Fprintf(&b, "\t{")
		for i, d := range row {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%d", d)
		}
		fmt.Fprintf(&b, "}, // %d\n", y)
	}

	return fmt.Sprintf(fileTemplate, pkg, b.String()), nil
}
