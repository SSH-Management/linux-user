package user

func mapSystemGroups(groups []string, m map[string]string) []string {
	newGroups := make([]string, 0, len(groups))

	for _, group := range groups {
		if replace, ok := m[group]; ok {
			newGroups = append(newGroups, replace)
		} else {
			newGroups = append(newGroups, group)
		}
	}

	return newGroups
}
