package main

import (
	"fmt"

	"github.com/gotify/plugin-api"
)

const messageTemplate = "Sender: %s\n\n```\n%s\n```"

func makeMarkdownMessage(title, message, remoteIP string) plugin.Message {
	return plugin.Message{
		Title: title,
		Message: fmt.Sprintf(messageTemplate,
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
