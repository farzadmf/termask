package mask

import (
	"testing"
)

func TestMask(t *testing.T) {
	t.Run("new resource", func(t *testing.T) {
		t.Run("should mask 'password' by default", func(t *testing.T) {
			masker := NewTFMasker([]string{}, false)
			got := getMaskOutput(t, masker, ` + password = "value"`)
			expected := ` + password = "***"`
			assertMatch(t, got, expected)
		})

		t.Run("should mask 'PaSSworD' by default", func(t *testing.T) {
			masker := NewTFMasker([]string{}, false)
			got := getMaskOutput(t, masker, ` + PaSSworD = "value"`)
			expected := ` + PaSSworD = "***"`
			assertMatch(t, got, expected)
		})

		t.Run("should mask 'My_PassWord' by default", func(t *testing.T) {
			masker := NewTFMasker([]string{}, false)
			got := getMaskOutput(t, masker, ` + My_PassWord = "value"`)
			expected := ` + My_PassWord = "***"`
			assertMatch(t, got, expected)
		})

		t.Run("should mask custom property case sensitive", func(t *testing.T) {
			masker := NewTFMasker([]string{"my_prop"}, false)
			got := getMaskOutput(t, masker, ` + my_prop = "value"`)
			expected := ` + my_prop = "***"`
			assertMatch(t, got, expected)
		})

		t.Run("should mask custom property ignoring case", func(t *testing.T) {
			masker := NewTFMasker([]string{"my_prop"}, true)
			got := getMaskOutput(t, masker, ` + My_PrOP = "value"`)
			expected := ` + My_PrOP = "***"`
			assertMatch(t, got, expected)
		})

		t.Run("should not mask when property doesn't match", func(t *testing.T) {
			masker := NewTFMasker([]string{"my_prop"}, true)
			got := getMaskOutput(t, masker, ` + other_prop = "value"`)
			expected := ` + other_prop = "value"`
			assertMatch(t, got, expected)
		})
	})

	t.Run("multi-line new resource", func(t *testing.T) {
		t.Run("should mask property containing 'password' by default", func(t *testing.T) {
			masker := NewTFMasker([]string{}, false)
			input := `
  + resource "new_resource" {
  + prop        = "value"
  + my_password = "secret"
}`

			expected := `
  + resource "new_resource" {
  + prop        = "value"
  + my_password = "***"
}`
			got := getMaskOutput(t, masker, input)
			assertMatch(t, got, expected)
		})
	})
}
