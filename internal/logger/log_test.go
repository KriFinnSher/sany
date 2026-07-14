package logger

import (
	"log/slog"
	"testing"
	"time"
)

func TestReplaceTimeFormat(t *testing.T) {
	timestamp := time.Date(2026, time.July, 13, 20, 37, 16, 604165000, time.FixedZone("MSK", 3*60*60))

	got := replaceTimeFormat(nil, slog.Time(slog.TimeKey, timestamp))

	if got.Value.Kind() != slog.KindString {
		t.Fatalf("time kind = %v, want %v", got.Value.Kind(), slog.KindString)
	}
	if got.Value.String() != "2026-07-13|20:37" {
		t.Errorf("time = %q, want %q", got.Value.String(), "2026-07-13|20:37")
	}
}
