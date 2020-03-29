package main

import (
	"fmt"
	"regexp"
)

// These values tell us what we matched against
const (
	None = iota
	NewOrRemove
	Replace
	ReplaceKnownAfterApply
	RemoveToNull
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

// Matcher is used to match a line against a pattern
type Matcher struct{}

// NewMatcher creates a matcher to match input lines against known patterns
func NewMatcher() Matcher {
	return Matcher{}
}

// Match tries to match a line against a pattern
// Returns what we matched against and the matches slice (if we have a match)
func (m Matcher) Match(line string) (int, []string) {
	if newOrRemoveRegex.MatchString(line) {
		return NewOrRemove, newOrRemoveRegex.FindStringSubmatch(line)
	}

	if replaceRegex.MatchString(line) {
		return Replace, replaceRegex.FindStringSubmatch(line)
	}

	if replaceKnownAfterApplyRegex.MatchString(line) {
		return ReplaceKnownAfterApply, replaceKnownAfterApplyRegex.FindStringSubmatch(line)
	}

	if removeToNullRegex.MatchString(line) {
		return RemoveToNull, removeToNullRegex.FindStringSubmatch(line)
	}

	return None, []string{}
}
