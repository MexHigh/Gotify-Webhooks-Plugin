# Gotify Webhooks Plugin

### Installation

Just download the latest `.so` file for your architecture from the [package registry](https://git.leon.wtf/leon/gotify-webhooks-plugin/-/packages) or build it yourself with `make build` (required Go and Docker). This uses Gotify's build tools to build against the latest version. The `.so` files are compiled to `build/gotify-webhooks*.so`.

Then simply move the `.so` file to the Gotify plugin directory and restart Gotify.

### Usage

Activate the Plugin, then go to the plugin's details panel to retrieve the **Webhook URL**. You can `POST` and `PUT` payload to it.

The plugin tries to determine the content type of the message by the following means (in order of preference):
- URL query parameter `content-type`, e.g. `https://<host>/plugin/1/custom/<id>/webhook?content-type=application/json`
- Standard request header `content-type`
- Nonstandard request header `x-content-type`

The following content types are currently supported:
- `application/json`
- `text/markdown` _(see [here](https://gotify.net/docs/msgextras#clientdisplay) to see how Gotify handles Markdown)_

If none of those are set or the content type is unknown, the plugin tries to parse the content as JSON. If this is not 
successfull, the request body will be displayed as-is or a corresponding error is shown.

The parsed payload is sent to the automatically created "Webhooks" application channel along with the senders IP address.
