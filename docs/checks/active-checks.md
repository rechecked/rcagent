# Active Checks

Active checks are checks where the request for data is sent from `check_rcagent.py` to the agent's Status API. The returned output, exit code, and perfdata is already in the Nagios plugins return format and the `check_rcagent.py` plugin just outputs what is returned, it does not do any data or check processing on the plugin side.

In general, active checks are configured on the monitoring system (Nagios Core) and each active check is ran as a full command with returned exit code, output, and perfdata, which is what the monitoring system uses to determine if things are OK or not.

## Downloading `check_rcagent.py`

If you don't have the plugin already, you can [download the latest version](https://github.com/rechecked/rcagent-plugins/releases/latest/download/check_rcagent.py).

There is a [GitHub rcagent-plugins repo](https://github.com/rechecked/rcagent-plugins) specifically for the plugins, so if you have any problems with the plugin, feel free to mention them in the issues.

## Installing `check_rcagent.py`

If you are using Nagios Core, you will need to copy the plugin into the plugins directory. Normally the directory is `/usr/local/nagios/libexec` but may be different on your specific system if you've customized the paths or installed via a repo.

## Using `check_rcagent.py`

Endpoints (such as CPU, Memory, Services, etc) are reached using `-e <endpoint>` while plugins are ran using `-p <plugin>`.

A typical call with the plugin will look similar to this:

```
./check_rcagent.py -H <host> -t <token> -e <endpoint> -w <warning> -c <critical>
```

Another useful feature is using `-q` to pass query arguments. This is necessary for things like services or processes when you need to send a name or value. You can add multiple `-q` values on the command line, like so:

```
./check_rcagent.py -H <host> -t <token> -e <endpoint> -q "name=test" -q "name2=test2"
```

!!! note

	There are more options than what we show here! You can use the `--help` option to see all available options.

### Endpoint Examples

The endpoint checks are available by default on any rcagent and does not require special setup, unlike plugins which have to be added after installation.

Here are a few examples of endpoint checks. If you want to see more endpoints you can use, check the Status API Reference section.

#### Disk Check

An example of a basic [disk](../../status-api/disk) check using the path value of `/` (root) and showing the return output. Note that the return output also includes any perfdata for the check. Some checks may not have perfdata associated with them.

Command:

```
./check_rcagent.py -H <host> -t <token> -e disk -q path=/
```

Output:

```
OK: Disk usage of / is 34.91% (12.23/35.04 GiB Total) | 'percent'=34.91% 'used'=12.23GiB 'free'=22.81GiB 'total'=35.04GiB
```

#### Service Check

When running against the [services](../../status-api/services) endpoint, you need to pass an `against` and `expected` parameter. The `against` parameter is the name of the service you want to run the check against. If the service status matches the `expected` value, it will be OK, otherwise it is CRITICAL.

```
./check_rcagent.py -H <host> -t <token> -e services -q "against=rcagent" -q "expected=running"
```

Output:

```
OK: rcagent is [running] (expected value is [running])
```

#### Memory Check

Most checks that return a number run with a `-w` (warning) and `-c` (critical) value. Memory checks are example of this. You can also adjust the units returned by sending in `-u`.

```
./check_rcagent.py -H <host> -t <token> -e memory/virtual -w 20 -c 60 -u GB

```

Output:

```
WARNING - Memory usage is 24.11% (16.49/68.38 GB Total) | 'percent'=24.11%;20;60 'available'=51.89GB;20;60 'used'=16.49GB;20;60 'free'=51.89GB;20;60 'total'=68.38GB;20;60
```

### Plugins Example

Plugins are added to the plugins directory. You specify the plugin you want to run using the `-p` flag rather than the `-e` endpoint flag.

When running a plugin use `--arg=""` for proper parsing rather than `-a` if you are using `-`/`--` in your arugment. An example of this is:

```
./check_rcagent.py -H <host> -t <token> -p check_test.sh --arg="--warning 10" --arg="-c 20"
```

Output:

```
--warning 10 -c 20
```

The above `check_test.sh` plugin just echos out whatever is passed to the plugin. So in the output we see the two arguments that we passed to the plugin.
