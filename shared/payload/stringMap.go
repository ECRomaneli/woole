package payload

type StringMap map[string]*StringList

func CloneStringMap(strMap StringMap) StringMap {
	clone := make(StringMap)

	for key, strList := range strMap {
		clone[key] = &StringList{Val: append([]string{}, strList.Val...)}
	}

	return clone
}
