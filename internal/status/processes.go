
package status

import (
    "github.com/shirou/gopsutil/v3/process"
    "github.com/rechecked/rcagent/internal/config"
)

type Process struct {
    Name    string `json:"name"`
    PID     int32  `json:"pid"`
    Exe     string `json:"exe"`
    Cmdline string `json:"cmdline"`
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
        procs = append(procs, Process{
            Name: name,
            Exe: exe,
            Cmdline: cmdline,
            PID: p.Pid,
        })
    }
    return procs
}
