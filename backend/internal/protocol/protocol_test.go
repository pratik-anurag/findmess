package protocol

import (
	"testing"
	"time"
)

func TestRSSIBuckets(t *testing.T) {
	cases := map[int]string{
		-45: RSSINear,
		-60: RSSINear,
		-61: RSSIMedium,
		-78: RSSIMedium,
		-79: RSSIFar,
	}
	for rssi, want := range cases {
		if got := BucketRSSI(rssi); got != want {
			t.Fatalf("rssi %d got %s want %s", rssi, got, want)
		}
	}
}

func TestBucketTime(t *testing.T) {
	got := BucketTime(time.Date(2026, 5, 28, 10, 22, 30, 0, time.UTC))
	want := time.Date(2026, 5, 28, 10, 15, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("bucket mismatch: %s want %s", got, want)
	}
}
