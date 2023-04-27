package status

import (
	"fmt"
	"strings"

	"github.com/rechecked/rcagent/internal/config"
	"github.com/shirou/gopsutil/v3/process"
)

type Process struct {
	Name       string   `json:"name"`
	PID        int32    `json:"pid"`
	Exe        string   `json:"exe"`
	Cmdline    string   `json:"cmdline"`
	Username   string   `json:"username"`
	CPUPercent float64  `json:"cpuPercent"`
	MemPercent float64  `json:"memPercent"`
	Status     []string `json:"status"`
}

type ProcessList struct {
	Count     int       `json:"count"`
	Processes []Process `json:"processes"`
}

func (p ProcessList) String() string {
	return fmt.Sprintf("Process count is %d", p.Count)
}

func (p ProcessList) CheckValue() float64 {
	return float64(p.Count)
}

func (p ProcessList) PerfData(warn, crit string) string {
	var perfdata []string
	data := fmt.Sprintf("'processes'=%0.f", p.CheckValue())
	perfdata = append(perfdata, createPerfData(data, warn, crit))
	return strings.Join(perfdata, " ")
}

func HandleProcesses(cv config.Values) interface{} {
	var procs []Process
	data, err := process.Processes()
	if err != nil {
		return []Process{}
	}

	for _, p := range data {
		name, _ := p.Name()
		exe, _ := p.Exe()
		cmdline, _ := p.Cmdline()
		username, _ := p.Username()
		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		status, _ := p.Status()
		if cv.Name != "" {
			if name != cv.Name {
				continue
			}
		}
		procs = append(procs, Process{
			Name:     name,
			Exe:      exe,
			Cmdline:  cmdline,
			PID:      p.Pid,
			Username: username,
			CPUPercent: cpuPercent,
			MemPercent: float64(memPercent),
			Status: status,
		})
	}

	pList := ProcessList{
		Count:     len(procs),
		Processes: procs,
	}

	return pList
}
