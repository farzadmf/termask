package mask

import (
	"bufio"
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
func NewTFMasker(props []string, ignoreCase bool, partial bool) *TFMasker {
	masker := TFMasker{
		propsStr: getMaskedPropStr(props, ignoreCase, partial),
	}

	masker.buildNewInfo()
	masker.buildRemoveToNullInfo()
	masker.buildReplaceInfo()
	masker.buildReplaceKnownAfterInfo()

	return &masker
}

// Mask scans the reader line by line and prints masked/unmasked output to the writer
func (m *TFMasker) Mask(config Config) {
	scanner := bufio.NewScanner(config.Reader)
	for scanner.Scan() {
		line := scanner.Text()
		output := line

		if m.newRemoveRegex.MatchString(line) {
			output = m.newRemoveRegex.ReplaceAllString(line, m.newRemoveGroups)
		}

		if m.replaceRegex.MatchString(line) {
			output = m.replaceRegex.ReplaceAllString(output, m.replaceGroups)
		}

		if m.replaceKnownAfterRegex.MatchString(line) {
			output = m.replaceKnownAfterRegex.ReplaceAllString(output, m.replaceKnownAfterGroups)
		}

		if m.removeNullRegex.MatchString(line) {
			output = m.removeNullRegex.ReplaceAllString(output, m.replaceKnownAfterGroups)
		}

		fmt.Fprintln(config.Writer, output)
	}
}

func (m *TFMasker) buildNewInfo() {
	newPattern := fmt.Sprintf(
		`^( *[+-]? *)(?P<prop>"?%s"?)( += +)(")(?P<value>%s+)(")$`,
		m.propsStr,
		valuePattern,
	)

	regex, groups := buildRegexAndGroups(newPattern, []string{"value"})

	m.newRemoveRegex = regex
	m.newRemoveGroups = groups
}

func (m *TFMasker) buildRemoveToNullInfo() {
	removeToNullPattern := fmt.Sprintf(
		`^( *-? *)(?P<prop>"?%s"?)( += +)(")(?P<value>%s+)(")( +-> +)(null)$`,
		m.propsStr,
		valuePattern,
	)

	regex, groups := buildRegexAndGroups(removeToNullPattern, []string{"value"})

	m.removeNullRegex = regex
	m.removeToNullGroups = groups
}

func (m *TFMasker) buildReplaceInfo() {
	replace := fmt.Sprintf(
		`^( *~ *)(?P<prop>"?%s"?)( += +)(")(?P<value>%s+)(")( +-> +)(")(?P<changed_value>%s+)(")( +[#].*)*$`,
		m.propsStr,
		valuePattern,
		valuePattern,
	)

	regex, groups := buildRegexAndGroups(replace, []string{"value", "changed_value"})

	m.replaceRegex = regex
	m.replaceGroups = groups
}

func (m *TFMasker) buildReplaceKnownAfterInfo() {
	replaceKnownAfterPattern := fmt.Sprintf(
		`^( *~ *)(?P<prop>"?%s"?)( += +)(")(?P<value>%s+)(")( +-> +)(\(known after apply\))( +[#].*)*$`,
		m.propsStr,
		valuePattern,
	)

	regex, groups := buildRegexAndGroups(replaceKnownAfterPattern, []string{"value"})

	m.replaceKnownAfterRegex = regex
	m.replaceKnownAfterGroups = groups
}
