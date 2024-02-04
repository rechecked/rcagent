# Memory

## `memory/virtual`

Returns the memory usage of the system.

#### Options

Parameter | Default | Description
----------|---------|------------
`units`  | [*defaultUnits*](../../config/options#defaultunits) | Sets the units that will be returned. Available units: `kB`, `KiB`, `MB`, `MiB`, `GB`, `GiB`, `TB`, `TiB`. Default is set by [`defaultUnits`](../../config/options#defaultunits) in the [`config.yml`](../../config/options).
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.
`warning` | | Optional warning threshold value for checks.
`critical` | | Optional critical threshold value for checks.

#### Example

=== "URL"

	```
	https://localhost:5995/status/memory/virtual?token=private&pretty=1
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/memory/virtual?token=private&pretty=1&check=1
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e memory/virtual -w 10 -c 20

	```

#### Response 

=== "URL"

	```
	{
		"total": 63.684268951416016,
		"available": 48.12907409667969,
		"free": 48.12907409667969,
		"used": 15.555194854736328,
		"usedPercent": 24.425490173410335,
		"units": "GiB"
	}
	```

=== "URL (Check)"

	```
	{
		"exitcode": 2,
		"output": "CRITICAL - Memory usage is 24.14% (15.37/63.68 GiB Total)",
		"perfdata": "'percent'=24.14%;10;20 'available'=48.31GiB;10;20 'used'=15.37GiB;10;20 'free'=48.31GiB;10;20 'total'=63.68GiB;10;20",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	CRITICAL - Memory usage is 24.16% (15.38/63.68 GiB Total) | 'percent'=24.16%;10;20 'available'=48.30GiB;10;20 'used'=15.38GiB;10;20 'free'=48.30GiB;10;20 'total'=63.68GiB;10;20
	```

## `memory/swap`

Returns the total swap memory usage of the system.

#### Options

Parameter | Default | Description
----------|---------|------------
`units`  | [*defaultUnits*](../../config/options#defaultunits) | Sets the units that will be returned. Available units: `kB`, `KiB`, `MB`, `MiB`, `GB`, `GiB`, `TB`, `TiB`. Default is set by [`defaultUnits`](../../config/options#defaultunits) in the [`config.yml`](../../config/options).
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.
`warning` | | Optional warning threshold value for checks.
`critical` | | Optional critical threshold value for checks.

#### Example

=== "URL"

	```
	https://localhost:5995/status/memory/swap?token=private&pretty=1
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/memory/swap?token=private&pretty=1&check=1&warning=10&critical=20
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e memory/swap -w 10 -c 20

	```

#### Response 

=== "URL"

	```
	{
		"total": 67.68426895141602,
		"free": 48.68544387817383,
		"used": 18.998825073242188,
		"usedPercent": 28.069779533084127,
		"units": "GiB"
	}
	```

=== "URL (Check)"

	```
	{
		"exitcode": 2,
		"output": "CRITICAL - Swap usage is 28.03% (18.97/67.68 GiB Total)",
		"perfdata": "'percent'=28.03%;10;20 'used'=18.97GiB;10;20 'free'=48.71GiB;10;20 'total'=67.68GiB;10;20",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	CRITICAL - Swap usage is 28.14% (19.04/67.68 GiB Total) | 'percent'=28.14%;10;20 'used'=19.04GiB;10;20 'free'=48.64GiB;10;20 'total'=67.68GiB;10;20
	```