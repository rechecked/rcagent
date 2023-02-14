
package status

import (
    "fmt"
    "strings"
    "github.com/shirou/gopsutil/v3/host"
    "github.com/rechecked/rcagent/internal/config"
)

type Users struct {
    count float64
}

type Version struct {
    Version string `json:"version"`
}

func (u Users) String() string {
    return fmt.Sprintf("Current users count is %0.0f", u.CheckValue())
}

func (u Users) CheckValue() float64 {
    return u.count
}

func (u Users) PerfData(warn, crit string) string {
    var perfdata []string
    data := fmt.Sprintf("'users'=%0.0f", u.CheckValue())
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    return strings.Join(perfdata, " ")
}

func (v Version) String() string {
    return fmt.Sprintf("rcagent version is %s", v.Version)
}

func (v Version) CheckValue() float64 {
    return 0.0
}

func (v Version) PerfData(warn, crit string) string {
    return ""
}

func HandleSystem(cv config.Values) interface{} {
    data, err := host.Info()
    if err != nil {
        return err
    }
    return data
}

func HandleTemps(cv config.Values) interface{} {
    temps, err := host.SensorsTemperatures()
    if err != nil {
        return err
    }
    return temps
}

func HandleVersion(cv config.Values) interface{} {
    return Version{
        Version: config.Version,
    }
}

func HandleUsers(cv config.Values) interface{} {
    users, err := getUsers()
    if err != nil {
        return err
    }
    if cv.Check {
        return Users{
            count: float64(len(users)),
        }
    }
    return users
}
