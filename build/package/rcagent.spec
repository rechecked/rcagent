Name:           rcagent
Version:        1.0.0
Release:        1%{?dist}
Vendor:         ReChecked
Summary:        Cross platform system status and monitoring agent for monitoring systems.
Group:          Network/Monitoring
License:        GPLv3
URL:            https://github.com/rechecked/rcagent
Source0:        %{name}-%{version}.tar.gz

%if 0%{?fedora} || 0%{?rhel} || 0%{?centos}
BuildRequires:  golang
BuildRequires:  systemd-rpm-macros
%endif

Requires(pre): 	shadow-utils

Provides:       %{name} = %{version}

%description
The rcagent is a Nagios-compatable monitoring agent and system status
reporter created to run active and passive checks on multiple operating systems.

%global debug_package %{nil}

%prep
%setup -q

%build
LDFLAGS="-X github.com/rechecked/rcagent/internal/config.PluginDir=%{_libdir}/%{name}/plugins \
-X github.com/rechecked/rcagent/internal/config.ConfigDir=%{_sysconfdir}/%{name}" \
%{__make} build

%install
mkdir -p $RPM_BUILD_ROOT/%{_libdir}/%{name}/plugins
%{__install} -Dpm 0755 build/bin/%{name} $RPM_BUILD_ROOT/%{_sbindir}/%{name}
%{__install} -Dpm 0755 build/package/config.yml $RPM_BUILD_ROOT/%{_sysconfdir}/%{name}/config.yml

%clean
%{__rm} -rf $RPM_BUILD_ROOT

%check
%{__make} test

%pre
getent group rcagent >/dev/null || groupadd -r rcagent
getent passwd rcagent >/dev/null || \
    useradd -r -g rcagent -d %{_libdir}/%{name} -s /sbin/nologin \
    -c "rcagent user account for running plugins" rcagent

%post
if [ $1 -eq 1 ]
then
    # Install sets up systemctl service so it only runs on install
    %{_sbindir}/%{name} -a install &> /dev/null
    
    # Disable on systemctl during install because we don't want to default enabled
    if command -v systemctl > /dev/null
    then
        systemctl disable %{name}.service &> /dev/null
    fi
fi

%preun
# On uninstall stop before removing
if [ $1 -eq 0 ]
then
	systemctl stop %{name}.service &> /dev/null
    %{_sbindir}/%{name} -a uninstall &> /dev/null
fi

%posttrans
# Restart the service if it needs to be restarted after upgrade
if [ $1 -eq 2 ]
then
    if command -v systemctl > /dev/null
    then
        systemctl is-active --quiet %{name}.service && systemctl restart %{name}.service &> /dev/null
    fi
fi

%files
%config(noreplace) %{_sysconfdir}/%{name}/config.yml
%{_sbindir}/%{name}

%dir %{_libdir}/%{name}/plugins

