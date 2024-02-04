# Disk

## `disk`

Returns the amount of disk usage for a particular disk `path`.

#### Options

Parameter | Default | Description
----------|---------|------------
`path` | | Use that path to define the disk, drive, or mounted location. Examples: `/`, `C:`
`units`  | [*defaultUnits*](../../config/options#defaultunits) | Sets the units that will be returned. Available units: `kB`, `KiB`, `MB`, `MiB`, `GB`, `GiB`, `TB`, `TiB`. Default is set by [`defaultUnits`](../../config/options#defaultunits) in the [`config.yml`](../../config/options).
`pretty` | *false* | Set to `1` or `true` to format the JSON returned using a pretty print function.
`check`  | *false* | To run the endpoint and get check results, set to `1` or `true`.
`warning` | | Optional warning threshold value for checks.
`critical` | | Optional critical threshold value for checks.

#### Example

=== "URL"

	```
	https://localhost:5995/status/disk?token=private&pretty=1&warning=70&critical=90&path=C:
	```

=== "URL (Check)"

	```
	https://localhost:5995/status/disk?token=private&pretty=1&check=1&warning=70&critical=90&path=C:
	```

=== "Plugin"

	```
	./check_rcagent.py -H localhost -t private -e disk -q "path=C:" -w 70 -c 90

	```

#### Response 

=== "URL"

	```
	{
		"path": "C:",
		"device": "C:",
		"fstype": "NTFS",
		"total": 3725.2324180603027,
		"free": 1911.6025123596191,
		"used": 1813.6299057006836,
		"usedPercent": 48.68501350165489,
		"units": "GiB"
	}
	```

=== "URL (Check)"

	```
	{
		"exitcode": 0,
		"output": "OK - Disk usage of C: is 48.69% (1813.63/3725.23 GiB Total)",
		"perfdata": "'percent'=48.69%;70;90 'used'=1813.63GiB;70;90 'free'=1911.60GiB;70;90 'total'=3725.23GiB;70;90",
		"longoutput": ""
	}
	```

=== "Plugin"

	```
	OK - Disk usage of C: is 48.69% (1813.63/3725.23 GiB Total) | 'percent'=48.69%;70;90 'used'=1813.63GiB;70;90 'free'=1911.60GiB;70;90 'total'=3725.23GiB;70;90
	```
