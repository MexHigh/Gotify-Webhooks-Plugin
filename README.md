# Gotify Webhooks Plugin

### Installation

Just download the latest `.so` file for your architecture from the [package registry](https://git.leon.wtf/leon/gotify-webhooks-plugin/-/packages) or build it yourself with `make build` (required Go and Docker). This uses Gotify's build tools to build against the latest version. The `.so` files are compiled to `build/gotify-webhooks*.so`.

Then simply move the `.so` file to the Gotify plugin directory and restart Gotify.

### Usage

Activate the Plugin, then go to the plugin's details panel to retrieve the **Webhook URL**. You can `POST` and `PUT` payload to it.

The plugin tries to parse the payload as JSON to create a more readable indentation. If it fails to do so, the payload is sent as-is.

The payloads are sent to the automatically created "Webhooks" application channel along with the senders IP address.