package utils

import (
	"strings"
)

func FuzzyMatch(target, query string) bool {
	if query == "" {
		return true
	}

	target = strings.ToLower(target)
	query = strings.ToLower(query)

	targetIdx := 0
	queryIdx := 0

	for queryIdx < len(query) && targetIdx < len(target) {
		if query[queryIdx] == target[targetIdx] {
			queryIdx++
		}
		targetIdx++
	}

	return queryIdx == len(query)
}
