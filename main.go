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
	PayloadTypeUnknown = iota
	PayloadTypeJSON
)

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		Name:        "Webhooks",
		Description: "Plugin that enables Gotify to receive webhooks",
		ModulePath:  "git.leon.wtf/leon/gotify-webhook-plugin",
		Author:      "Leon Schmidt",
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

// GetDisplay implements plugin.Displayer
func (p *Plugin) GetDisplay(location *url.URL) string {
	baseHost := ""
	if location != nil {
		baseHost = fmt.Sprintf("%s://%s", location.Scheme, location.Host)
	}
	return fmt.Sprintf("Use this URL to recieve webhooks: %s%s%s", baseHost, p.basePath, routeName)
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
			))
			return
		}

		// try to parse json
		payloadType := PayloadTypeUnknown
		var data interface{}
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			payloadType = PayloadTypeUnknown
		} else {
			payloadType = PayloadTypeJSON
		}

		switch payloadType {
		case PayloadTypeJSON:
			// re-indent JSON
			jsonStr, err := json.MarshalIndent(data, "", "    ")
			if err != nil {
				p.msgHandler.SendMessage(makeMarkdownMessage(
					"Error re-marshalling payload",
					err.Error(),
					c.ClientIP(),
				))
				return
			}
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Recieved webhook",
				string(jsonStr),
				c.ClientIP(),
			))
		// TODO add more types?
		case PayloadTypeUnknown:
			// just send the string
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Recieved non-JSON webhook",
				string(bytes),
				c.ClientIP(),
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
