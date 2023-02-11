// +build !windows

package status

import (
    "fmt"
    "github.com/shirou/gopsutil/v3/host"
)

func getUsers() ([]host.UserStat, error) {
    users, err := host.Users()
    fmt.Printf("%s", err)
    if err != nil {
        fmt.Printf("%s", err)
        return users, err
    }
    return users, nil
}