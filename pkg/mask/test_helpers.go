package mask

import (
	"bytes"
	"strings"
	"testing"
)

type maskerTestCase struct {
	description string
	props       []string
	ignoreCase  bool
	partial     bool
	input       string
	want        string
}

func runTFMaskerTests(t *testing.T, cases []maskerTestCase) {
	t.Helper()

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			t.Helper()

			masker := NewTFMasker(test.props, test.ignoreCase, test.partial)
			got := getMaskOutput(t, masker, test.input)
			assertMatch(t, got, test.want)
		})
	}
}

func runJSONMaskerTests(t *testing.T, cases []maskerTestCase) {
	t.Helper()

	for _, test := range cases {
		t.Run(test.description, func(t *testing.T) {
			t.Helper()

			masker := NewJSONMasker(test.props, test.ignoreCase, test.partial)
			got := getMaskOutput(t, masker, test.input)
			assertMatch(t, got, test.want)
		})
	}
}

func getMaskOutput(t *testing.T, masker Masker, input string) string {
	t.Helper()

	output := bytes.Buffer{}
	config := Config{
		Writer: &output,
		Reader: strings.NewReader(input),
	}

	masker.Mask(config)
	return output.String()
}

func assertMatch(t *testing.T, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}
}
