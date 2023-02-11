
package status

import (
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/rechecked/rcagent/internal/config"
)

func HandleDocker(cv config.Values) interface{} {
	dockerIds, _ := docker.GetDockerIDList()
 	return dockerIds
}