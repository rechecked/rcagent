# Checks

Checks are passive checks configured to run at certain intervals while rcagent is running. These checks are sent to some other location (like NRDP) using [senders](../options#senders).

In the YAML file, `checks` is a list of checks to run. You can have as many passive check configurations as you want.

## Special Values

### `$LOCAL_HOSTNAME`

This value is populated with the hostname of the system the rcagent is running on.

### `$HOST_ADDRESS`

The value of the host's `address` field. This variable can only be used n the plugin `args` but is available for the host and the services.

## Config Options

Options with a * next to them are **required**.

### `hostname` *

The hostname associated with the passive check.

### `address`

For hosts, you can define an address variable which sets `$HOST_ADDRESS` which can be used in plugin `args`.

### `servicename`

The service name (or service description, if a service check) associated with the passive check.

### `interval` *

The interval in which to run the check. It can be in format: `Xs` (seconds), `Xm` (minutes), `Xh` (hours) where `X` is a number.

### `endpoint ` *

The endpoint to use for the check, just like an active check. Example is `memory/virtual` or `services`.

### `options`

You can pass all the normal URL-style parameters in the options, such as warning/critical value and more. See the example in the [`config.yml`](../options) file for formatting.

## Example Check Config

```
checks:
  - hostname: $LOCAL_HOSTNAME
    interval: 5m
    endpoint: system/version
    options:
      warning: 10
      critical: 20
  - hostname: $LOCAL_HOSTNAME
    servicename: CPU Usage
    interval: 30s
    endpoint: cpu/percent
    options:
      warning: 10
      critical: 20
```