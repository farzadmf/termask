package mask

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

const (
	jsonLine = `( *?)(")(?P<prop>[a-zA-Z0-9%._-]+)(")(:)( )(")(?P<value>[a-zA-Z0-9%._-]+)(")(.*)`
)

var (
	jsonLineStr = fmt.Sprintf("^%s$", jsonLine)

	jsonLineRegex = regexp.MustCompile(jsonLineStr)
)

// JSONMasker takes care of masking JSON input
type JSONMasker struct {
	propsRegex *regexp.Regexp
}

// NewJSONMasker creates a new JSON masker
func NewJSONMasker(props []string, ignoreCase bool) JSONMasker {
	return JSONMasker{
		propsRegex: regexp.MustCompile(getMaskedPropStr(props, ignoreCase)),
	}
}

// Mask masks values for properties matching the specified list or props
func (m JSONMasker) Mask(config Config) {
	scanner := bufio.NewScanner(config.Reader)
	for scanner.Scan() {
		line := scanner.Text()

		names := jsonLineRegex.SubexpNames()

		if jsonLineRegex.MatchString(line) {
			matches := jsonLineRegex.FindAllStringSubmatch(line, -1)[0]
			propIndex := getGroupIndex(names, "prop")
			valueIndex := getGroupIndex(names, "value")

			if m.propsRegex.MatchString(matches[propIndex]) {
				matches[valueIndex] = "***"
			}

			fmt.Fprintln(config.Writer, strings.Join(matches[1:], ""))
		} else {
			fmt.Fprintln(config.Writer, line)
		}
	}
}
