package mask

import (
	"testing"
)

func TestMaskNewResource(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "should mask different combinations of 'password' by default",
			props:       []string{},
			ignoreCase:  false,
			input: `
+ resource "resource_type" "new_resource" {
  + prop        = "value"
  + my_password = "secret"
  + password    = "secret"
  + password2   = "secret"
  + PasSWorD    = "secret"
  + MyPassWORD2 = "SECRET"
}`,
			want: `
+ resource "resource_type" "new_resource" {
  + prop        = "value"
  + my_password = "***"
  + password    = "***"
  + password2   = "***"
  + PasSWorD    = "***"
  + MyPassWORD2 = "***"
}`,
		},
		{
			description: "masks specified property (case sensitive)",
			props:       []string{"my_prop"},
			ignoreCase:  false,
			input: `
+ resource "resource_type" "new_resource" {
  + my_prop     = "value"
  + My_Prop     = "value2"
}`,
			want: `
+ resource "resource_type" "new_resource" {
  + my_prop     = "***"
  + My_Prop     = "value2"
}`,
		},
		{
			description: "masks combinations of specified property when ignoring case",
			props:       []string{"my_prop"},
			ignoreCase:  true,
			input: `
+ resource "resource_type" "new_resource" {
  + my_prop     = "value"
  + My_Prop     = "value2"
}`,
			want: `
+ resource "resource_type" "new_resource" {
  + my_prop     = "***"
  + My_Prop     = "***"
}`,
		},
		{
			description: "prints output as is when no match",
			props:       []string{},
			ignoreCase:  false,
			input: `
+ resource "resource_type" "new_resource" {
  + my_prop     = "value"
}`,
			want: `
+ resource "resource_type" "new_resource" {
  + my_prop     = "value"
}`,
		},
	}

	runTFMaskerTests(t, cases)
}

func TestMaskChangedResource(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "only masks props containing 'password' by default",
			props:       []string{},
			ignoreCase:  false,
			input: `
~ resource "changed_resource" {
  ~ prop      = "value" -> "new_value"     # comment
  ~ password  = "secret" -> "new_secret"
  ~ password2 = "secret2" -> "new_secret2" # with comment
}`,
			want: `
~ resource "changed_resource" {
  ~ prop      = "value" -> "new_value"     # comment
  ~ password  = "***" -> "***"
  ~ password2 = "***" -> "***" # with comment
}`,
		},
		{
			description: "additionally masks specified prop",
			props:       []string{"prop"},
			ignoreCase:  false,
			input: `
~ resource "changed_resource" {
  ~ prop      = "value" -> "new_value"       # comment
  ~ prop2     = "value2" -> "new_value2"     # comment
  ~ password  = "secret" -> "new_secret"
  ~ password2 = "secret2" -> "new_secret2"   # with comment
}`,
			want: `
~ resource "changed_resource" {
  ~ prop      = "***" -> "***"       # comment
  ~ prop2     = "value2" -> "new_value2"     # comment
  ~ password  = "***" -> "***"
  ~ password2 = "***" -> "***"   # with comment
}`,
		},
	}

	runTFMaskerTests(t, cases)
}

func TestMaskRemovedResource(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "masks 'password' and specified prop",
			props:       []string{"prop"},
			ignoreCase:  false,
			input:       ` - password = "secret" -> null`,
			// - resource "removed_resource" {
			//   - prop      = "value" -> null
			//   - prop2     = "value2" -> null
			//   - password  = "secret" -> null
			//   - password2 = "secret2" -> null
			// }`,
			want: ` - password = "***" -> null`,
			// - resource "removed_resource" {
			//   - prop      = "***" -> null
			//   - prop2     = "value2" -> null
			//   - password  = "***" -> null
			//   - password2 = "***" -> null
			// }`,
		},
	}

	runTFMaskerTests(t, cases)
}

func TestMaskReplaceKnownAfterApply(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "masks 'password' and specified prop",
			props:       []string{"prop"},
			ignoreCase:  false,
			input: `
~ resource "known_after_resource" {
  ~ prop      = "value" -> (known after apply)       # comment
  ~ prop2     = "value2" -> (known after apply)     # comment
  ~ password  = "secret" -> (known after apply)
  ~ password2 = "secret2" -> (known after apply)   # with comment
}`,
			want: `
~ resource "known_after_resource" {
  ~ prop      = "***" -> (known after apply)       # comment
  ~ prop2     = "value2" -> (known after apply)     # comment
  ~ password  = "***" -> (known after apply)
  ~ password2 = "***" -> (known after apply)   # with comment
}`,
		},
	}

	runTFMaskerTests(t, cases)
}
