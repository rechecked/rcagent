
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
