
package status

import (
    "fmt"
    "strings"
    "github.com/shirou/gopsutil/v3/mem"
    "github.com/rechecked/rcagent/internal/config"
)

type MemoryStatus struct {
    Total       float64 `json:"total"`
    Available   float64 `json:"available"`
    Free        float64 `json:"free"`
    Used        float64 `json:"used"`
    UsedPercent float64 `json:"usedPercent"`
    Units       string  `json:"units"`
}

type SwapStatus struct {
    Total       float64 `json:"total"`
    Free        float64 `json:"free"`
    Used        float64 `json:"used"`
    UsedPercent float64 `json:"usedPercent"`
    Units       string  `json:"units"`
}

func (m MemoryStatus) String() string {
    return fmt.Sprintf("Memory usage is %.2f%% (%.2f/%.2f %s Total)", m.UsedPercent,
        m.Used, m.Total, m.Units)
}

func (s SwapStatus) String() string {
    return fmt.Sprintf("Swap usage is %.2f%% (%.2f/%.2f %s Total)", s.UsedPercent,
        s.Used, s.Total, s.Units)
}

func (m MemoryStatus) CheckValue() float64 {
    return m.UsedPercent
}

func (s SwapStatus) CheckValue() float64 {
    return s.UsedPercent
}

func (m MemoryStatus) PerfData(warn, crit string) string {
    var perfdata []string
    data := fmt.Sprintf("'percent'=%0.2f%%", m.UsedPercent)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'available'=%0.2f%s", m.Available, m.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'used'=%0.2f%s", m.Used, m.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'free'=%0.2f%s", m.Free, m.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'total'=%0.2f%s", m.Total, m.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    return strings.Join(perfdata, " ")
}

func (s SwapStatus) PerfData(warn, crit string) string {
    var perfdata []string
    data := fmt.Sprintf("'percent'=%0.2f%%", s.UsedPercent)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'used'=%0.2f%s", s.Used, s.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'free'=%0.2f%s", s.Free, s.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'total'=%0.2f%s", s.Total, s.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    return strings.Join(perfdata, " ")
}

func HandleMemory(cv config.Values) interface{} {
    data, err := memoryUsage(cv.Units())
    if err != nil {
        return MemoryStatus{}
    }
    return data
}

func HandleSwap(cv config.Values) interface{} { 
    data, err := swapUsage(cv.Units())
    if err != nil {
        return SwapStatus{}
    }
    return data
}

func memoryUsage(units string) (MemoryStatus, error) {
    v, err := mem.VirtualMemory()
    // usedPercent is rounded for some reason, so we are
    // calculating that ourselves which is just used/total*100
    mem := MemoryStatus{
        Total: ConvertToUnit(v.Total, units),
        Available: ConvertToUnit(v.Available, units),
        Used: ConvertToUnit(v.Used, units),
        Free: ConvertToUnit(v.Free, units),
        UsedPercent: (float64(v.Used) / float64(v.Total) *100),
        Units: units,
    }
    return mem, err
}

func swapUsage(units string) (SwapStatus, error) {
    s, err := mem.SwapMemory()
    swap := SwapStatus{
        Total: ConvertToUnit(s.Total, units),
        Used: ConvertToUnit(s.Used, units),
        Free: ConvertToUnit(s.Free, units),
        UsedPercent: s.UsedPercent,
        Units: units,
    }
    return swap, err
}
