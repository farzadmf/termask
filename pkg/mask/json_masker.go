package mask

import (
	"fmt"
	"regexp"
)

// JSONMasker takes care of masking JSON input
type JSONMasker struct {
	lineRegex         *regexp.Regexp
	lineReplaceGroups string
	propsStr          string
}

// NewJSONMasker creates a new JSON masker
func NewJSONMasker(props []string, ignoreCase bool) *JSONMasker {
	masker := JSONMasker{
		propsStr: getMaskedPropStr(props, ignoreCase),
	}

	masker.buildLineInfo()

	return &masker
}

// Mask masks values for properties matching the specified list or props
func (m *JSONMasker) Mask(config Config) {
	input := getInput(config.Reader)

	output := m.lineRegex.ReplaceAllString(input, m.lineReplaceGroups)

	fmt.Fprint(config.Writer, output)
}

func (m *JSONMasker) buildLineInfo() {
	linePattern := `(?m)^( *?)(")(?P<prop>PROPS)(")(:)( )(")(?P<value>[a-zA-Z0-9%._-]+)(")(.*)$`

	regex, groups := buildInfo(linePattern, m.propsStr, "value")

	m.lineRegex = regex
	m.lineReplaceGroups = groups
}
