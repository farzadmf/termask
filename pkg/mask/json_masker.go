package mask

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

const (
	jsonLine = `(?m)^( *?)(")(?P<prop>PROPS)(")(:)( )(")(?P<value>[a-zA-Z0-9%._-]+)(")(.*)$`
)

// JSONMasker takes care of masking JSON input
type JSONMasker struct {
	jsonLineRegex *regexp.Regexp
}

// NewJSONMasker creates a new JSON masker
func NewJSONMasker(props []string, ignoreCase bool) JSONMasker {
	propsStr := getMaskedPropStr(props, ignoreCase)
	jsonLineStr := strings.Replace(jsonLine, "PROPS", propsStr, -1)

	return JSONMasker{
		jsonLineRegex: regexp.MustCompile(jsonLineStr),
	}
}

// Mask masks values for properties matching the specified list or props
func (m JSONMasker) Mask(config Config) {
	scanner := bufio.NewScanner(config.Reader)
	var input []string
	for scanner.Scan() {
		line := scanner.Text()
		input = append(input, line)
	}
	inputStr := strings.Join(input, "\n")

	names := m.jsonLineRegex.SubexpNames()
	valueIndex := getGroupIndex(names, "value")

	var replaceGroups []string
	for i := 1; i < len(names); i++ {
		if i == valueIndex {
			replaceGroups = append(replaceGroups, "***")
		} else {
			replaceGroups = append(replaceGroups, fmt.Sprintf("${%d}", i))
		}
	}

	inputStr = m.jsonLineRegex.ReplaceAllString(inputStr, strings.Join(replaceGroups, ""))
	fmt.Fprintln(config.Writer, inputStr)
}
