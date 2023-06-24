# Config Options

Configuration for rcagent is done through YAML config files. On install the `config.yml` file holds a default configuration. This section covers all the config options available.

## Main Options

### `address`

Address to bind the rcagent api. If left blank, it will bind to all available addresses. The default is set to blank.

**Default:**

### `port`

Port that the rcagent uses for the api.

**Default:** `5995`

### `token`

The security token for the rcagent Status API. This should be changed when the agent is installed and is required in the URL of any request to the Status API.

**Default:** `private`

### `tls`

This section defines the TLS certificates when using https (default) for the Status API.

By default the **adhoc** setting tells the rcagent to create a self-signed certificate when it starts if one does not already exist.

!!! note

	When using the ReChecked Manager, **adhoc** will mean that the rcagent will request a signed certificate from the organization it is connected to.

**Options:** **adhoc** or **[file path]**

**Default:**
```
tls:
  cert: adhoc
  key: adhoc
```

Example file path:
```
tls:
  cert: /path/to/my_cert.pem
  key: /path/to/my_key.key
```

### `defaultUnits`

Default units sets the units returned by the Status API when no units is passed during the API call.

**Default:** `GiB`

### `debug`

Turns debugging in the log files on.

**Default:** `false`

### `excludeFsTypes`

A list of file types to exclude for disk checks. There can be a lot of types that you do not wish to monitor and don't want to show up in the wizard or be returned by the Status API.

**Default:** `["aufs", "autofs", "binfmt_misc", "cgroup", "configfs",
  "debugfs", "devpts", "devtmpfs", "encryptfs", "efivarfs", "fuse", "fuseblk",
  "fusectl", "hugetlbfs","mqueue", "overlayfs", "proc", "pstore", "rpc_pipefs",
  "securityfs", "selinuxfs", "sysfs", "tmpfs", "tracefs", "bpf",
  "fuse.vmware-vmblock", "fuse.gvfsd-fuse"]`

## Plugin Options

### `runPluginsAsRoot`

If set to true, rcagent will run all plugins as root rather than the rcagent user.

!!! note
	
	On windows systems, all plugins are ran as the system user and this option is not used.

**Default:** `false`

### `pluginDir`

The plugin directory that plugins are stored in. The defaults below are automatically checked for plugins, so if no value is set, it'll look at all locations specified.

=== "Linux"

    **Default:** `/usr/lib64/rcagent/plugins` OR `/usr/lib/rcagent/plugins`

=== "Windows"

    **Default:** `C:\Program Files\rcagent\plugins`

=== "macOS"

    **Default:** `/usr/lib64/rcagent/plugins` OR `/usr/lib/rcagent/plugins`

### `pluginTypes`

In order to properly run plugins, we define a list of plugin types based on the file extension. Depending on the extension of the plugin, it will run a specific command for that plugin.

Commands are broken up into each individual token, so you need to define an array. In the default, the `-` symbol in the YAML file creates an array but you can alternatively do:
```
.ps1: ["powershell", "-ExecutionPolicy", "Bypass", "-File"]
```

As an example, if we were running `check_test.ps1 -w 20 -c 30` then the command rcagent would run would be:
```
powershell -ExecutionPolicy Bypass -File check_test.ps1 -w 20 -c 30
```

!!! note

	You can use `$pluginName` and `$pluginArgs` to place the location of those values in the comand. If they are not specified, they are placed at the end of the command.

**Default:**
```
pluginTypes:
  .sh:
    - /bin/sh
  .py:
    - python3
  .pl:
    - perl
  .php:
    - php
  .ps1:
    - powershell
    - -ExecutionPolicy
    - Bypass
    - -File
  .vbs:
    - cscript
    - $pluginName
    - $pluginArgs
    - //NoLogo
  .bat:
    - cmd
    - /c
```

## Passive Check Options

### `senders`

Senders send passive check data somewhere that is set up to recieve passive check data. You can have as many senders as you want, but you will need at least one to send passive checks to. Senders are added as an array in YAML.

There are no senders defined by default.

Example single sender:
```
senders:
  - name: NRDP Server 1
    url: http://192.168.0.100/nrdp/
    token: sometoken
    type: nrdp
```

Example multiple senders:
```
senders:
  - name: NRDP Server 1
    url: http://192.168.0.100/nrdp/
    token: sometoken
    type: nrdp
    senders:
  - name: NRDP Server 2
    url: http://192.168.0.102/nrdp/
    token: sometoken
    type: nrdp

```

!!! note

	NRDP is currently the only sender type.

All the following options are **required** for each sender defined.

#### `name`

The name you want to give the recieving location.

#### `url`

The URL the sender should send the data to.

#### `name`

The token to pass in the request.

#### `type`

Currently, there is only the nrdp sender so type should always be `nrdp`.


### `checks`

Checks are passive checks configured to run at certain intervals while rcagent is running. These checks are sent to some other location (like NRDP) using senders, defined above.

[See the full details about how to set up checks and what options are available.](../checks)

## Other Options

### `manager`

If you are using ReChecked Manager you will need to fill out some options in this section. [See the full details about what options are available.](../manager)

## Example Config

An example configuration. Another helpful example is the [default `config.yml` in the repo](https://github.com/rechecked/rcagent/blob/main/build/package/config.yml).

```
address:
port: 5995 
token: private
tls:
  cert: adhoc
  key: adhoc

defaultUnits: GiB
debug: false

excludeFsTypes: ["aufs", "autofs", "binfmt_misc", "cgroup", "configfs",
  "debugfs", "devpts", "devtmpfs", "encryptfs", "efivarfs", "fuse", "fuseblk",
  "fusectl", "hugetlbfs","mqueue", "overlayfs", "proc", "pstore", "rpc_pipefs",
  "securityfs", "selinuxfs", "sysfs", "tmpfs", "tracefs", "bpf",
  "fuse.vmware-vmblock", "fuse.gvfsd-fuse"]

pluginTypes:
  # Linux
  .sh:
    - /bin/sh
  .py:
    - python3
  .pl:
    - perl
  .php:
    - php
  # Windows
  .ps1:
    - powershell
    - -ExecutionPolicy
    - Bypass
    - -File
  .vbs:
    - cscript
    - $pluginName
    - $pluginArgs
    - //NoLogo
  .bat:
    - cmd
    - /c

senders:
  - name: NRDP Server 1
    url: http://<ip>/nrdp/
    token: <token>
    type: nrdp

checks:
  - hostname: $HOST
    interval: 5m
    endpoint: system/version
    options:
      warning: 10
      critical: 20
  - hostname: $HOST
    servicename: CPU Usage
    interval: 30s
    endpoint: cpu/percent
    options:
      warning: 10
      critical: 20
```