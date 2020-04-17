package mask

import (
	"fmt"
	"strings"
)

func getMaskedPropStr(props []string, ignoreCase bool) string {
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

	return masked
}

func getGroupIndex(groupNames []string, name string) (index int) {
	for i, group := range groupNames {
		if group == name {
			index = i
			break
		}
	}

	return
}
