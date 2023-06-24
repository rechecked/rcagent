# Processes

## `processes`

Returns a list of processes and a total count. You can filter the list down by the process name. The check runs against the total count of running processes with `warning` and `critical`.

#### Options

Parameter | Default | Description
----------|---------|------------
`name` | | Optional value to specify what the process name is to filter down processes.
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.
`warning` | | Optional warning threshold value for checks.
`critical` | | Optional critical threshold value for checks.

#### Example

=== "URL"

	```
	https://localhost:5995/status/processes?token=private&pretty=1&name=cmd.exe
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/processes?token=private&pretty=1&check=1&warning=2&critical=5&name=cmd.exe
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e processes -q "name=cmd.exe" -w 2 -c 5
	```

#### Reponse

=== "URL"

	```
	{
		"count": 3,
		"processes": [
			{
				"name": "cmd.exe",
				"pid": 692,
				"exe": "C:\\Windows\\System32\\cmd.exe",
				"cmdline": "C:\\WINDOWS\\system32\\cmd.exe /d /s /c react-scripts start",
				"username": "DESKTOP-S0Q8L0D\\Jake",
				"cpuPercent": 0,
				"memPercent": 0.0081224599853158,
				"status": [
					""
				]
			},
			...
		]
	}
	```

=== "URL (Check)"

	```
	{
		"exitcode": 1,
		"output": "WARNING - Process count is 3",
		"perfdata": "'processes'=3;2;5",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	WARNING - Process count is 3 | 'processes'=3;2;5
	```
