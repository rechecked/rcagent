# Active Checks

Active checks are checks where the request for data is sent from `check_rcagent.py` to the agent's Status API. The returned output, exit code, and perfdata is already in the Nagios plugins return format and the `check_rcagent.py` plugin just outputs what is returned, it does not do any data or check processing on the plugin side.

## Downloading `check_rcagent.py`

If you don't have the plugin already, you can [download the latest version](https://github.com/rechecked/rcagent-plugins/releases/latest/download/check_rcagent.py).

There is a [GitHub rcagent-plugins repo](https://github.com/rechecked/rcagent-plugins) specifically for the plugins, so if you have any problems with the plugin, feel free to mention them in the issues.

## Using `check_rcagent.py`

