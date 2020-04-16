package match

import (
	"fmt"
	"regexp"
)

const (
	tfNewOrRemovedProp = "( +)( *?[+-] *?)( +)"
	tfChangedProp      = "( +)( *?[~] *?)( +)"
	tfRemovedProp      = "( +)( *?[-] *?)( +)"

	tfValue           = "([\"<])(.*?)([>\"])"
	tfPropEquals      = "([\"a-zA-Z0-9%._-]+)( +)=( +)"
	tfValueChange     = "( +)->( +)"
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

// TFMatcher is used to match a terraform line against a pattern
type TFMatcher struct{}

// NewTFMatcher creates a matcher to match input lines against known patterns
func NewTFMatcher() TFMatcher {
	return TFMatcher{}
}

// Match tries to match a line against a pattern
// Returns what we matched against and the matches slice (if we have a match)
func (m TFMatcher) Match(line string) (int, []string) {
	if tfNewOrRemoveRegex.MatchString(line) {
		return TFNewOrRemove, tfNewOrRemoveRegex.FindStringSubmatch(line)
	}

	if tfReplaceRegex.MatchString(line) {
		return TFReplace, tfReplaceRegex.FindStringSubmatch(line)
	}

	if tfReplaceKnownAfterApplyRegex.MatchString(line) {
		return TFReplaceKnownAfterApply, tfReplaceKnownAfterApplyRegex.FindStringSubmatch(line)
	}

	if tfRemoveToNullRegex.MatchString(line) {
		return TFRemoveToNull, tfRemoveToNullRegex.FindStringSubmatch(line)
	}

	return None, []string{}
}
