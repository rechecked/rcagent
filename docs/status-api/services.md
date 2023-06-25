# Services

## `services`

Will return the current status of the service, generally this is `running` or `stopped`. On linux systems there may also be other statuses available to check against.

#### Options

Parameter | Default | Description
----------|---------|------------
`against` | | The name of the service to check against.
`expected` | | The expected value to check for. If the value matches, it returns `OK`, if it doesn't, it returns `CRITICAL`.
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.

#### Example

=== "URL"

	```
	https://localhost:5995/status/services?token=private&pretty=1
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/services?token=private&pretty=1&expected=running&against=rcagent&check=1
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e services -q "against=rcagent" -q "expected=running"
	```

#### Reponse

=== "URL"

	```
	[
		{
			"name": "AdobeARMservice",
			"status": "running"
		},
		{
			"name": "AdobeUpdateService",
			"status": "running"
		},
		...
	]
	```

=== "URL (Check)"

	```
	{
		"exitcode": 2,
		"output": "CRITICAL - rcagent is [stopped] (expected value is [running])",
		"perfdata": "",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	CRITICAL - rcagent is [stopped] (expected value is [running])
	```
