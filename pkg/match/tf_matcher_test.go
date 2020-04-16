package match

import (
	"testing"
)

func TestTFMatch(t *testing.T) {
	matcher := NewTFMatcher()

	t.Run("should match new resource", func(t *testing.T) {
		line := ` + my_prop = "my_value"`
		valueIndex, matches := matcher.Match(line)

		assertPropAndValue(t, valueIndex)
		assertMatch(t, matches)
	})

	t.Run("should match removed resource", func(t *testing.T) {
		line := ` - my_prop = "my_value"`
		valueIndex, matches := matcher.Match(line)

		assertPropAndValue(t, valueIndex)
		assertMatch(t, matches)
	})

	t.Run("should match renamed resource", func(t *testing.T) {
		line := ` ~ my_prop = "old_value" -> "new_value"`
		valueIndex, matches := matcher.Match(line)

		assertPropAndValue(t, valueIndex)
		assertMatch(t, matches)
	})

	t.Run("should match removed resource to null", func(t *testing.T) {
		line := ` - my_prop = "old_value" -> null`
		valueIndex, matches := matcher.Match(line)

		assertPropAndValue(t, valueIndex)
		assertMatch(t, matches)
	})

	t.Run("should not match known after apply for new resource", func(t *testing.T) {
		line := ` ~ my_prop = (known after apply)`
		valueIndex, matches := matcher.Match(line)

		assertPropAndValue(t, valueIndex)
		assertNoMatch(t, matches)
	})

	t.Run("should match known after apply for changed resource", func(t *testing.T) {
		line := ` ~ my_prop = "old_value" -> (known after apply)`
		valueIndex, matches := matcher.Match(line)

		assertPropAndValue(t, valueIndex)
		assertMatch(t, matches)
	})
}
