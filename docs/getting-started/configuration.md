# Configuration

## Main Config File

The agent uses YAML configuration files to define it's main config and also define senders, checks, and other details about passive checks. Depending on the system you are running on and where you installed rcagent, the file will be in a different directory. The default directories are:

- Linux: /etc/rcagent/config.yml
- Windows: C:\Program Files\rcagent\config.yml
- MacOS: /etc/rcagent/config.yml

To get started, all you have to do for now is change the token value to something secure and private. By default the token is **private**.

```
token: private # <=== CHANGE THIS ====
```

Replace the token with whatever value you’d like. Refer back to step one to start the service and enable it if you haven’t already.

For a full list of values you can use in the configuration, check the configuration file options documentation.

## Check Configuration

Now that the agent has been installed and you have a secure token, we can set up checks. Depending on your situation, you may want to use active vs passive checks. One reason to use passive checks is if you do not want to open your firewall, you can send out passive checks via NRDP with HTTP connections.

### Active Checks

#### Manual

You can manually run active checks using the `check_rcagent.py` script.



#### Nagios XI

If you’re running Nagios XI you can download the Nagios XI Config Wizard and install it on your system through `Admin > Manage Config Wizards > Upload & Install` and select the `rcagent.zip` config wizard file.

Once installed, the config wizard will let you set up active checks though an interactive interface.

#### Nagios Core

For Nagios Core you will need to save the `check_rcagent.py` into your `/usr/local/nagios/libexec` directory before adding configurations to your Nagios Core system.

You’ll need to make a command first, in your commands.cfg file. Normally this is stored in either `/etc/nagios/` or `/usr/local/nagios/etc/`. You can also add it wherever you store your commands.

```
define command {
    command_name    check_rcagent
    command_line    $USER1$/check_rcagent.py -H $HOSTADDRESS$ $ARG1$
}
```

Passing `$ARG1$` at the end allows us to manage how we want to pass arguments without needing to make extra commands for more complicated checks.

Then you’ll be able to create hosts and services, the below is an example of a service using our above command and passing arguments.

```
define host {
    host_name               RCAgent Test Host
    address                 192.168.1.100
    check_command           check_rcagent!-t private -e system/version
    max_check_attempts      5
    check_interval          5
    retry_interval          1
    check_period            24x7
    contacts                admin
    notification_interval   60
    notification_period     24x7
    notifications_enabled   1
    icon_image              rcagent.png
    statusmap_image         rcagent.png
    register                1
}

define service {
    host_name               RCAgent Test Host
    service_description     CPU Usage
    check_command           check_rcagent!-t private -e cpu/percent -w 20 -c 40
    max_check_attempts      5
    check_interval          5
    retry_interval          1
    check_period            24x7
    notification_interval   60
    notification_period     24x7
    contacts                admin
    register                1
}
```

### Passive Checks

You can add individual passive checks to be sent over NRDP by adding the following to your `config.yml`.

Create a senders section and add the NRDP server to send to with the token:

```
senders:
  - name: NRDP Server 1
    url: http://<ip>/nrdp/
    token: <token>
    type: nrdp
```

!!! note

	All passive checks will be sent to all senders created at this time.

Next, set up the passive checks you wish to send, in this example we will send a simple version check for the host and a cpu usage check for the service. You must create the checks section if it doesn’t already exist.

```
checks:
  - hostname: $HOST
    interval: 5m
    endpoint: cpu/percent
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

The $HOST variable is the hostname of the system the rcagent is running on and gets populated during runtime. The options section allows you to pass parameters just like the URL for active checks via the status API. This is why we pass warning/critical values in this way.

For a full list of options for checks and senders check the config file reference section.
