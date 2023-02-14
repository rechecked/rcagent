#!/bin/sh

# Verify admin permissions
if ! id -Gn $(whoami) | grep -q -w "admin"; then
    echo "Error: You must have admin permissions to run this script."
    exit 1
fi

DIR="/usr/local/rcagent"
CFGDIR="/etc/rcagent"

echo "Uninstalling ReChecked Agent"

# Disable service
launchctl stop io.rechecked.rcagent

# Remove service
/usr/local/rcagent/rcagent -a uninstall

# Remove install directory
rm -rf $DIR
rm -rf $CFGDIR
