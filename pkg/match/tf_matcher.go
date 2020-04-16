package match

import (
	"fmt"
	"regexp"
)

const (
	newOrRemovedProp = "( +)( *?[+-] *?)( +)"
	changedProp      = "( +)( *?[~] *?)( +)"
	removedProp      = "( +)( *?[-] *?)( +)"

	value           = "([\"<])(.*?)([>\"])"
	propEquals      = "([\"a-zA-Z0-9%._-]+)( +)=( +)"
	valueChange     = "( +)->( +)"
	comment         = "( +[#].*)*"
	null            = "(null)"
	knownAfterApply = "(\\(known after apply\\))"
)

var (
	newOrRemoveStr            = fmt.Sprintf("^%s%s%s$", newOrRemovedProp, propEquals, value)
	replaceStr                = fmt.Sprintf("^%s%s%s%s%s%s$", changedProp, propEquals, value, valueChange, value, comment)
	replaceKnownAfterApplyStr = fmt.Sprintf("^%s%s%s%s%s%s$", changedProp, propEquals, value, valueChange, knownAfterApply, comment)
	removeToNullStr           = fmt.Sprintf("^%s%s%s%s%s$", removedProp, propEquals, value, valueChange, null)

	newOrRemoveRegex            = regexp.MustCompile(newOrRemoveStr)
	replaceRegex                = regexp.MustCompile(replaceStr)
	replaceKnownAfterApplyRegex = regexp.MustCompile(replaceKnownAfterApplyStr)
	removeToNullRegex           = regexp.MustCompile(removeToNullStr)
)

// TFMatcher is used to match a terraform line against a pattern
type TFMatcher struct{}

// NewTFMatcher creates a matcher to match input lines against known patterns
func NewTFMatcher() TFMatcher {
	return TFMatcher{}
}

// Match tries to match a line against a pattern
// Returns what we matched against and the matches slice (if we have a match)
func (m TFMatcher) Match(line string) (int, []string) {
	if newOrRemoveRegex.MatchString(line) {
		return TFNewOrRemove, newOrRemoveRegex.FindStringSubmatch(line)
	}

	if replaceRegex.MatchString(line) {
		return TFReplace, replaceRegex.FindStringSubmatch(line)
	}

	if replaceKnownAfterApplyRegex.MatchString(line) {
		return TFReplaceKnownAfterApply, replaceKnownAfterApplyRegex.FindStringSubmatch(line)
	}

	if removeToNullRegex.MatchString(line) {
		return TFRemoveToNull, removeToNullRegex.FindStringSubmatch(line)
	}

	return None, []string{}
}
