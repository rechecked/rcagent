package status

import (
	"github.com/rechecked/rcagent/internal/config"
	"github.com/shirou/gopsutil/v3/docker"
)

func HandleDocker(cv config.Values) interface{} {
	dockerIds, err := docker.GetDockerIDList()
	if err != nil {
		return []string{}
	}
	return dockerIds
}
