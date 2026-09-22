package bs

import "testing"

func TestCalendarDataIntegrity(t *testing.T) {
	wantYears := MaxBSYear - MinBSYear + 1
	if len(calendarData) != wantYears {
		t.Fatalf("len(calendarData) = %d, want %d", len(calendarData), wantYears)
	}

	for i, row := range calendarData {
		year := MinBSYear + i
		if len(row) != monthsPerYear {
			t.Fatalf("year %d has %d months, want %d", year, len(row), monthsPerYear)
		}
		total := 0
		for m, days := range row {
			if days < 29 || days > 32 {
				t.Errorf("year %d month %d has %d days, want a value in 29..32", year, m+1, days)
			}
			total += int(days)
		}
		if total < 360 || total > 370 {
			t.Errorf("year %d totals %d days, want a value in 360..370", year, total)
		}
	}
}

func TestReferenceDate(t *testing.T) {
	ad, err := BSToAD(referenceBS)
	if err != nil {
		t.Fatalf("BSToAD(referenceBS): %v", err)
	}
	if !ad.Equal(referenceAD) {
		t.Errorf("BSToAD(referenceBS) = %s, want %s", ad.Format("2006-01-02"), referenceAD.Format("2006-01-02"))
	}

	bsDate, err := ADToBS(referenceAD)
	if err != nil {
		t.Fatalf("ADToBS(referenceAD): %v", err)
	}
	if bsDate != referenceBS {
		t.Errorf("ADToBS(referenceAD) = %v, want %v", bsDate, referenceBS)
	}
}
