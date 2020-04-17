package mask

import (
	"bytes"
	"strings"
	"testing"
)

func getMaskOutputTrimmed(t *testing.T, masker Masker, input string) string {
	t.Helper()

	output := bytes.Buffer{}
	config := Config{
		Writer: &output,
		Reader: strings.NewReader(input),
	}

	masker.Mask(config)
	return strings.TrimSpace(output.String())
}
