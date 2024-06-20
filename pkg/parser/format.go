package parser

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
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

func DetermineFeedFormat(urlPath url.URL) (format string, err error) {
	resp, err := http.Head(urlPath.String())
	if err != nil {
		return unknownFormat, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Errorf("error occurred while closing response body: %s", err.Error())
		}
	}(resp.Body)

	contentType := resp.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "text/xml"):
		return rssFormat, nil
	case strings.Contains(contentType, "application/json"):
		return jsonFormat, nil
	case strings.Contains(contentType, "text/html"):
		return htmlFormat, nil
	default:
		return unknownFormat, nil
	}
}
