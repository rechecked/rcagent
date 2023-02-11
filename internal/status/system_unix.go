// +build !windows

package status

import (
    "github.com/shirou/gopsutil/v3/host"
)

func getUsers() ([]host.UserStat, error) {
    users, err := host.Users()
    if err != nil {
        return users, err
    }
    return users, nil
}