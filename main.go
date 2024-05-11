package main

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gotify/plugin-api"
)

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		Name:        "Webhook Plugin",
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
func (c *Plugin) Enable() error {
	return nil
}

// Disable implements plugin.Plugin
func (c *Plugin) Disable() error {
	return nil
}

// GetDisplay implements plugin.Displayer
func (c *Plugin) GetDisplay(location *url.URL) string {
	if c.userCtx.Admin {
		return "You are an admin! You have super cow powers."
	} else {
		return "You are **NOT** an admin! You can do nothing:("
	}
}

// SetMessageHandler implements plugin.Messenger
func (c *Plugin) SetMessageHandler(h plugin.MessageHandler) {
	// invoced during initialization
	c.msgHandler = h
}

// RegisterWebhook implements plugin.Webhooker
func (c *Plugin) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	c.basePath = basePath
	mux.POST("/webhook", func(con *gin.Context) {
		// TODO
		c.msgHandler.SendMessage(plugin.Message{
			Title:   "Recieved webhook!",
			Message: "Bimbim bambam",
		})
	})
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
