[![Go Test/Build](https://github.com/rechecked/rcagent/actions/workflows/go.yml/badge.svg)](https://github.com/rechecked/rcagent/actions/workflows/go.yml)
[![GitHub release](https://img.shields.io/github/release/rechecked/rcagent?include_prereleases=&sort=semver&color=blue)](https://github.com/rechecked/rcagent/releases/)
[![License](https://img.shields.io/badge/License-GPLv3-blue)](https://github.com/rechecked/rcagent/blob/main/LICENSE)
[![rcagent - Documentation](https://img.shields.io/badge/rcagent-Documentation-informational)](https://rechecked.io/docs/rcagent)

ReChecked Agent (rcagent) is light-weight, cross-platform, API-based agent that is compatible with Nagios (using NRDP and check_rcagent) check-based monitoring systems.

### For Binary Installs

We currently build for: Windows, macOS (12+), CentOS Stream/RHEL (8+), Debian (10+), Ubuntu (18+ LTS)

Follow the [installation guide in the rcagent documentation](https://rechecked.io/docs/rcagent/getting-started/installation/).

Download from [GitHub releases](https://github.com/rechecked/rcagent/releases).

### For Source Install

To build from source, ensure `make` and `golang` are installed and download the souce, then run:

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

Read more about how to use rcagent in the [documentation](https://rechecked.io/docs/rcagent/).

- **Nagios XI** - Download and install the latest [ReChecked Nagios XI Config Wizard](https://github.com/rechecked/rcagent-nagiosxi/releases/latest/download/rcagent.zip) in the XI GUI.
- **Nagios Core** - Download the [check_rcagent.py](https://github.com/rechecked/rcagent-plugins/releases/latest/download/check_rcagent.py) file, place it in `/usr/local/nagios/libexec`.

### Config Options

For a full list of config options in `config.yml`, check the [config options documentation](https://rechecked.io/docs/rcagent/config/options/).

### Building Documentation

If you'd like the build your own set of docs (you'll need to host them in order for them to work properly) you can run the following:

Install mkdocs-material and run the mkdocs build:

```
pip install mkdocs-material
mkdocs build
```