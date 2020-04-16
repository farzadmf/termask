package match

import (
	"testing"
)

func assertNoPropAndValue(t *testing.T, propIndex, valueIndex int) {
	t.Helper()

	if propIndex != -1 || valueIndex != -1 {
		t.Errorf("expected 'prop' and 'value' indices to be -1, got %d and %d", propIndex, valueIndex)
	}
}

func assertPropAndValue(t *testing.T, propIndex, valueIndex int) {
	t.Helper()

	if propIndex < 0 || valueIndex < 0 {
		t.Errorf("expected 'prop' and 'value' indices to be > 0, got %d and %d", propIndex, valueIndex)
	}
}

func assertNoMatch(t *testing.T, matches []string) {
	t.Helper()

	if len(matches) > 0 {
		t.Errorf("expected no match; match count was %d", len(matches))
	}
}

func assertMatch(t *testing.T, matches []string) {
	t.Helper()

	if len(matches) == 0 {
		t.Errorf("expected match count to be positive; got %d", len(matches))
	}
}
