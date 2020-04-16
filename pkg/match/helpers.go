package match

func getValueIndex(groupNames []string) (index int) {
	index = -1
	for i := 0; i < len(groupNames); i++ {
		if groupNames[i] == "value" {
			index = i
		}
	}

	return
}
