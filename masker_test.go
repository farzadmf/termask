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

	t.Run("should mask password by default", func(t *testing.T) {
		masker = NewMasker(NewMatcher(), []string{}, false)
		input = ` + "my_password" = "password"`
		config.reader = strings.NewReader(input)

		masker.Mask(config)
		t.Log("***** output", output.String())
	})
}
