
package status

import (
    "fmt"
    "strings"
    "github.com/shirou/gopsutil/v3/disk"
    "github.com/rechecked/rcagent/internal/config"
)

type Disk struct {
    Path        string  `json:"path"`
    Device      string  `json:"device"`
    Fstype      string  `json:"fstype"`
    Total       float64 `json:"total"`
    Free        float64 `json:"free"`
    Used        float64 `json:"used"`
    UsedPercent float64 `json:"usedPercent"`
    Units       string  `json:"units"`
}

type Inodes struct {
    Path        string  `json:"path"`
    Device      string  `json:"device"`
    Fstype      string  `json:"fstype"`
    Total       float64 `json:"total"`
    Used        float64 `json:"used"`
    Free        float64 `json:"free"`
    UsedPercent float64 `json:"usedPercent"`
}

func (d Disk) String() string {
    return fmt.Sprintf("Disk usage of %s is %.2f%% (%.2f/%.2f %s Total)", d.Path,
        d.UsedPercent, d.Used, d.Total, d.Units)
}

func (i Inodes) String() string {
    return fmt.Sprintf("Inode usage of %s is %.2f%% (%.0f/%.0f Total)", i.Path,
        i.UsedPercent, i.Used, i.Total)
}

func (d Disk) PerfData(warn, crit string) string {
    var perfdata []string
    data := fmt.Sprintf("'percent'=%0.2f%%", d.UsedPercent)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'used'=%0.2f%s", d.Used, d.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'free'=%0.2f%s", d.Free, d.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'total'=%0.2f%s", d.Total, d.Units)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    return strings.Join(perfdata, " ")
}

func (i Inodes) PerfData(warn, crit string) string {
    var perfdata []string
    data := fmt.Sprintf("'percent'=%0.2f%%", i.UsedPercent)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'used'=%0.0f", i.Used)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'free'=%0.0f", i.Free)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    data = fmt.Sprintf("'total'=%0.0f", i.Total)
    perfdata = append(perfdata, createPerfData(data, warn, crit))
    return strings.Join(perfdata, " ")
}

func (d Disk) CheckValue() float64 {
    return d.UsedPercent
}

func (i Inodes) CheckValue() float64 {
    return i.UsedPercent
}

func HandleDisks(cv config.Values) interface{} { 

    disks, _ := getDisks()

    // Find the specific disk if we are passing a path
    if cv.Path != "" {
        disk, err := getDiskFromPath(disks, cv.Path, cv.Units(), false)
        if err != nil {
            return err
        }
        return disk
    }

    return disks
}

func HandleInodes(cv config.Values) interface{} {
    disks, _ := getDisks()

     // Find the specific disk if we are passing a path
    if cv.Path != "" {
        disk, _ := getDiskFromPath(disks, cv.Path, "", true)
        return disk
    } else {
        var inodes []Inodes
        for _, d := range disks {
            i, _ := getDiskFromPath(disks, d.Mountpoint, "", true)
            inode, _ := i.(Inodes)
            inodes = append(inodes, inode)
        }
        return inodes
    }

    return disks
}

func getDisks() ([]disk.PartitionStat, error) {
    var disks []disk.PartitionStat
    d, err := disk.Partitions(true)
    if err != nil {
        return []disk.PartitionStat{}, err
    }
    for _, disk := range d {
        if !config.Contains(config.Settings.ExcludeFsTypes, disk.Fstype) {
            disks = append(disks, disk)
        }
    }
    return disks, nil
}

func getDiskFromPath(disks []disk.PartitionStat, path, units string, inodes bool) (interface{}, error) {

    for _, d := range disks {
        if d.Mountpoint == path {
            disk, _ := disk.Usage(d.Mountpoint)

            if inodes {
                return Inodes{
                    Path: disk.Path,
                    Device: d.Device,
                    Fstype: disk.Fstype,
                    Total: float64(disk.InodesTotal),
                    Free: float64(disk.InodesFree),
                    Used: float64(disk.InodesUsed),
                    UsedPercent: disk.InodesUsedPercent,
                }, nil
            }
            return Disk{
                Path: disk.Path,
                Device: d.Device,
                Fstype: disk.Fstype,
                Total: ConvertToUnit(disk.Total, units),
                Free: ConvertToUnit(disk.Free, units),
                Used: ConvertToUnit(disk.Used, units),
                UsedPercent: disk.UsedPercent,
                Units: units,
            }, nil
        }
    }

    return Disk{}, fmt.Errorf("The path (%s) does not exist as a mountpount.", path)
}
