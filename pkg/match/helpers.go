package match

import (
	"testing"
)

func verifyMatchType(t *testing.T, expected, got int) {
	t.Helper()

	if expected != got {
		t.Errorf("expected %d for match type, got %d", expected, got)
	}
}

func verifyNoMatch(t *testing.T, matches []string) {
	t.Helper()

	if len(matches) > 0 {
		t.Errorf("expected no match; match count was %d", len(matches))
	}
}

func verityMatch(t *testing.T, matches []string) {
	t.Helper()

	if len(matches) == 0 {
		t.Errorf("expected match count to be positive; got %d", len(matches))
	}
}
