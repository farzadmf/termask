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
			description: "masks 'pasword' by default when there is no plus sign",
			input:       `password = "secret"`,
			want:        `password = "***"`,
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
		{
			description: "masks Azure connection string using partial match and ignore case",
			props:       []string{"connectionstring"},
			ignoreCase:  true,
			partial:     true,
			input: `"My__StorageAccountConnectionString"          = "DefaultEndpointsProtocol=https;` +
				`AccountName=account;AccountKey=2S8/T4B4cquIjr6w==;EndpointSuffix=core.windows.net"`,
			want: `"My__StorageAccountConnectionString"          = "***"`,
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
		{
			description: "masks Azure connection string using partial match and ignore case",
			props:       []string{"connectionstring"},
			ignoreCase:  true,
			partial:     true,
			input: `~ "My__StorageAccountConnectionString"          = "DefaultEndpointsProtocol=https;` +
				`AccountName=account;AccountKey=2S8/T4B4cquIjr6w==;EndpointSuffix=core.windows.net" -> ` +
				`"Ac=ac;Key=2D5/==;End=windows.net"`,
			want: `~ "My__StorageAccountConnectionString"          = "***" -> "***"`,
		},
	}

	runTFMaskerTests(t, cases)
}

func TestMaskRemovedResource(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "masks 'password' by default",
			input:       ` - password = "secret" -> null`,
			want:        ` - password = "***" -> null`,
		},
		{
			description: "masks specified prop with more than one specified",
			props:       []string{"prop", "my_prop"},
			input:       ` + my_prop = "secret"`,
			want:        ` + my_prop = "***"`,
		},
	}

	runTFMaskerTests(t, cases)
}

func TestMaskReplaceKnownAfterApply(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "masks 'My_Password_2' by default",
			input:       ` ~ My_Password_2 = "secret" -> (known after apply)`,
			want:        ` ~ My_Password_2 = "***" -> (known after apply)`,
		},
		{
			description: "masks specified prop when matched partially",
			props:       []string{"prop"},
			partial:     true,
			input:       ` ~ prop_value = "secret" -> (known after apply)`,
			want:        ` ~ prop_value = "***" -> (known after apply)`,
		},
	}

	runTFMaskerTests(t, cases)
}
