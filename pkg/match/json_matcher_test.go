package match

import (
	"testing"
)

func TestJSONMatch(t *testing.T) {
	m := NewJSONMatcher()

	t.Run("matches json line", func(t *testing.T) {
		propIndex, valueIndex, matches := m.Match(`"prop": "value"`)

		assertPropAndValue(t, propIndex, valueIndex)
		assertMatch(t, matches)
	})

	t.Run("does not match a non-json line", func(t *testing.T) {
		propIndex, valueIndex, matches := m.Match(`prop = value`)

		assertNoPropAndValue(t, propIndex, valueIndex)
		assertNoMatch(t, matches)
	})
}
