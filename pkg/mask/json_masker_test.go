package mask

import (
	"testing"
)

func TestJSONMask(t *testing.T) {
	t.Run("should mask 'password' by default", func(t *testing.T) {
		masker := NewJSONMasker([]string{}, false)
		got := getMaskOutput(t, masker, `  "password": "secret"`)
		expected := `  "password": "***"`
		assertMatch(t, got, expected)
	})

	t.Run("should not mask an unspecified property", func(t *testing.T) {
		masker := NewJSONMasker([]string{}, false)
		got := getMaskOutput(t, masker, `  "prop": "value"`)
		expected := `  "prop": "value"`
		assertMatch(t, got, expected)
	})

	t.Run("should mask specified property", func(t *testing.T) {
		masker := NewJSONMasker([]string{"prop"}, false)
		got := getMaskOutput(t, masker, `  "prop": "value"`)
		expected := `  "prop": "***"`
		assertMatch(t, got, expected)
	})

	t.Run("multi-line string", func(t *testing.T) {
		t.Run("should mask specified property", func(t *testing.T) {
			masker := NewJSONMasker([]string{"prop"}, false)
			input := `{
  "prop": "value",
  "otherProp": "otherValue"
}`

			expected := `{
  "prop": "***",
  "otherProp": "otherValue"
}`
			got := getMaskOutput(t, masker, input)
			assertMatch(t, got, expected)
		})

		t.Run("should not mask unspecified properties", func(t *testing.T) {
			masker := NewJSONMasker([]string{}, false)
			input := `{
  "prop": "value",
  "otherProp": "otherValue"
}`

			expected := `{
  "prop": "value",
  "otherProp": "otherValue"
}`
			got := getMaskOutput(t, masker, input)
			assertMatch(t, got, expected)
		})
	})
}
