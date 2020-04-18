package mask

import (
	"fmt"
	"regexp"
)

// TFMasker reads from its reader, and masks lines matched by the matcher
type TFMasker struct {
	propsStr string

	newRemoveRegex         *regexp.Regexp
	replaceRegex           *regexp.Regexp
	replaceKnownAfterRegex *regexp.Regexp
	removeNullRegex        *regexp.Regexp

	newRemoveGroups         string
	replaceGroups           string
	replaceKnownAfterGroups string
	removeToNullGroups      string
}

// NewTFMasker creates a new masker using the specified reader and matcher
func NewTFMasker(props []string, ignoreCase bool) *TFMasker {
	masker := TFMasker{
		propsStr: getMaskedPropStr(props, ignoreCase),
	}

	masker.buildNewRemoveInfo()
	masker.buildRemoveToNullInfo()
	masker.buildReplaceInfo()
	masker.buildReplaceKnownAfterInfo()

	return &masker
}

// Mask scans the reader line by line and prints masked/unmasked output to the writer
func (m *TFMasker) Mask(config Config) {
	input := getInput(config.Reader)

	var output string
	output = m.newRemoveRegex.ReplaceAllString(input, m.newRemoveGroups)
	output = m.replaceRegex.ReplaceAllString(output, m.replaceGroups)
	output = m.replaceKnownAfterRegex.ReplaceAllString(output, m.replaceKnownAfterGroups)
	output = m.removeNullRegex.ReplaceAllString(output, m.replaceKnownAfterGroups)

	fmt.Fprint(config.Writer, output)
}

func (m *TFMasker) buildNewRemoveInfo() {
	newRemovePattern := fmt.Sprintf(
		`(?m)^( +[+-] +)(?P<prop>%s)( += +)(")(?P<value>[a-zA-Z0-9%%._-]+)(")$`,
		m.propsStr,
	)

	regex, groups := buildInfo(newRemovePattern, []string{"value"})

	m.newRemoveRegex = regex
	m.newRemoveGroups = groups
}

func (m *TFMasker) buildRemoveToNullInfo() {
	removeToNullPattern := fmt.Sprintf(
		`(?m)^( +- +)(?P<prop>%s)( += +)(")(?P<value>[a-zA-Z0-9%%._-]+)(")( +-> +)(null)$`,
		m.propsStr,
	)

	regex, groups := buildInfo(removeToNullPattern, []string{"value"})

	m.removeNullRegex = regex
	m.removeToNullGroups = groups
}

func (m *TFMasker) buildReplaceInfo() {
	replace := fmt.Sprintf(
		`(?m)^( +~ +)(?P<prop>%s)( += +)(")(?P<value>[a-zA-Z0-9%%._-]+)(")`+
			`( +-> +)(")(?P<changed_value>[a-zA-Z0-9%%._-]+)(")( +[#].*)*$`,
		m.propsStr,
	)

	regex, groups := buildInfo(replace, []string{"value", "changed_value"})

	m.replaceRegex = regex
	m.replaceGroups = groups
}

func (m *TFMasker) buildReplaceKnownAfterInfo() {
	replaceKnownAfterPattern := fmt.Sprintf(
		`(?m)^( +~ +)(?P<prop>%s)( += +)(")(?P<value>[a-zA-Z0-9%%._-]+)(")`+
			`( +-> +)(\(known after apply\))( +[#].*)*$`,
		m.propsStr,
	)

	regex, groups := buildInfo(replaceKnownAfterPattern, []string{"value"})

	m.replaceKnownAfterRegex = regex
	m.replaceKnownAfterGroups = groups
}
