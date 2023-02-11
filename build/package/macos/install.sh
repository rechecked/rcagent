#!/bin/sh

# Verify admin permissions
if ! id -Gn $(whoami) | grep -q -w "admin"; then
    echo "Error: You must have admin permissions to run this script."
    exit 1
fi

DIR="/usr/local/rcagent"
CFGDIR="/etc/rcagent"
UPGRADE=0

# Check installation
if [ -d $DIR ]; then
    UPGRADE=1
    echo "Upgrading ReChecked Agent"
else
    echo "Installing ReChecked Agent"
fi

if [ $UPGRADE -eq 1 ]; then

    # Do upgrade

    # TODO: Add upgrade

else

    # Do install

    # Make directories
    mkdir -p $DIR
    mkdir $DIR/plugins
    mkdir $CFGDIR

    # Fix apple message
    xattr -d -r com.apple.quarantine $DIR

    # Copy files
    cp config.yml $CFGDIR/config.yml
    cp uninstall.sh $DIR/uninstall.sh
    cp rcagent $DIR/rcagent

    # Create symlink
    ln -s $DIR/rcagent /usr/local/bin/rcagent

    # Create rcagent user/group for plugins
    if ! dscl . -read /Groups/rcagent > /dev/null 2>&1; then
        PrimaryGroupID=`dscl . -list /Groups PrimaryGroupID | awk '{print $2}' | sort -ug | tail -1`
        let PrimaryGroupID=PrimaryGroupID+1
        dscl . -create /Groups/rcagent
        dscl . -create /Groups/rcagent RecordName "_rcagent rcagent"
        dscl . -create /Groups/rcagent PrimaryGroupID $PrimaryGroupID
        dscl . -create /Groups/rcagent RealName "rcagent"
        dscl . -create /Groups/rcagent Password "*"
    fi
    if ! dscl . -read /Users/rcagent > /dev/null 2>&1; then
        UniqueID=`dscl . -list /Users UniqueID | awk '{print $2}' | sort -ug | tail -1`
        let UniqueID=UniqueID+1
        dscl . -create /Users/rcagent
        dscl . -create /Users/rcagent UserShell /usr/bin/false
        dscl . -create /Users/rcagent UniqueID $UniqueID
        dscl . -create /Users/rcagent RealName "rcagent"
        dscl . -create /Users/rcagent PrimaryGroupID $PrimaryGroupID
        dscl . -create /Users/rcagent Password "*"
        dscl . -create /Users/rcagent NFSHomeDirectory $DIR
    fi

    # Run the rcagent service install
    rcagent -a install
    if [[ $(xattr -l /Library/LaunchDaemons/io.rechecked.rcagent.plist) ]]; then
        xattr -d com.apple.quarantine /Library/LaunchDaemons/io.rechecked.rcagent.plist
    fi
    launchctl load io.rechecked.rcagent

    echo "Installation complete!"
    echo "- Configure rcagent token: $CFGDIR/config.yml"
    echo "- Start rcagent: launchctl start io.rechecked.rcagent"

fi
