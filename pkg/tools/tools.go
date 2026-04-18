package tools

import "strconv"

func AtoiOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func ValueOrDefault[T comparable](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}
