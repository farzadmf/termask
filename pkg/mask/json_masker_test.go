package mask

import (
	"testing"
)

func TestJSONMask(t *testing.T) {
	cases := []maskerTestCase{
		{
			description: "masks combinations of 'password' by default",
			props:       []string{},
			ignoreCase:  false,
			input: `
{
  "password": "secret",
  "myPassWOrd": "secret",
  "mPassWORD2": "secret"
}`,
			want: `
{
  "password": "***",
  "myPassWOrd": "***",
  "mPassWORD2": "***"
}`,
		},
		{
			description: "prints input as is when no match",
			props:       []string{},
			ignoreCase:  false,
			input: `
{
  "prop": "value",
  "prop2": "value",
}`,
			want: `
{
  "prop": "value",
  "prop2": "value",
}`,
		},
		{
			description: "masks 'password' and specified property (case sensitive)",
			props:       []string{"prop"},
			ignoreCase:  false,
			input: `
{
  "prop": "value",
  "PROP": "value",
  "password": "secret"
}`,
			want: `
{
  "prop": "***",
  "PROP": "value",
  "password": "***"
}`,
		},
		{
			description: "masks combinations of specified property when ignoring case",
			props:       []string{"prop"},
			ignoreCase:  true,
			input: `
{
  "prop": "value",
  "PrOP": "value",
  "prop2": "value"
}`,
			want: `
{
  "prop": "***",
  "PrOP": "***",
  "prop2": "value"
}`,
		},
	}

	runJSONMaskerTests(t, cases)
}
