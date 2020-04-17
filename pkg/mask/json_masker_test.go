package mask

import (
	"testing"
)

func TestJSONMask(t *testing.T) {
	t.Run("should mask 'password' by default", func(t *testing.T) {
		masker := NewJSONMasker([]string{}, false)
		output := getMaskOutputTrimmed(t, masker, `  "password": "secret"`)
		if output != `"password": "***"` {
			t.Errorf("'password' value not masked; got '%s'", output)
		}
	})

	t.Run("should not mask an unspecified property", func(t *testing.T) {
		masker := NewJSONMasker([]string{}, false)
		output := getMaskOutputTrimmed(t, masker, `  "prop": "value"`)
		if output != `"prop": "value"` {
			t.Errorf("unspecified property 'prop' should not be masked; got '%s'", output)
		}
	})

	t.Run("should mask specified property", func(t *testing.T) {
		masker := NewJSONMasker([]string{"prop"}, false)
		output := getMaskOutputTrimmed(t, masker, `  "prop": "value"`)
		if output != `"prop": "***"` {
			t.Errorf("property 'prop' should be masked; got '%s'", output)
		}
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
			output := getMaskOutputTrimmed(t, masker, input)
			if output != expected {
				t.Errorf("specified property was not masked, expected '%s', got '%s'", expected, output)
			}
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
			output := getMaskOutputTrimmed(t, masker, input)
			if output != expected {
				t.Errorf("specified property was not masked, expected '%s', got '%s'", expected, output)
			}
		})
	})
}
