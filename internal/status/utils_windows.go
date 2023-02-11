// +build windows

package status

import (
    "os/exec"
)

func setUser(cmd *exec.Cmd) {

    // Do nothing for now, on Windows we are running as the system
    // user and as a service

}