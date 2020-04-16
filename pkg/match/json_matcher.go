package match

import (
	"fmt"
	"regexp"
)

const (
	jsonLine = `(?P<leading_ws> *?)(?P<prop_begin_quote>")` +
		`(?P<prop>[a-zA-Z0-9%._-]+)(?P<prop_end_quote>")(?P<colon>:)(?P<ws_after_colon> )` +
		`(?P<value_begin_quote>")(?P<value>[a-zA-Z0-9%._-]+)(?P<value_end_quote>")(?P<trailing_chars>.*)`
)

var (
	jsonLineStr = fmt.Sprintf("^%s$", jsonLine)

	jsonLineRegex = regexp.MustCompile(jsonLineStr)
)

// JSONMatcher is used to match input lines against JSON
type JSONMatcher struct{}

// NewJSONMatcher creates a new matcher to match JSON lines
func NewJSONMatcher() JSONMatcher {
	return JSONMatcher{}
}

// Match match a line against JSON and returns the result
func (m JSONMatcher) Match(line string) (int, []string) {
	if jsonLineRegex.MatchString(line) {
		return JSONLine, jsonLineRegex.FindStringSubmatch(line)
	}

	return None, []string{}
}
