# CPU

## `cpu/percent`

Returns the current total CPU percentage.

#### Options

Parameter | Default | Description
----------|---------|------------
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.
`warning` | | Optional warning threshold value for checks.
`critical` | | Optional critical threshold value for checks.

#### Example

=== "URL"

	```
	https://localhost:5995/status/cpu/percent?token=private&pretty=1
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/cpu/percent?token=private&pretty=1&check=1&warning=70&critical=90
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e cpu/percent -w 70 -c 90

	```

#### Response 

=== "URL"

	```
	{
		"percent": [
			0.17331022530329288
		],
		"units": "%"
	}
	```

=== "URL (Check)"

	```
	{
		"exitcode": 0,
		"output": "OK - CPU usage is 0.52%",
		"perfdata": "'percent'=0.52%;70;90",
		"longoutput": "Intel(R) Core(TM) i9-10980XE CPU @ 3.00GHz [Total Cores: 36]"
	}
	```

=== "Plugin"

	```
	OK - CPU usage is 0.26% | 'percent'=0.26%;70;90
	Intel(R) Core(TM) i9-10980XE CPU @ 3.00GHz [Total Cores: 36]
	```