# System

## `system`

The system endpoint returns information about the system itself. You can't use this endpoint for a check, but it may be useful to get system related info.

#### Options

Parameter | Default | Description
----------|---------|------------
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.

#### Example

=== "URL"

	```
	https://localhost:5995/status/system?token=private&pretty=1
	```

#### Response

=== "URL"

	```
	{
		"hostname": "DESKTOP-S0Q8L0D",
		"uptime": 148601,
		"bootTime": 1687249918,
		"procs": 270,
		"os": "windows",
		"platform": "Microsoft Windows 11 Pro",
		"platformFamily": "Standalone Workstation",
		"platformVersion": "10.0.22621.1848 Build 22621.1848",
		"kernelVersion": "10.0.22621.1848 Build 22621.1848",
		"kernelArch": "x86_64",
		"virtualizationSystem": "",
		"virtualizationRole": "",
		"hostId": "12b0bdae-2d4e-4e32-8d5b-7af4f3521ea2"
	}
	```

## `system/users`

Return the current users. For the check, it will show the current number of users. The `warning` and `critical` thresholds are against the total number of users. 

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
	https://localhost:5995/status/system/users?token=private&pretty=1
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/system/users?token=private&pretty=1&check=1&warning=5&critical=10
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e system/users -w 5 -c 10
	```

#### Response

=== "URL"
	
	```
	[
		{
			"username": "jake",
			"domain": "DESKTOP-S0Q8L0D",
			"isLocal": true,
			"isAdmin": true,
			"logonType": 2,
			"logonTime": "2023-06-20T03:33:04.6377448-05:00",
			"dnsDomainName": ""
		}
	]
	```

=== "URL (Check)"

	```
	{
		"exitcode": 0,
		"output": "OK - Current users count is 1",
		"perfdata": "'users'=1;5;10",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	OK - Current users count is 1 | 'users'=1;5;10
	```

## `system/version`

Get the version of the rcagent. Currently there is no threshold so it will always return OK.

#### Options

Parameter | Default | Description
----------|---------|------------
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.

#### Example

=== "URL"

	```
	https://localhost:5995/status/system/version?token=private&pretty=1
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/system/version?token=private&pretty=1&check=1
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e system/version
	```

#### Response

=== "URL"
	
	```
	{
		"version": "1.1.0"
	}
	```

=== "URL (Check)"

	```
	{
		"exitcode": 0,
		"output": "OK - rcagent version is 1.1.0",
		"perfdata": "",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	OK - rcagent version is 1.1.0
	```
