package match

import (
	"testing"
)

func TestMatch(t *testing.T) {
	matcher := NewTFMatcher()

	verifyMatchType := func(t *testing.T, expected, got int) {
		t.Helper()

		if expected != got {
			t.Errorf("expected %d for match type, got %d", expected, got)
		}
	}

	verifyNoMatch := func(t *testing.T, matches []string) {
		if len(matches) > 0 {
			t.Errorf("expected no match; match count was %d", len(matches))
		}
	}

	verityMatch := func(t *testing.T, matches []string) {
		if len(matches) == 0 {
			t.Errorf("expected match count to be positive; got %d", len(matches))
		}
	}

	t.Run("should match new resource", func(t *testing.T) {
		line := ` + my_prop = "my_value"`
		matchType, matches := matcher.Match(line)

		verifyMatchType(t, TFNewOrRemove, matchType)
		verityMatch(t, matches)
	})

	t.Run("should match removed resource", func(t *testing.T) {
		line := ` - my_prop = "my_value"`
		matchType, matches := matcher.Match(line)

		verifyMatchType(t, TFNewOrRemove, matchType)
		verityMatch(t, matches)
	})

	t.Run("should match renamed resource", func(t *testing.T) {
		line := ` ~ my_prop = "old_value" -> "new_value"`
		matchType, matches := matcher.Match(line)

		verifyMatchType(t, TFReplace, matchType)
		verityMatch(t, matches)
	})

	t.Run("should match removed resource to null", func(t *testing.T) {
		line := ` - my_prop = "old_value" -> null`
		matchType, matches := matcher.Match(line)

		verifyMatchType(t, TFRemoveToNull, matchType)
		verityMatch(t, matches)
	})

	t.Run("should not match known after apply for new resource", func(t *testing.T) {
		line := ` ~ my_prop = (known after apply)`
		matchType, matches := matcher.Match(line)

		verifyMatchType(t, None, matchType)
		verifyNoMatch(t, matches)
	})

	t.Run("should match known after apply for changed resource", func(t *testing.T) {
		line := ` ~ my_prop = "old_value" -> (known after apply)`
		matchType, matches := matcher.Match(line)

		verifyMatchType(t, TFReplaceKnownAfterApply, matchType)
		verityMatch(t, matches)
	})
}
