package config

import (
	"errors"
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/kardianos/service"
	"github.com/shirou/gopsutil/v3/host"
)

type HostInfo struct {
	Hostname  string
	MachineId string
	OS        string
	Platform  string
}

var Log service.Logger

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !errors.Is(err, os.ErrNotExist)
}

func Contains(s []string, val string) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}
	return false
}

func UsingManager() bool {
	return Settings.Manager.APIKey != ""
}

func LogDebug(v ...interface{}) {
	if DebugMode {
		Log.Info(v...)
	}
}

func LogDebugf(format string, a ...interface{}) {
	if DebugMode {
		Log.Infof(format, a...)
	}
}

func GetMachineId() string {
	return GetHostInfo().MachineId
}

func GetHostInfo() HostInfo {

	hostname, _ := os.Hostname()
	machineId, _ := machineid.ProtectedID("rcagent")
	host, _ := host.Info()

	i := HostInfo{
		Hostname:  hostname,
		MachineId: machineId,
		OS:        host.OS,
		Platform:  host.Platform,
	}

	return i
}
