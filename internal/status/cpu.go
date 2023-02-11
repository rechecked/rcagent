
package status

import (
    "fmt"
    "time"
    "strings"
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/rechecked/rcagent/internal/config"
)

type CPUStatus struct {
    Percent []float64      `json:"percent"`
    Units   string         `json:"units"`
    info    []cpu.InfoStat
}

func (c CPUStatus) String() string {
    return fmt.Sprintf("CPU usage is %.2f%%", c.CheckValue())
}

func (c CPUStatus) PerfData(warn, crit string) string {
    var perfdata []string
    data := fmt.Sprintf("'percent'=%0.2f%%", c.CheckValue())
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    return strings.Join(perfdata, " ")
}

func (c CPUStatus) LongOutput() string {

    var models string
    cpus := make(map[string]int)

    // Get total cores
    cores := 0
    for _, i := range c.info {
        cores += int(i.Cores)
        cpus[i.ModelName]++
    }

    for c, n := range cpus {
        if n > 1 {
            models += fmt.Sprintf("%d x %s", n, c)
        } else {
            models += fmt.Sprintf("%s", c)
        }
    }

    return fmt.Sprintf("%s [Total Cores: %d]", models, cores)
}

func (c CPUStatus) CheckValue() float64 {

    // Get average of all percents
    var avg float64
    for _, p := range c.Percent {
        avg += p
    }
    avg = avg / float64(len(c.Percent))

    return avg
}

func HandleCPU(cv config.Values) interface{} {
    waitTime := 500 * time.Millisecond
    // Custom wait time for cpu usage (delta=<seconds>)
    /*
    if cv.Delta > 0 {
        waitTime = time.Duration(cv.Delta) * time.Second
    }
    */
    p, err := cpu.Percent(waitTime, false)
    i, err := cpu.Info()
    if err != nil {
        return nil
    }
    return CPUStatus{
        Percent: p,
        Units: "%",
        info: i,
    }
}
