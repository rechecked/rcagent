//go:build !windows
// +build !windows

package status

import (
	"github.com/rechecked/rcagent/internal/config"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func setUser(cmd *exec.Cmd) {

	// Run plugins as rcagent on *nix systems (unless runPluginsAsRoot is turned on)
	if !config.Settings.RunPluginsAsRoot {
		u, err := user.Lookup("rcagent")
		if err != nil {
			return
		}

		uid, _ := strconv.ParseInt(u.Uid, 10, 32)
		gid, _ := strconv.ParseInt(u.Gid, 10, 32)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Credential: &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)},
		}
	}

}
