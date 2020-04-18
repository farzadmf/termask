package mask

import (
	"testing"
)

func TestMaskNewResource(t *testing.T) {
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

func TestMaskChangedResource(t *testing.T) {
	t.Run("should only mask props containing 'password' by default", func(t *testing.T) {
		masker := NewTFMasker([]string{}, false)
		input := `
~ resource "changed_resource" {
  ~ prop      = "value" -> "new_value"     # comment
  ~ password  = "secret" -> "new_secret"
  ~ password2 = "secret2" -> "new_secret2" # with comment
}`
		expected := `
~ resource "changed_resource" {
  ~ prop      = "value" -> "new_value"     # comment
  ~ password  = "***" -> "***"
  ~ password2 = "***" -> "***" # with comment
}`
		got := getMaskOutput(t, masker, input)
		assertMatch(t, got, expected)
	})

	t.Run("should additionally mask specified prop", func(t *testing.T) {
		masker := NewTFMasker([]string{"prop"}, false)
		input := `
~ resource "changed_resource" {
  ~ prop      = "value" -> "new_value"       # comment
  ~ prop2     = "value2" -> "new_value2"     # comment
  ~ password  = "secret" -> "new_secret"
  ~ password2 = "secret2" -> "new_secret2"   # with comment
}`
		expected := `
~ resource "changed_resource" {
  ~ prop      = "***" -> "***"       # comment
  ~ prop2     = "value2" -> "new_value2"     # comment
  ~ password  = "***" -> "***"
  ~ password2 = "***" -> "***"   # with comment
}`
		got := getMaskOutput(t, masker, input)
		assertMatch(t, got, expected)
	})
}

func TestMaskRemovedResource(t *testing.T) {
	t.Run("should mask 'password' and specified prop", func(t *testing.T) {
		masker := NewTFMasker([]string{"prop"}, false)
		input := `
- resource "removed_resource" {
  - prop      = "value" -> null
  - prop2     = "value2" -> null
  - password  = "secret" -> null
  - password2 = "secret2" -> null
}`
		expected := `
- resource "removed_resource" {
  - prop      = "***" -> null
  - prop2     = "value2" -> null
  - password  = "***" -> null
  - password2 = "***" -> null
}`
		got := getMaskOutput(t, masker, input)
		assertMatch(t, got, expected)
	})
}

func TestMaskReplaceKnownAfterApply(t *testing.T) {
	t.Run("should mask 'password' and specified prop", func(t *testing.T) {
		masker := NewTFMasker([]string{"prop"}, false)
		input := `
~ resource "known_after_resource" {
  ~ prop      = "value" -> (known after apply)       # comment
  ~ prop2     = "value2" -> (known after apply)     # comment
  ~ password  = "secret" -> (known after apply)
  ~ password2 = "secret2" -> (known after apply)   # with comment
}`
		expected := `
~ resource "known_after_resource" {
  ~ prop      = "***" -> (known after apply)       # comment
  ~ prop2     = "value2" -> (known after apply)     # comment
  ~ password  = "***" -> (known after apply)
  ~ password2 = "***" -> (known after apply)   # with comment
}`
		got := getMaskOutput(t, masker, input)
		assertMatch(t, got, expected)
	})
}
