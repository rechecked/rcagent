//go:build !windows
// +build !windows

package status

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
)

func getUsers() ([]host.UserStat, error) {
	users, err := host.Users()
	if err != nil {
		fmt.Printf("%s", err)
		return users, err
	}
	return users, nil
}
