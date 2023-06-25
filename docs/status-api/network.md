# Network

## `network`

Returns network usage on an interface. By default, the API endpoint will show the current counter value of bytes send/recieved. If you'd like to see a per second version, you can set the `delta` value to get the per second version. Delta is calculated as change per second since the last call.

Checks that are ran using the `check` command or through the plugin will always set `delta=1` so it will always return units per second when running checks.

#### Options

Parameter | Default | Description
----------|---------|------------
`against` | *total* | The value to check the warning and critical thresholds against. Options: `in`, `out`, `total`.
`delta` | *0* | Setting delta to `1` will return the delta values since the last call. This value will be set to `1` when `check` is set to `1` or `true`.
`units`  | [*defaultUnits*](../../config/options#defaultunits) | Sets the units that will be returned. Available units: `kB`, `KiB`, `MB`, `MiB`, `GB`, `GiB`, `TB`, `TiB`. Default is set by [`defaultUnits`](../../config/options#defaultunits) in the [`config.yml`](../../config/options).
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.
`warning` | | Optional warning threshold value for checks.
`critical` | | Optional critical threshold value for checks.

#### Counter Example

##### Example

=== "URL"

	```
	https://localhost:5995/status/network?token=private&pretty=1&name=Ethernet
	```

##### Response

=== "URL"

	```
	{
		"hardwareAddr": "2c:f0:5d:8a:96:6b",
		"addrs": [
			{
				"addr": "fe80::8b88:df2:9ba0:f362/64"
			},
			{
				"addr": "192.168.1.3/24"
			}
		],
		"name": "Ethernet",
		"bytesSent": 2872277064,
		"bytesRecv": 76592601758,
		"packetsSent": 16771499,
		"packetsRecv": 55809453,
		"errin": 0,
		"errout": 0,
		"dropin": 449377,
		"dropout": 0,
		"fifoin": 0,
		"fifoout": 0
	}
	```

#### Delta Example

This example shows using the delta value

##### Example

=== "URL"

	```
	https://localhost:5995/status/network?token=private&pretty=1&name=Ethernet&delta=1&units=kB
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/network?token=private&pretty=1&name=Ethernet&delta=1&units=kB&check=1&warning=100&critical=200
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e network -q "name=Ethernet" -w 100 -c 200 -u kB

	```

##### Response

=== "URL"

	```
	{
		"outTotal": 2376.871,
		"outPerSec": 30.64,
		"inTotal": 40.689,
		"inPerSec": 1789.866,
		"units": "kB"
	}
	```

=== "URL (Check)"

	```
	{
		"exitcode": 2,
		"output": "CRITICAL - Current network traffic 362.19 kB/s",
		"perfdata": "'in'=343.81kB/s;100;200 'out'=18.38kB/s;100;200",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	CRITICAL - Current network traffic 1268.70 kB/s | 'in'=1244.99kB/s;100;200 'out'=23.70kB/s;100;200
	```
