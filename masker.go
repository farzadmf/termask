package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// MaskConfig is used to specify the reader to mask and writer to write the output
type MaskConfig struct {
	reader io.Reader
	writer io.Writer
}

// Masker reads from its reader, and masks lines matched by the matcher
type Masker struct {
	MaskedProps []string

	matcher    Matcher
	propsRegex *regexp.Regexp
}

// NewMasker creates a new masker using the specified reader and matcher
func NewMasker(m Matcher, props []string, ignoreCase bool) Masker {
	masked := "((?i).*password)"

	if len(props) > 0 {
		var caseString string
		if ignoreCase {
			caseString = "(?i)"
		} else {
			caseString = "(?-i)"
		}
		masked = fmt.Sprintf("^(%s|%s%s)$", masked, caseString, strings.Join(props, "|"))
	}

	masker := Masker{
		matcher:    m,
		propsRegex: regexp.MustCompile(masked),
	}

	return masker
}

// Mask scans the reader line by line and prints masked/unmasked output to the writer
func (m Masker) Mask(config MaskConfig) {
	scanner := bufio.NewScanner(config.reader)
	for scanner.Scan() {
		line := scanner.Text()
		match, matches := m.matcher.Match(line)

		switch match {
		case NewOrRemove:
			fmt.Fprintln(config.writer, m.maskNewOrRemove(matches))
		case Replace:
			fmt.Fprintln(config.writer, m.maskReplace(matches))
		case ReplaceKnownAfterApply:
			fmt.Fprintln(config.writer, m.maskKnownAfterApply(matches))
		case RemoveToNull:
			fmt.Fprintln(config.writer, m.maskRemoveToNull(matches))
		case None:
			fmt.Fprintln(config.writer, line)
		}
	}
}

// maskNewOrRemove masks a property value when a resource is being removed or added
func (m Masker) maskNewOrRemove(matches []string) string {
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
func (m Masker) maskReplace(matches []string) string {
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
func (m Masker) maskKnownAfterApply(matches []string) string {
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
func (m Masker) maskRemoveToNull(matches []string) string {
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
