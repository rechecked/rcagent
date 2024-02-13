# Passive Checks

Passive checks are a type of check that run on the agent side and the results are sent to the monitotring system. In our case, we normally send passive checks to Nagios via NRDP. You can use NRDP with both Nagios Core and Nagios XI which includes it pre-configured by default.

## Why Use Passive Checks?

One of the biggest reasons for using passive checks is that it will offload some of the resource requirements from the monitoring system. Active checks are ran through a subprocess by Nagios Core and will use various amounts of system resources depending on the kind of check.

Passive checks are more like a distributed system; running specific checks for the Nagios Core instance and returning just the output, result code, and perfdata back to Nagios Core. The resources required to read that data is much less than used to do the checks themselves.

!!! note

	There is one difference with passive checks compared to active checks that you should remember. You will need to set up [freshness checks](https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/4/en/freshness.html) for passive checks to ensure that if you are no longer recieving a check from rcagent, you know there is a problem you need to deal with!

## Setup NRDP Sender(s)

You can configure one or more [`senders`](../../config/options/#senders) to send passive checks. Currently the only available sender is for NRDP which comes installed by default on Nagios XI but can easily be added to Nagios Core if you follow the [NRDP readme on GitHub](https://github.com/NagiosEnterprises/nrdp).

If you have multiple senders, all checks performed will be sent to all the senders you have added.

Example sender config:

```
senders:
  - name: NRDP Server 1
    url: http://192.168.0.100/nrdp/
    token: sometoken
    type: nrdp
```

## Adding Passive Checks

It is possible to add as many passive checks as you want to your system. One of the nice things about passive checks is that you can run plugins and offload some of the work from your main monitoring system.

Each check can be set with a pecific `interval` time, which is per check. There are also other [`checks` config options available](../../config/checks).

### Example Host Check

The variable `$LOCAL_HOSTNAME` will be replaced with the hostname of the rcagent system

```
checks:
  - hostname: $LOCAL_HOSTNAME
    interval: 5m
    endpoint: system/version
    options:
      warning: 10
      critical: 20
```

### Example Service Checks

Example of service checks, including running a plugin as a passive service check:

```
checks:
  - hostname: $LOCAL_HOSTNAME
    servicename: Custom Plugin
    interval: 5m
    endpoint: plugins
    options:
      plugin: check_test.ps1
      args:
        - -m "hello and test!"
        - --dir /test/dir
  - hostname: $LOCAL_HOSTNAME
    servicename: CPU Usage
    interval: 30s
    endpoint: cpu/percent
    options:
      warning: 10
      critical: 20
  - hostname: $LOCAL_HOSTNAME
    servicename: Memory Usage
    interval: 5m
    endpoint: memory/virtual
    options:
      warning: 80
      critical: 90
  - hostname: $LOCAL_HOSTNAME
    servicename: Disk Usage - C:
    interval: 1h
    endpoint: disk
    options:
      path: C:
      warning: 70
      critical: 90
```