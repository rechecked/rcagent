
package status

import (
    "github.com/shirou/gopsutil/v3/host"
    "github.com/rechecked/rcagent/internal/config"
)

func HandleSystem(cv config.Values) interface{} {
    data, err := host.Info()
    if err != nil {
        return nil
    }
    return data
}

