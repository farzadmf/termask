package mask

import (
	"testing"
)

func TestMaskNewResource(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "masks 'password' by default",
			input:       ` + password = "secret"`,
			want:        ` + password = "***"`,
		},
		{
			description: "masks 'PaSSWorD' by default",
			input:       `+ PaSSWorD = "secret"`,
			want:        `+ PaSSWorD = "***"`,
		},
		{
			description: "masks 'password2' by default",
			input:       ` + password2 = "secret"`,
			want:        ` + password2 = "***"`,
		},
		{
			description: "masks 'myPassWORD' by default",
			input:       ` + myPassWORD = "secret"`,
			want:        ` + myPassWORD = "***"`,
		},
		{
			description: "masks specified property",
			props:       []string{"my_prop"},
			input:       ` + my_prop = "secret"`,
			want:        ` + my_prop = "***"`,
		},
		{
			description: "does not mask specified property when casing does not match",
			props:       []string{"my_prop"},
			input:       ` + My_Prop = "value"`,
			want:        ` + My_Prop = "value"`,
		},
		{
			description: "masks property ignoring case",
			props:       []string{"my_prop"},
			ignoreCase:  true,
			input:       ` + My_Prop = "value"`,
			want:        ` + My_Prop = "***"`,
		},
		{
			description: "does not mask when nothing matches",
			input:       ` + my_prop = "value"`,
			want:        ` + my_prop = "value"`,
		},
		{
			description: "masks new prop with quotes",
			props:       []string{"my_prop"},
			input:       `   + "my_prop" = "value"`,
			want:        `   + "my_prop" = "***"`,
		},
		{
			// This is the case where this prop is a sub-section of a parent section being added
			description: "masks new prop with quotes and no '+'",
			props:       []string{"my_prop"},
			input:       `  "my_prop" = "value"`,
			want:        `  "my_prop" = "***"`,
		},
		{
			description: "masks partially matched property",
			props:       []string{"other_prop", "my_prop"},
			partial:     true,
			input:       `   + "my_prop_partial" = "value"`,
			want:        `   + "my_prop_partial" = "***"`,
		},
	}

	runTFMaskerTests(t, cases)
}

func TestMaskChangedResource(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "masks 'password' with comment",
			input:       ` ~ password = "secret" -> "new_secret"   # comment`,
			want:        ` ~ password = "***" -> "***"   # comment`,
		},
		{
			description: "masks 'MyPasWorD2' with quotes and without comment",
			input:       ` ~ "MyPassWorD2" = "secret" -> "new_secret"`,
			want:        ` ~ "MyPassWorD2" = "***" -> "***"`,
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
			want:        ` - password = "***" -> null`,
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
