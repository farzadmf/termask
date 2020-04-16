package match

import (
	"testing"
)

func TestJSONMatch(t *testing.T) {
	m := NewJSONMatcher()

	t.Run("matches json line", func(t *testing.T) {
		match, matches := m.Match(`"prop": "value"`)

		verifyMatchType(t, JSONLine, match)
		verityMatch(t, matches)
	})

	t.Run("does not match a non-json line", func(t *testing.T) {
		match, matches := m.Match(`prop = value`)

		verifyMatchType(t, None, match)
		verifyNoMatch(t, matches)
	})
}
