# Load

!!! note

	This endpoint is not available on Windows systems.

## `load`

Returns the current load values for the system. You can check against different load values or the highest load value using the `against` parameter.

#### Options

Parameter | Default | Description
----------|---------|------------
`against` | *load1* | What to run the warning and critical thresholds against. Options: `load1`, `load5`, `load15`, or `hightest`. 
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.
`warning` | | Optional warning threshold value for checks.
`critical` | | Optional critical threshold value for checks.

#### Example

=== "URL"

	```
	https://localhost:5995/status/load?token=private&pretty=1
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/load?token=private&pretty=1&check=1&against=highest&warning=10&critical=20
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e load -q "against=highest" -w 10 -c 20

	```

#### Response 

=== "URL"

	```
	{
		"load1": 1.37,
		"load5": 0.9,
		"load15": 0.49
	}
	```

=== "URL (Check)"

	```
	{
		"exitcode": 0,
		"output": "OK - Load average is 0.35, 0.67, 0.45",
		"perfdata": "'load1'=0.35;10;20 'load5'=0.67;10;20 'load15'=0.45;10;20",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	OK - Load average is 0.54, 0.74, 0.46 | 'load1'=0.54;10;20 'load5'=0.74;10;20 'load15'=0.46;10;20
	```
