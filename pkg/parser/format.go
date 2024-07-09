package parser

import (
	"strings"
)

const (
	rssFormat     = "rss"
	jsonFormat    = "json"
	htmlFormat    = "html"
	unknownFormat = "unknown"
)

// DetermineFileFormat determines the file format based on its extension
func DetermineFileFormat(filename string) (format string) {
	if strings.HasSuffix(filename, ".xml") {
		return rssFormat
	} else if strings.HasSuffix(filename, ".json") {
		return jsonFormat
	} else if strings.HasSuffix(filename, ".html") {
		return htmlFormat
	}
	return unknownFormat
}
