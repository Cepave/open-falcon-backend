package utils

func UniqueArrayOfStrings(arrayOfStrings []string) []string {
	mapOfUnique := make(map[string]bool)

	for _, strValue := range arrayOfStrings {
		mapOfUnique[strValue] = true
	}

	result := make([]string, 0, len(mapOfUnique))
	for k := range mapOfUnique {
		result = append(result, k)
	}

	return result
}
