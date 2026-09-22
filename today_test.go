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
	want, err := ADToBS(time.Now())
	if err != nil {
		t.Fatalf("ADToBS(time.Now()): %v", err)
	}
	if got != want {
		t.Errorf("TodayBS() = %v, want %v (ADToBS(time.Now()))", got, want)
	}
}
