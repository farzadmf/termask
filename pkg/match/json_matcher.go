package match

import (
	"fmt"
	"regexp"
)

const (
	jsonLine = `( *?)(")(?P<prop>[a-zA-Z0-9%._-]+)(")(:)( )(")(?P<value>[a-zA-Z0-9%._-]+)(")(.*)`
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
func (m JSONMatcher) Match(line string) (propIndex, valueIndex int, matches []string) {
	propIndex = -1
	valueIndex = -1

	if jsonLineRegex.MatchString(line) {
		valueIndex = getValueIndex(jsonLineRegex.SubexpNames())
		matches = jsonLineRegex.FindAllStringSubmatch(line, -1)[0]
	}

	return
}
