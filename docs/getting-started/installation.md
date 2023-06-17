# Installation

We currently build for some of the most popular Linux distros, but we cannot build for all of them. If you don't see your distro here, you can always build it [manually from source](#other-source).

## Linux

### Repo (Recommended)

We recommend using the package manager to install rcagent. You can do this on RPM and DEB based systems with the following steps.

#### CentOS/RHEL

Import GPG Key

```
sudo rpm -v --import https://repo.rechecked.io/rechecked-pub.gpg
```

Add the YUM repo

```
. /etc/os-release
sudo yum-config-manager --add-repo https://repo.rechecked.io/rpm/el$VERSION_ID/rechecked.repo
```

If it says yum-config-manager not found, you will need to install yum-utils.

Install the agent:

```
sudo yum install rcagent
```

#### Debian/Ubuntu

Set up the GPG Key

Add the GPG key to trusted list:

```
wget -qO - https://repo.rechecked.io/rechecked-pub.gpg | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/rechecked-pub.gpg > /dev/null
```

Add repo to apt

```
echo "deb https://repo.rechecked.io/deb $(lsb_release -cs) main" > /etc/apt/sources.list.d/rechecked.list
```

Install the agent:

```
sudo apt-get update
sudo apt-get install rcagent
```

### Manual

#### RPM Based

Download the latest .rpm from the [download page](https://rechecked.io/download/) for the operating system.

```
sudo rpm -ivh rcagent-<version>.rpm
```

#### DEB Based

Download the latest .deb version from the [download page](https://rechecked.io/download/) for the operating system.

```
sudo dpkg -i rcagent-<version>.deb
```

### Running the Agent

The service is not set up to be enabled on restart/boot by default since it requires configuration. Once you’ve updated the config with a [secure authentication token](../configuration/#security-token), enable it and start it:

```
systemctl enable rcagent
systemctl start rcagent
```

## Windows

Download the windows installer from the [download](https://rechecked.io/download/) page.

The msi can be installed manually or you can install the package using the command line with the following command:

```
msiexec /i rcagent-<version>.msi 
```

!!! note
	
	The installer will **not** open the firewall port for you, if you are going to be using active checks, you will need to open the firewall for the port you are going to use. The agent uses port 5995 by default.

On windows, the service is not running after install but it is **set to automatically start on boot**. Once you are finished configuring, start the service by running as admin:

```
sc.exe start rcagent
```

## MacOS

Download the dmg installer from the [download](https://rechecked.io/download/) page.

Mount the dmg file by double clicking it. You will need to run the following command with root permissions or via sudo:

```
sudo zsh /Volumes/rcagent-<version>/install.sh
```

!!! note

	The install.sh script is also used for upgrading on macOS.

The agent is not started by default, but will start on reboot. Once you’ve configured the agent’s config.yml file, you can enable and start it with:

```
sudo launchctl start io.rechecked.rcagent
```

## Other (Source)

If you'd like to run rcagent on a system we don't build for, you'll have to install from source. This is normally a fairly straightforward process, but note that some features may not work properly if they are not implemented on the system you are running it on.

!!! note

	You will need to have a Go compiler installed to build rcagent from source. We recommend installing Go from the official website, [go.dev](https://go.dev), prior to continuing this section.

Download the [source files](https://github.com/rechecked/rcagent/archive/refs/heads/main.zip) from GitHub:
```
https://github.com/rechecked/rcagent/archive/refs/heads/main.zip
```

Unzip the files into a directory, navigate to the directory. The following commands must be ran inside the source directory.

To build rcagent, just run:

```
make build
```

To install the source version, run the following:

```
make install
```

If you'd like to run the source version as a service, you can install the service by running:

```
/usr/local/rcagent/rcagent -a install
```

You should now have a running rcagent.
