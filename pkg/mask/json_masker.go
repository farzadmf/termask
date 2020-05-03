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
func NewJSONMasker(props []string, ignoreCase bool, partial bool) *JSONMasker {
	masker := JSONMasker{
		propsStr: getMaskedPropStr(props, ignoreCase, partial),
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
	linePattern := fmt.Sprintf(
		`(?m)^( *?)(")(?P<prop>%s)(")(:)( )(")(?P<value>[a-zA-Z0-9%%._-]+)(")(.*)$`,
		m.propsStr,
	)

	regex, groups := buildInfo(linePattern, []string{"value"})

	m.lineRegex = regex
	m.lineReplaceGroups = groups
}
