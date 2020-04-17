package mask

import (
	"bytes"
	"strings"
	"testing"
)

func TestMask(t *testing.T) {
	var masker TFMasker
	var output bytes.Buffer

	input := ""
	config := Config{
		Writer: &output,
	}

	t.Run("should mask 'password' by default", func(t *testing.T) {
		masker = NewTFMasker([]string{}, false)
		input = ` + password = "value"`
		output = bytes.Buffer{}
		config.Reader = strings.NewReader(input)

		masker.Mask(config)
		trimmedOutput := strings.TrimSpace(output.String())
		if trimmedOutput != `+ password = "***"` {
			t.Errorf("'password' value was not masked; got '%s'", trimmedOutput)
		}
	})

	t.Run("should mask 'PaSSworD' by default", func(t *testing.T) {
		masker = NewTFMasker([]string{}, false)
		input = ` + PaSSworD = "value"`
		output = bytes.Buffer{}
		config.Reader = strings.NewReader(input)

		masker.Mask(config)
		trimmedOutput := strings.TrimSpace(output.String())
		if trimmedOutput != `+ PaSSworD = "***"` {
			t.Errorf("'PaSSworD' value was not masked; got '%s'", trimmedOutput)
		}
	})

	t.Run("should mask 'My_PassWord' by default", func(t *testing.T) {
		masker = NewTFMasker([]string{}, false)
		input = ` + My_PassWord = "value"`
		output = bytes.Buffer{}
		config.Reader = strings.NewReader(input)

		masker.Mask(config)
		trimmedOutput := strings.TrimSpace(output.String())
		if trimmedOutput != `+ My_PassWord = "***"` {
			t.Errorf("'My_PassWord' value was not masked; got '%s'", trimmedOutput)
		}
	})

	t.Run("should mask custom property case sensitive", func(t *testing.T) {
		props := []string{"my_prop"}
		masker = NewTFMasker(props, false)
		input = ` + my_prop = "value"`
		output = bytes.Buffer{}
		config.Reader = strings.NewReader(input)

		masker.Mask(config)
		trimmedOutput := strings.TrimSpace(output.String())
		if trimmedOutput != `+ my_prop = "***"` {
			t.Errorf("did not mask custom property; got '%s'", trimmedOutput)
		}
	})

	t.Run("should mask custom property ignoring case", func(t *testing.T) {
		props := []string{"my_prop"}
		masker = NewTFMasker(props, true)
		input = ` + My_PrOP = "value"`
		output = bytes.Buffer{}
		config.Reader = strings.NewReader(input)

		masker.Mask(config)
		trimmedOutput := strings.TrimSpace(output.String())
		if trimmedOutput != `+ My_PrOP = "***"` {
			t.Errorf("did not mask custom property, case insensitive; got '%s'", trimmedOutput)
		}
	})

	t.Run("should not mask when property doesn't match", func(t *testing.T) {
		props := []string{"my_prop"}
		masker = NewTFMasker(props, true)
		input = ` + other_prop = "value"`
		output = bytes.Buffer{}
		config.Reader = strings.NewReader(input)

		masker.Mask(config)
		trimmedOutput := strings.TrimSpace(output.String())
		if trimmedOutput != `+ other_prop = "value"` {
			t.Errorf("did not print property as is when no match; got '%s'", trimmedOutput)
		}
	})
}
