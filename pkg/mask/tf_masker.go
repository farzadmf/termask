package mask

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	matcher "github.com/farzadmf/termask/pkg/match"
)

const (
	tfNewOrRemovedProp = "( +)( *?[+-] *?)( +)"
	tfChangedProp      = "( +)( *?[~] *?)( +)"
	tfRemovedProp      = "( +)( *?[-] *?)( +)"

	tfValue           = "([\"<])(?P<value>.*?)([>\"])"
	tfPropEquals      = "(?P<prop>[\"a-zA-Z0-9%._-]+)( +)(=)( +)"
	tfValueChange     = "( +)(->)( +)"
	tfComment         = "( +[#].*)*"
	tfNull            = "(null)"
	tfKnownAfterApply = "(\\(known after apply\\))"
)

var (
	tfNewOrRemoveStr            = fmt.Sprintf("^%s%s%s$", tfNewOrRemovedProp, tfPropEquals, tfValue)
	tfReplaceStr                = fmt.Sprintf("^%s%s%s%s%s%s$", tfChangedProp, tfPropEquals, tfValue, tfValueChange, tfValue, tfComment)
	tfReplaceKnownAfterApplyStr = fmt.Sprintf("^%s%s%s%s%s%s$", tfChangedProp, tfPropEquals, tfValue, tfValueChange, tfKnownAfterApply, tfComment)
	tfRemoveToNullStr           = fmt.Sprintf("^%s%s%s%s%s$", tfRemovedProp, tfPropEquals, tfValue, tfValueChange, tfNull)

	tfNewOrRemoveRegex            = regexp.MustCompile(tfNewOrRemoveStr)
	tfReplaceRegex                = regexp.MustCompile(tfReplaceStr)
	tfReplaceKnownAfterApplyRegex = regexp.MustCompile(tfReplaceKnownAfterApplyStr)
	tfRemoveToNullRegex           = regexp.MustCompile(tfRemoveToNullStr)
)

// TFMasker reads from its reader, and masks lines matched by the matcher
type TFMasker struct {
	matcher    matcher.Matcher
	propsRegex *regexp.Regexp
}

// NewTFMasker creates a new masker using the specified reader and matcher
func NewTFMasker(m matcher.Matcher, props []string, ignoreCase bool) TFMasker {
	return TFMasker{
		matcher:    m,
		propsRegex: regexp.MustCompile(getMaskedPropStr(props, ignoreCase)),
	}
}

// Mask scans the reader line by line and prints masked/unmasked output to the writer
func (m TFMasker) Mask(config Config) {
	scanner := bufio.NewScanner(config.Reader)
	for scanner.Scan() {
		// line := scanner.Text()

		// names := jsonLineRegex.SubexpNames()
		// matches := jsonLineRegex.FindAllStringSubmatch(line, -1)[0]

		// -------------------

		// match, matches := m.matcher.Match(line)

		// switch match {
		// case matcher.TFNewOrRemove:
		// 	fmt.Fprintln(config.Writer, m.maskNewOrRemove(matches))
		// case matcher.TFReplace:
		// 	fmt.Fprintln(config.Writer, m.maskReplace(matches))
		// case matcher.TFReplaceKnownAfterApply:
		// 	fmt.Fprintln(config.Writer, m.maskKnownAfterApply(matches))
		// case matcher.TFRemoveToNull:
		// 	fmt.Fprintln(config.Writer, m.maskRemoveToNull(matches))
		// case matcher.None:
		// 	fmt.Fprintln(config.Writer, line)
		// }
	}
}

// maskNewOrRemove masks a property value when a resource is being removed or added
func (m TFMasker) maskNewOrRemove(matches []string) string {
	leadingWhitespace := matches[1]
	plus := matches[2]
	spaceAfterPlus := matches[3]
	property := matches[4]
	spaceBeforeEqual := matches[5]
	spaceAfterEqual := matches[6]
	firstQuote := matches[7]
	value := matches[8]
	secondQuote := matches[9]

	if m.propsRegex.MatchString(property) {
		value = strings.Repeat("*", 3)
	}

	return fmt.Sprintf("%s%s%s%s%s=%s%s%s%s",
		leadingWhitespace, plus, spaceAfterPlus, property, spaceBeforeEqual,
		spaceAfterEqual, firstQuote, value, secondQuote,
	)
}

// maskReplace masks a property value when a resource is being replaced
func (m TFMasker) maskReplace(matches []string) string {
	leadingWhitespace := matches[1]
	plus := matches[2]
	spaceAfterPlus := matches[3]
	property := matches[4]
	spaceBeforeEqual := matches[5]
	spaceAfterEqual := matches[6]
	firstQuote := matches[7]
	value := matches[8]
	secondQuote := matches[9]
	spaceBeforeArrow := matches[10]
	spaceAfterArrow := matches[11]
	changeFirstQuote := matches[12]
	changeValue := matches[13]
	changeSecondQuote := matches[14]
	comment := matches[15]

	if m.propsRegex.MatchString(property) {
		value = strings.Repeat("*", 3)
		changeValue = strings.Repeat("*", 3)
	}

	return fmt.Sprintf("%s%s%s%s%s=%s%s%s%s%s->%s%s%s%s%s",
		leadingWhitespace, plus, spaceAfterPlus, property, spaceBeforeEqual,
		spaceAfterEqual, firstQuote, value, secondQuote, spaceBeforeArrow,
		spaceAfterArrow, changeFirstQuote, changeValue, changeSecondQuote, comment,
	)
}

// maskKnownAfterApply takes care of masking values when we have '... -> (known after apply)'
func (m TFMasker) maskKnownAfterApply(matches []string) string {
	leadingWhitespace := matches[1]
	plus := matches[2]
	spaceAfterPlus := matches[3]
	property := matches[4]
	spaceBeforeEqual := matches[5]
	spaceAfterEqual := matches[6]
	firstQuote := matches[7]
	value := matches[8]
	secondQuote := matches[9]
	spaceBeforeArrow := matches[10]
	spaceAfterArrow := matches[11]
	knownAfterApply := matches[12]
	comment := matches[13]

	if m.propsRegex.MatchString(property) {
		value = strings.Repeat("*", 3)
	}

	return fmt.Sprintf("%s%s%s%s%s=%s%s%s%s%s->%s%s%s",
		leadingWhitespace, plus, spaceAfterPlus, property, spaceBeforeEqual,
		spaceAfterEqual, firstQuote, value, secondQuote, spaceBeforeArrow,
		spaceAfterArrow, knownAfterApply, comment,
	)
}

// maskRemoveToNull masks values when a resource begin removed and we have '... -> null'
func (m TFMasker) maskRemoveToNull(matches []string) string {
	leadingWhitespace := matches[1]
	plus := matches[2]
	spaceAfterPlus := matches[3]
	property := matches[4]
	spaceBeforeEqual := matches[5]
	spaceAfterEqual := matches[6]
	firstQuote := matches[7]
	value := matches[8]
	secondQuote := matches[9]
	spaceBeforeArrow := matches[10]
	spaceAfterArrow := matches[11]
	null := matches[12]

	if m.propsRegex.MatchString(property) {
		value = strings.Repeat("*", 3)
	}

	return fmt.Sprintf("%s%s%s%s%s=%s%s%s%s%s->%s%s",
		leadingWhitespace, plus, spaceAfterPlus, property, spaceBeforeEqual,
		spaceAfterEqual, firstQuote, value, secondQuote, spaceBeforeArrow,
		spaceAfterArrow, null,
	)
}
