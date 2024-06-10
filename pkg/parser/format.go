package parser

import (
	"strings"
)

// DetermineFileFormat determines the file format based on its extension
func DetermineFileFormat(filename string) (format string) {
	if strings.HasSuffix(filename, ".xml") {
		return "rss"
	} else if strings.HasSuffix(filename, ".json") {
		return "json"
	} else if strings.HasSuffix(filename, ".html") {
		return "html"
	}
	return "unknown"
}
