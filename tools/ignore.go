package tools

import "strings"

func ShouldIgnoreDir(path string) bool {
	ignoreDirs := []string{
		"node_modules",
		".git",
		"vendor",
		"dist",
		"static",
		"templates",
	}
	for _, dir := range ignoreDirs {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}
