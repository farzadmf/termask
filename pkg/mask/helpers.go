package mask

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func buildRegexAndGroups(pattern string, maskedNames []string) (regex *regexp.Regexp, replaceGroups string) {
	regex = regexp.MustCompile(pattern)

	names := regex.SubexpNames()
	var indices []int
	for _, name := range maskedNames {
		indices = append(indices, getGroupIndex(names, name))
	}

	replaceGroups = buildReplaceGroups(names, indices...)

	return
}

func getInput(reader io.Reader) string {
	scanner := bufio.NewScanner(reader)
	var input []string
	for scanner.Scan() {
		line := scanner.Text()
		input = append(input, line)
	}
	return strings.Join(input, "\n")
}

func getMaskedPropStr(props []string, ignoreCase bool, partial bool) string {
	masked := `(?i)"?.*?password.*?"?`

	for _, prop := range props {
		var caseString string
		if ignoreCase {
			caseString = "(?i)"
		} else {
			caseString = "(?-i)"
		}

		if partial {
			masked = fmt.Sprintf(`%s|"?%s.*%s.*"?`, masked, caseString, prop)
		} else {
			masked = fmt.Sprintf(`%s|"?%s%s"?`, masked, caseString, prop)
		}
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

func buildReplaceGroups(names []string, maskedIndices ...int) string {
	var replaceGroups []string
	for i := 1; i < len(names); i++ {
		appended := false
		for _, index := range maskedIndices {
			if i == index {
				replaceGroups = append(replaceGroups, "***")
				appended = true
				break
			}
		}
		if !appended {
			replaceGroups = append(replaceGroups, fmt.Sprintf("${%d}", i))
		}
	}

	return strings.Join(replaceGroups, "")
}
