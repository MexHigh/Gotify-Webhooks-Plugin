# Gotify Webhooks Plugin

### Installation

Build with `make build` (required Go and Docker). This uses Gotify's build tools to build against the latest version.

Then move the file matching your server's architecture (`build/gotify-webhook*.so`) to the Gotify plugin directory. Restart Gotify.

### Usage

Activate the Plugin, then go to the plugin's details panel to retrieve the **Webhook URL**. You can `POST` and `PUT` payload to it.

The plugin tries to parse the payload as JSON to create a more readable indentation. If it fails to do so, the payload is sent as-is.

The payloads are sent to the automatically created "Webhooks" application channel along with the senders IP address.