package mask

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

const (
	tfNewOrRemovedProp = "( +[+-] +)"
	tfChangedProp      = "( +~ +)"
	tfRemovedProp      = "( +- +)"

	tfValue           = `(")(?P<value>[a-zA-Z0-9%._-]+)(")`
	tfChangedValue    = `(")(?P<changed_value>[a-zA-Z0-9%._-]+)(")`
	tfPropEquals      = "(?P<prop>[a-zA-Z0-9%._-]+)( += +)"
	tfValueChange     = "( +-> +)"
	tfComment         = "( +[#].*)*"
	tfNull            = "(null)"
	tfKnownAfterApply = "(\\(known after apply\\))"
)

var (
	tfNewOrRemoveStr            = fmt.Sprintf("^%s%s%s$", tfNewOrRemovedProp, tfPropEquals, tfValue)
	tfReplaceStr                = fmt.Sprintf("^%s%s%s%s%s%s$", tfChangedProp, tfPropEquals, tfValue, tfValueChange, tfChangedValue, tfComment)
	tfReplaceKnownAfterApplyStr = fmt.Sprintf("^%s%s%s%s%s%s$", tfChangedProp, tfPropEquals, tfValue, tfValueChange, tfKnownAfterApply, tfComment)
	tfRemoveToNullStr           = fmt.Sprintf("^%s%s%s%s%s$", tfRemovedProp, tfPropEquals, tfValue, tfValueChange, tfNull)

	tfNewOrRemoveRegex            = regexp.MustCompile(tfNewOrRemoveStr)
	tfReplaceRegex                = regexp.MustCompile(tfReplaceStr)
	tfReplaceKnownAfterApplyRegex = regexp.MustCompile(tfReplaceKnownAfterApplyStr)
	tfRemoveToNullRegex           = regexp.MustCompile(tfRemoveToNullStr)
)

// TFMasker reads from its reader, and masks lines matched by the matcher
type TFMasker struct {
	propsRegex *regexp.Regexp
}

// NewTFMasker creates a new masker using the specified reader and matcher
func NewTFMasker(props []string, ignoreCase bool) TFMasker {
	return TFMasker{
		propsRegex: regexp.MustCompile(getMaskedPropStr(props, ignoreCase)),
	}
}

// Mask scans the reader line by line and prints masked/unmasked output to the writer
func (m TFMasker) Mask(config Config) {
	scanner := bufio.NewScanner(config.Reader)
	for scanner.Scan() {
		line := scanner.Text()
		var matches []string

		if tfNewOrRemoveRegex.MatchString(line) {
			names := tfNewOrRemoveRegex.SubexpNames()
			matches = tfNewOrRemoveRegex.FindAllStringSubmatch(line, -1)[0]

			m.maskPropAndValue(matches, names)
		} else if tfReplaceRegex.MatchString(line) {
			names := tfReplaceRegex.SubexpNames()
			matches = tfReplaceRegex.FindAllStringSubmatch(line, -1)[0]

			m.maskPropAndValues(matches, names)
		} else if tfReplaceKnownAfterApplyRegex.MatchString(line) {
			names := tfReplaceKnownAfterApplyRegex.SubexpNames()
			matches = tfReplaceKnownAfterApplyRegex.FindAllStringSubmatch(line, -1)[0]

			m.maskPropAndValue(matches, names)
		} else if tfRemoveToNullRegex.MatchString(line) {
			names := tfRemoveToNullRegex.SubexpNames()
			matches = tfRemoveToNullRegex.FindAllStringSubmatch(line, -1)[0]

			m.maskPropAndValue(matches, names)
		}

		if len(matches) > 0 {
			line = strings.Join(matches[1:], "")
		}

		fmt.Fprintln(config.Writer, line)
	}
}

func (m TFMasker) maskPropAndValue(matches, names []string) {
	propIndex := getGroupIndex(names, "prop")
	valueIndex := getGroupIndex(names, "value")

	if m.propsRegex.MatchString(matches[propIndex]) {
		matches[valueIndex] = "***"
	}
}

func (m TFMasker) maskPropAndValues(matches, names []string) {
	propIndex := getGroupIndex(names, "prop")
	valueIndex := getGroupIndex(names, "value")
	changedValueIndex := getGroupIndex(names, "changedValue")

	if m.propsRegex.MatchString(matches[propIndex]) {
		matches[valueIndex] = "***"
		matches[changedValueIndex] = "***"
	}
}
