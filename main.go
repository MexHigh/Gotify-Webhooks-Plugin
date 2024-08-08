package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gotify/plugin-api"
)

const routeName = "webhook"

const (
	ContentTypeUnknown = iota
	ContentTypeJSON
	ContentTypeMarkdown
)

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		Name:        "Webhooks",
		Description: "Plugin that enables Gotify to receive generic webhooks",
		ModulePath:  "git.leon.wtf/leon/gotify-webhook-plugin",
		Author:      "Leon Schmidt <mail@leon-schmidt.dev>",
		Website:     "https://leon-schmidt.dev",
	}
}

// Plugin is plugin instance
type Plugin struct {
	userCtx    plugin.UserContext
	msgHandler plugin.MessageHandler
	basePath   string
}

// Enable implements plugin.Plugin
func (p *Plugin) Enable() error {
	return nil
}

// Disable implements plugin.Plugin
func (p *Plugin) Disable() error {
	return nil
}

const helpMessageTemplate = "Use this **webhook URL**: %s\n\n" +
	"You can set the content type of the payload via the `content-type` query parameter (e.g. `%s?content-type=application/json`) or the `content-type` or `x-content-type` request headers.\n\n" +
	"The following content types are supported: `application/json`, `text/markdown`"

// GetDisplay implements plugin.Displayer
func (p *Plugin) GetDisplay(location *url.URL) string {
	baseHost := ""
	if location != nil {
		baseHost = fmt.Sprintf("%s://%s", location.Scheme, location.Host)
	}
	webhookURL := baseHost + p.basePath + routeName
	return fmt.Sprintf(helpMessageTemplate, webhookURL, webhookURL)
}

// SetMessageHandler implements plugin.Messenger
func (p *Plugin) SetMessageHandler(h plugin.MessageHandler) {
	// invoced during initialization
	p.msgHandler = h
}

// RegisterWebhook implements plugin.Webhooker
func (p *Plugin) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	p.basePath = basePath

	webhookHandler := func(c *gin.Context) {
		// read body
		bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Error reading request body",
				err.Error(),
				c.ClientIP(),
				false,
			))
			return
		}

		contentType := getContentTypeFromRequest(c.Request)
		// if content type is unknown, try JSON anyway (some clients cannot set a content-type)
		if contentType == ContentTypeUnknown {
			var data interface{}
			err = json.Unmarshal(bytes, &data)
			if err == nil {
				contentType = ContentTypeJSON
			}
		}

		switch contentType {
		case ContentTypeJSON:
			// try to parse json to verify format
			var data interface{}
			err = json.Unmarshal(bytes, &data)
			if err != nil {
				p.msgHandler.SendMessage(makeMarkdownMessage(
					"Error parsing JSON message",
					err.Error(),
					c.ClientIP(),
					false,
				))
				return
			}
			// re-indent JSON
			jsonStr, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				p.msgHandler.SendMessage(makeMarkdownMessage(
					"Error re-marshalling payload",
					err.Error(),
					c.ClientIP(),
					true,
				))
				return
			}
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Recieved webhook with JSON",
				string(jsonStr),
				c.ClientIP(),
				true,
			))
		case ContentTypeMarkdown:
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Recieved webhook with Markdown",
				string(bytes),
				c.ClientIP(),
				false,
			))
		case ContentTypeUnknown:
			// just send the string
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Recieved webhook with unknown content type",
				string(bytes),
				c.ClientIP(),
				true,
			))
		}
	}

	mux.POST("/"+routeName, webhookHandler)
	mux.PUT("/"+routeName, webhookHandler)
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &Plugin{
		userCtx: ctx,
	}
}

func main() {
	panic("this should be built as go plugin")
}
