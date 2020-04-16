package match

func getValueIndex(groupNames []string) (propIndex, valueIndex int) {
	propIndex = -1
	valueIndex = -1

	for i := 0; i < len(groupNames); i++ {
		if groupNames[i] == "value" {
			valueIndex = i
		} else if groupNames[i] == "prop" {
			propIndex = i
		}
	}

	return
}
