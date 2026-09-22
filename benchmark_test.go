package bs

import (
	"testing"
	"time"
)

func BenchmarkADToBS(b *testing.B) {
	date := time.Date(2026, 9, 22, 0, 0, 0, 0, time.UTC)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := ADToBS(date); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBSToAD(b *testing.B) {
	d := Date{2083, 6, 6}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := BSToAD(d); err != nil {
			b.Fatal(err)
		}
	}
}
