package bs

import (
	"testing"
	"time"
)

func TestTodayBS(t *testing.T) {
	got, err := TodayBS()
	if err != nil {
		t.Fatalf("TodayBS: %v", err)
	}
	want, err := ADToBS(time.Now().In(nepalTime))
	if err != nil {
		t.Fatalf("ADToBS(time.Now().In(nepalTime)): %v", err)
	}
	if got != want {
		t.Errorf("TodayBS() = %v, want %v (Nepal's current date)", got, want)
	}
}

// TestTodayBSIgnoresLocalTimezone guards the actual bug this function
// exists to avoid: on a server whose OS timezone is UTC (a common default
// for cloud VMs and containers), naively using time.Now() without
// converting to Nepal time would report the wrong calendar day for part of
// every day. This locks in an instant, deep into the UTC evening, where
// Nepal (UTC+05:45) has already moved into the next calendar day while UTC
// has not, and checks TodayBS's underlying conversion against Nepal time
// rather than raw local time.
func TestTodayBSIgnoresLocalTimezone(t *testing.T) {
	utcLateEvening := time.Date(2026, time.September, 22, 20, 0, 0, 0, time.UTC)
	inNepal := utcLateEvening.In(nepalTime)
	if inNepal.Format("2006-01-02") == utcLateEvening.Format("2006-01-02") {
		t.Fatalf("test setup invalid: chosen instant doesn't straddle the UTC/Nepal day boundary")
	}
	got, err := ADToBS(inNepal)
	if err != nil {
		t.Fatalf("ADToBS: %v", err)
	}
	wrongIfUsingUTC, err := ADToBS(utcLateEvening)
	if err != nil {
		t.Fatalf("ADToBS: %v", err)
	}
	if got == wrongIfUsingUTC {
		t.Fatalf("expected Nepal's date to differ from raw UTC's date at this instant, both gave %v", got)
	}
}
