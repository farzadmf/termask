package main

import (
	"bytes"
	"strings"
	"testing"
)

var (
	output bytes.Buffer
	input  = ""

	config = MaskConfig{
		writer: &output,
	}
)

func TestMask(t *testing.T) {
	var masker Masker

	t.Run("should mask 'password' by default", func(t *testing.T) {
		masker = NewMasker(NewMatcher(), []string{}, false)
		input = ` + "password" = "value"`
		output = bytes.Buffer{}
		config.reader = strings.NewReader(input)

		masker.Mask(config)
		outputString := strings.TrimSpace(output.String())
		if outputString != `+ "password" = "***"` {
			t.Error("'password' value was not masked")
		}
	})

	t.Run("should mask 'PaSSworD' by default", func(t *testing.T) {
		masker = NewMasker(NewMatcher(), []string{}, false)
		input = ` + "PaSSworD" = "value"`
		output = bytes.Buffer{}
		config.reader = strings.NewReader(input)

		masker.Mask(config)
		outputString := strings.TrimSpace(output.String())
		if outputString != `+ "PaSSworD" = "***"` {
			t.Error("'PaSSworD' value was not masked")
		}
	})

	t.Run("should mask 'My_PassWord' by default", func(t *testing.T) {
		masker = NewMasker(NewMatcher(), []string{}, false)
		input = ` + "My_PassWord" = "value"`
		output = bytes.Buffer{}
		config.reader = strings.NewReader(input)

		masker.Mask(config)
		outputString := strings.TrimSpace(output.String())
		if outputString != `+ "My_PassWord" = "***"` {
			t.Error("'My_PassWord' value was not masked")
		}
	})

	t.Run("should mask custom property case sensitive", func(t *testing.T) {
		props := []string{"my_prop"}
		masker = NewMasker(NewMatcher(), props, false)
		input = ` + "my_prop" = "value"`
		output = bytes.Buffer{}
		config.reader = strings.NewReader(input)

		masker.Mask(config)
		outputString := strings.TrimSpace(output.String())
		t.Logf("HELLO '%s'", outputString)
		if outputString != `+ "my_prop" = "***"` {
			t.Error("did not mask custom property")
		}
	})
}
