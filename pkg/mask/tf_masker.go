package mask

import (
	"fmt"
	"regexp"
)

const (
// a tfNewOrRemovedProp = "( +[+-] +)"
// b tfChangedProp      = "( +~ +)"
// c tfRemovedProp      = "( +- +)"

// d tfValue        = `(")(?P<value>[a-zA-Z0-9%._-]+)(")`
// e tfChangedValue = `(")(?P<changed_value>[a-zA-Z0-9%._-]+)(")`
// tfPropEquals      = "(?P<prop>[a-zA-Z0-9%._-]+)( += +)"
// f tfPropEquals      = "(?P<prop>PROPS)( += +)"
// g tfValueChange     = "( +-> +)"
// h tfComment         = "( +[#].*)*"
// i tfNull            = "(null)"
// j tfKnownAfterApply = "(\\(known after apply\\))"
)

var (
	tfNewOrRemoveStr = `( +[+-] +)(?P<prop>PROPS)( += +)(")(?P<value>[a-zA-Z0-9%._-]+)(")`
	tfReplaceStr     = `( +~ +)(?P<prop>PROPS)( += +)(")(?P<value>[a-zA-Z0-9%._-]+)(")` +
		`( +-> +)(")(?P<changed_value>[a-zA-Z0-9%._-]+)(")( +[#].*)*`
	tfReplaceKnownAfterApplyStr = `( +~ +)(?P<prop>PROPS)( += +)(")(?P<value>[a-zA-Z0-9%._-]+)(")` +
		`( +-> +)(\\(known after apply\\))( +[#].*)*`
	tfRemoveToNullStr = `( +- +)(?P<prop>PROPS)( += +)(")(?P<value>[a-zA-Z0-9%._-]+)(")( +-> +)(null)`
)

// TFMasker reads from its reader, and masks lines matched by the matcher
type TFMasker struct {
	propsStr string
	// propsRegex                  *regexp.Regexp
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
		// propsRegex: regexp.MustCompile(getMaskedPropStr(props, ignoreCase)),
	}

	masker.buildNewRemoveInfo()
	masker.buildRemoveToNullInfo()
	masker.buildReplaceInfo()
	masker.buildReplaceKnownAfterInfo()

	return &masker
}

func (m *TFMasker) buildReplaceInfo() {
	replace := `( +~ +)(?P<prop>PROPS)( += +)(")(?P<value>[a-zA-Z0-9%._-]+)(")` +
		`( +-> +)(")(?P<changed_value>[a-zA-Z0-9%._-]+)(")( +[#].*)*`

	regex, groups := buildInfo(replace, m.propsStr, "value", "changed_value")

	m.replaceRegex = regex
	m.replaceGroups = groups
}

func (m *TFMasker) buildNewRemoveInfo() {
	newRemovePattern := `( +[+-] +)(?P<prop>PROPS)( += +)(")(?P<value>[a-zA-Z0-9%._-]+)(")`

	regex, groups := buildInfo(newRemovePattern, m.propsStr, "value")

	m.newRemoveRegex = regex
	m.newRemoveGroups = groups
}

func (m *TFMasker) buildReplaceKnownAfterInfo() {
	replaceKnownAfterPattern := `( +~ +)(?P<prop>PROPS)( += +)(")(?P<value>[a-zA-Z0-9%._-]+)(")` +
		`( +-> +)(\\(known after apply\\))( +[#].*)*`

	regex, groups := buildInfo(replaceKnownAfterPattern, m.propsStr, "value")

	m.replaceKnownAfterRegex = regex
	m.replaceKnownAfterGroups = groups
}

func (m *TFMasker) buildRemoveToNullInfo() {
	tfRemoveToNullStr = `( +- +)(?P<prop>PROPS)( += +)(")(?P<value>[a-zA-Z0-9%._-]+)(")( +-> +)(null)`

	regex, groups := buildInfo(tfNewOrRemoveStr, m.propsStr, "value")

	m.removeNullRegex = regex
	m.removeToNullGroups = groups
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

	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	var matches []string

	// 	if tfNewOrRemoveRegex.MatchString(line) {
	// 		names := tfNewOrRemoveRegex.SubexpNames()
	// 		matches = tfNewOrRemoveRegex.FindAllStringSubmatch(line, -1)[0]

	// 		m.maskPropAndValue(matches, names)
	// 	} else if tfReplaceRegex.MatchString(line) {
	// 		names := tfReplaceRegex.SubexpNames()
	// 		matches = tfReplaceRegex.FindAllStringSubmatch(line, -1)[0]

	// 		m.maskPropAndValues(matches, names)
	// 	} else if tfReplaceKnownAfterApplyRegex.MatchString(line) {
	// 		names := tfReplaceKnownAfterApplyRegex.SubexpNames()
	// 		matches = tfReplaceKnownAfterApplyRegex.FindAllStringSubmatch(line, -1)[0]

	// 		m.maskPropAndValue(matches, names)
	// 	} else if tfRemoveToNullRegex.MatchString(line) {
	// 		names := tfRemoveToNullRegex.SubexpNames()
	// 		matches = tfRemoveToNullRegex.FindAllStringSubmatch(line, -1)[0]

	// 		m.maskPropAndValue(matches, names)
	// 	}

	// 	if len(matches) > 0 {
	// 		line = strings.Join(matches[1:], "")
	// 	}

	// 	fmt.Fprintln(config.Writer, line)
	// }
}

// func (m TFMasker) getNewOrRemoveReplaceGroups() string {
// 	names := m.newOrRemoveRegex.SubexpNames()
// 	valueIndex := getGroupIndex(names, "value")

// 	maskedIndices := []int{}
// 	var replaceGroups []string
// 	for i := 1; i < len(names); i++ {
// 		appended := false
// 		for _, index := range maskedIndices {
// 			if i == index {
// 				replaceGroups = append(replaceGroups, "***")
// 				appended = true
// 				break
// 			}
// 		}
// 		if !appended {
// 			replaceGroups = append(replaceGroups, fmt.Sprintf("${%d}", i))
// 		}
// 	}
// }

// func (m TFMasker) maskPropAndValue(matches, names []string) {
// 	propIndex := getGroupIndex(names, "prop")
// 	valueIndex := getGroupIndex(names, "value")

// 	if m.propsRegex.MatchString(matches[propIndex]) {
// 		matches[valueIndex] = "***"
// 	}
// }

// func (m TFMasker) maskPropAndValues(matches, names []string) {
// 	propIndex := getGroupIndex(names, "prop")
// 	valueIndex := getGroupIndex(names, "value")
// 	changedValueIndex := getGroupIndex(names, "changedValue")

// 	if m.propsRegex.MatchString(matches[propIndex]) {
// 		matches[valueIndex] = "***"
// 		matches[changedValueIndex] = "***"
// 	}
// }
