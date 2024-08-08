package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gotify/plugin-api"
)

const (
	messageTemplate          = "Sender: %s\n\n%s"
	messageTemplateCodeBlock = "Sender: %s\n\n```\n%s\n```"
)

func makeMarkdownMessage(title, message, remoteIP string, withinCodeBlock bool) plugin.Message {
	tmpl := messageTemplate
	if withinCodeBlock {
		tmpl = messageTemplateCodeBlock
	}

	return plugin.Message{
		Title: title,
		Message: fmt.Sprintf(tmpl,
			remoteIP,
			message,
		),
		Extras: map[string]interface{}{
			"client::display": map[string]interface{}{
				"contentType": "text/markdown",
			},
		},
	}
}

// getContentTypeFromRequest tries to read the configured
// content-type from the request in the following order:
// "?content-type=" URL query parameter, "Content-Type" request header,
// "X-Content-Type" request header (for clients that can send custom
// headers, but cannot set the content type)
func getContentTypeFromRequest(req *http.Request) int {
	var (
		fromQuery        = req.URL.Query().Get("content-type")
		fromHeader       = req.Header.Get("content-type")
		fromNonStdHeader = req.Header.Get("x-content-type")
	)

	var foundType string
	if fromNonStdHeader != "" {
		foundType = fromNonStdHeader
	}
	if fromHeader != "" {
		foundType = fromHeader
	}
	if fromQuery != "" {
		foundType = fromQuery
	}

	switch strings.ToLower(foundType) {
	case "application/json":
		return ContentTypeJSON
	case "text/markdown":
		return ContentTypeMarkdown
	default:
		return ContentTypeUnknown
	}
}
