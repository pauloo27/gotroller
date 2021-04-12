package utils

import "strconv"

func EnforceSize(text string, maxLen int) string {
	if maxLen <= 0 || len(text) <= maxLen {
		return text
	}

	return text[0:maxLen-3] + "..."
}

func AtoiOrDefault(str string, defaultValue int) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return i
}
