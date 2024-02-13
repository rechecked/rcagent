//go:build !windows
// +build !windows

package status

import (
	"github.com/shirou/gopsutil/v3/host"

	"github.com/rechecked/rcagent/internal/config"
)

func getUsers() ([]host.UserStat, error) {
	users, err := host.Users()
	if err != nil {
		config.Log.Error(err)
		return users, err
	}
	return users, nil
}
