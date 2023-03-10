package status

import (
	"fmt"
	"github.com/rechecked/rcagent/internal/config"
	"github.com/shirou/gopsutil/v3/process"
	"strings"
)

type Process struct {
	Name    string `json:"name"`
	PID     int32  `json:"pid"`
	Exe     string `json:"exe"`
	Cmdline string `json:"cmdline"`
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
		if cv.Name != "" {
			if name != cv.Name {
				continue
			}
		}
		procs = append(procs, Process{
			Name:    name,
			Exe:     exe,
			Cmdline: cmdline,
			PID:     p.Pid,
		})
	}

	pList := ProcessList{
		Count:     len(procs),
		Processes: procs,
	}

	return pList
}
