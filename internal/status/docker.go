package status

import (
	"github.com/rechecked/rcagent/internal/config"
	"github.com/shirou/gopsutil/v3/docker"
)

func HandleDocker(cv config.Values) interface{} {
	dockerIds, _ := docker.GetDockerIDList()
	return dockerIds
}
