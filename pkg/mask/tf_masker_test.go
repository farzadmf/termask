package mask

import (
	"testing"
)

func TestMask(t *testing.T) {
	t.Run("should mask 'password' by default", func(t *testing.T) {
		masker := NewTFMasker([]string{}, false)
		output := getMaskOutputTrimmed(t, masker, ` + password = "value"`)
		if output != `+ password = "***"` {
			t.Errorf("'password' value was not masked; got '%s'", output)
		}
	})

	t.Run("should mask 'PaSSworD' by default", func(t *testing.T) {
		masker := NewTFMasker([]string{}, false)
		output := getMaskOutputTrimmed(t, masker, ` + PaSSworD = "value"`)
		if output != `+ PaSSworD = "***"` {
			t.Errorf("'PaSSworD' value was not masked; got '%s'", output)
		}
	})

	t.Run("should mask 'My_PassWord' by default", func(t *testing.T) {
		masker := NewTFMasker([]string{}, false)
		trimmedOutput := getMaskOutputTrimmed(t, masker, ` + My_PassWord = "value"`)
		if trimmedOutput != `+ My_PassWord = "***"` {
			t.Errorf("'My_PassWord' value was not masked; got '%s'", trimmedOutput)
		}
	})

	t.Run("should mask custom property case sensitive", func(t *testing.T) {
		masker := NewTFMasker([]string{"my_prop"}, false)
		trimmedOutput := getMaskOutputTrimmed(t, masker, ` + my_prop = "value"`)
		if trimmedOutput != `+ my_prop = "***"` {
			t.Errorf("did not mask custom property; got '%s'", trimmedOutput)
		}
	})

	t.Run("should mask custom property ignoring case", func(t *testing.T) {
		masker := NewTFMasker([]string{"my_prop"}, true)
		trimmedOutput := getMaskOutputTrimmed(t, masker, ` + My_PrOP = "value"`)
		if trimmedOutput != `+ My_PrOP = "***"` {
			t.Errorf("did not mask custom property, case insensitive; got '%s'", trimmedOutput)
		}
	})

	t.Run("should not mask when property doesn't match", func(t *testing.T) {
		masker := NewTFMasker([]string{"my_prop"}, true)
		trimmedOutput := getMaskOutputTrimmed(t, masker, ` + other_prop = "value"`)
		if trimmedOutput != `+ other_prop = "value"` {
			t.Errorf("did not print property as is when no match; got '%s'", trimmedOutput)
		}
	})
}
