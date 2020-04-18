package mask

import (
	"bytes"
	"strings"
	"testing"
)

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

func assertMatch(t *testing.T, got, expected string) {
	t.Helper()

	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}
