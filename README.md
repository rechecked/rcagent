[![Go Test/Build](https://github.com/rechecked/rcagent/actions/workflows/go.yml/badge.svg)](https://github.com/rechecked/rcagent/actions/workflows/go.yml)

ReChecked Agent (rcagent) is light-weight, cross-platform, API-based agent that is compatible with Nagios (using NRDP and check_rcagent) check-based monitoring systems.

### For Binary Installs

We currently build for: Windows, macOS (12+), CentOS Stream/RHEL (8+), Debian (10+), Ubuntu (18+ LTS)

We recommend [using the repo install](https://repo.rechecked.io/) for CentOS/RHEL, Debian, and Ubuntu.

Download from [GitHub releases](https://github.com/rechecked/rcagent/releases) or [download from rechecked.io](https://rechecked.io/download)

### For Source Install

To build from source, ensure golang is installed and download the souce, then run:

```
make build
```

To install the source version, run the following:

```
make install
```

If you'd like to run the source version as a service, you can install the service by running

```
/usr/local/rcagent/rcagent -a install
```

### Using rcagent

You can read how to [get started](https://rechecked.io/quick-start-guide/) and you can install the config wizard or plugin below.

- **Nagios XI** - Download and install the latest [ReChecked Nagios XI Config Wizard](https://rechecked.io/download) in the XI GUI.
- **Nagios Core** - Download the [check_rcagent.py](https://rechecked.io/download) file, place it in `/usr/local/nagios/libexec`.

### Config Options

For a full list of config options in config.yml, check the [config options page](https://rechecked.io/config-options/).
