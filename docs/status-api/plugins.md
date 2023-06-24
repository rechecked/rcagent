# Plugins

The plugins endpoint is the most versatile part of the agent. Just like with Nagios Core, you can use any plugins, including scripts, as long as they adhere to the [Nagios Plugin Guidelines](https://nagios-plugins.org/doc/guidelines.html) and return text output on stdout, an exit code, and optionally some formatted performance data.

## Adding Plugins

By default, rcagent does not come with any plugins. You can add the plugins into the plugin directory after install. You also may need to update the config file if the plugin needs extra arguments in the run command.

### Plugins Directory

The default plugins directories. You can also set your [`pluginDir`](../../config/options#plugindir) in the config file.

=== "Linux"

	`/usr/lib64/rcagent/plugins` OR `/usr/lib/rcagent/plugins`

=== "Windows"

	`C:\Program Files\rcagent\plugins`

=== "macOS"

	`/usr/lib64/rcagent/plugins` OR `/usr/lib/rcagent/plugins`

### Plugin Types

By default `.sh`, `.py`, `.pl`, `.php` plugin types are supported with the default config. Plugins that are not definied in the [`pluginTypes`](../../config/options#plugintypes) section of the config will be ran as:

```
<plugin_name> <plugin_args>
```

So if you need to pass extra options to make the plugins work, you'll need to edit your config and add those specific plugin types.

## `plugins`

The plugins endpoint allows you view the plugins installed, and to run plugins and pass arguments to the plugin. You can pass multiple arguments with multiple `arg` values in the parameters or `--arg` values in the plugin.

If you wanted to pass an argument such as `-m "hello"` you'd pass it as `arg=-m "hello"` in the URL parameters or as `--arg="-p 9950"` as a plugin argument.

#### Options

Parameter | Default | Description
----------|---------|------------
`plugin` | | The name of the plugin you want to run.
`arg` | | The arguments you want to supply to the plugin. You can pass multiple `arg` values, check the URL (Check) or Plugin example to see how.
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.


#### Example

=== "URL"

	```
	https://localhost:5995/status/plugins?token=private&pretty=1
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/plugins?token=private&pretty=1&plugin=check_test.ps1&arg=test&arg=test2
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -p check_test.ps1 --arg="test" --arg="test2"

	```

#### Response 

=== "URL"

	```
	{
		"plugins": [
			"check_test.ps1",
			"check_test.sh"
		]
	}
	```

=== "URL (Check)"

	```
	{
		"output": "Testing plugin running and output\nArg1: test\nArg2: test2",
		"exitcode": 0
	}
	```

=== "Plugin"

	```
	Testing plugin running and output
	Arg1: test
	Arg2: test2
	```
