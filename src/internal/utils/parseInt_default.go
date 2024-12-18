package utils

import "strconv"

func ParseIntDefault(input string, defaultValue int) int {
	parsedValue, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}
