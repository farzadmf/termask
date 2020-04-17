package mask

import (
	"bytes"
	"strings"
	"testing"
)

func TestJSONMask(t *testing.T) {
	t.Run("should mask 'password' by default", func(t *testing.T) {
		masker := NewJSONMasker([]string{}, false)
		output := getMaskOutputTrimmed(t, masker, `  "password": "secret"`)
		if output != `"password": "***"` {
			t.Errorf("'password' value not masked; got '%s'", output)
		}
	})

	t.Run("should not mask an unspecified property", func(t *testing.T) {
		masker := NewJSONMasker([]string{}, false)
		output := getMaskOutputTrimmed(t, masker, `  "prop": "value"`)
		if output != `"prop": "value"` {
			t.Errorf("unspecified property 'prop' should not be masked; got '%s'", output)
		}
	})

	t.Run("should mask specified property", func(t *testing.T) {
		masker := NewJSONMasker([]string{"prop"}, false)
		output := getMaskOutputTrimmed(t, masker, `  "prop": "value"`)
		if output != `"prop": "***"` {
			t.Errorf("property 'prop' should be masked; got '%s'", output)
		}
	})
}

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
