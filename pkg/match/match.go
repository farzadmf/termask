package match

// Matcher is used to match a line against a pattern
type Matcher interface {
	Match(line string) (groupNames, matches []string)
}

// These values tell us what we matched against
const (
	None = iota

	TFNewOrRemove
	TFReplace
	TFReplaceKnownAfterApply
	TFRemoveToNull

	JSONLine
)
