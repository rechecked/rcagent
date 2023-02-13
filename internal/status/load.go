
package status

import (
    "fmt"
    "strings"
    "github.com/shirou/gopsutil/v3/load"
    "github.com/rechecked/rcagent/internal/config"
)

type Load struct {
    Load1   float64 `json:"load1"`
    Load5   float64 `json:"load5"`
    Load15  float64 `json:"load15"`
    against string
}

func (l Load) String() string {
    return fmt.Sprintf("Load average is %.2f, %.2f, %.2f", l.Load1, l.Load5, l.Load15)    
}

func (l Load) PerfData(warn, crit string) string {
    var perfdata []string
    data := fmt.Sprintf("'load1'=%0.2f", l.Load1)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'load5'=%0.2f", l.Load5)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'load15'=%0.2f", l.Load15)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    return strings.Join(perfdata, " ")
}

func (l Load) CheckValue() float64 {

    // Check if we are running a check against something
    // other than the default (load1)
    switch l.against {
    case "highest":
        var highest float64
        tmp := []float64{l.Load1, l.Load5, l.Load15}
        for _, t := range tmp {
            if t > highest {
                highest = t
            }
        }
        return highest
    case "load5":
        return l.Load5
    case "load15":
        return l.Load15
    }

    return l.Load1
}

func HandleLoad(cv config.Values) interface{} {
    avg, _ := load.Avg()
    load := Load{
        Load1: avg.Load1,
        Load5: avg.Load5,
        Load15: avg.Load15,
        against: cv.Against,
    }

    return load
}