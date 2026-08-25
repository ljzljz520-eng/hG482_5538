package domain

import "testing"

func TestNormalizeAndKey(t *testing.T) {
	if got := NormalizeName("  lIN  mei "); got != "Lin Mei" {
		t.Fatalf("got %q", got)
	}
	if got := NormalizePhone("138-0001"); got != "1380001" {
		t.Fatalf("got %q", got)
	}
	if got := RecordKey(" c1 ", "2026-01-01"); got != "c1:2026-01-01" {
		t.Fatalf("got %q", got)
	}
}
