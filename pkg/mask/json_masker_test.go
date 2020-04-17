package mask

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONMask(t *testing.T) {
	var masker JSONMasker
	var output bytes.Buffer

	input := ""
	config := Config{
		Writer: &output,
	}

	t.Run("should mask 'password' by default", func(t *testing.T) {
		masker = NewJSONMasker([]string{}, false)
		input = `  "password": "secret"`
		output = bytes.Buffer{}
		config.Reader = strings.NewReader(input)

		masker.Mask(config)
		trimmedOutput := strings.TrimSpace(output.String())
		if trimmedOutput != `"password": "***"` {
			t.Errorf("'password' value not masked; got %q", trimmedOutput)
		}
	})
}
