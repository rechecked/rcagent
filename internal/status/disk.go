package status

import (
	"fmt"
	"github.com/rechecked/rcagent/internal/config"
	"github.com/shirou/gopsutil/v3/disk"
	"strings"
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
	disks, _ := getDisks(cv.Units())

	// Find the specific disk if we are passing a path
	if cv.Path != "" {
		for _, disk := range disks {
			if cv.Path == disk.Path {
				return disk
			}
		}
	}

	return disks
}

func HandleInodes(cv config.Values) interface{} {
	disks, _ := getDisksInodes(cv.Units())

	// Find the specific disk if we are passing a path
	if cv.Path != "" {
		for _, disk := range disks {
			if cv.Path == disk.Path {
				return disk
			}
		}
	}

	return disks
}

func getDisks(units string) ([]Disk, error) {
	var disks []Disk
	d, err := disk.Partitions(true)
	if err != nil {
		return disks, err
	}
	for _, i := range d {
		if !config.Contains(config.Settings.ExcludeFsTypes, i.Fstype) {
			u, _ := disk.Usage(i.Mountpoint)
			disks = append(disks, Disk{
				Path:        i.Mountpoint,
				Device:      i.Device,
				Fstype:      i.Fstype,
				Total:       ConvertToUnit(u.Total, units),
				Free:        ConvertToUnit(u.Free, units),
				Used:        ConvertToUnit(u.Used, units),
				UsedPercent: u.UsedPercent,
				Units:       units,
			})
		}
	}
	return disks, nil
}

func getDisksInodes(units string) ([]Inodes, error) {
	var inodes []Inodes
	d, err := disk.Partitions(true)
	if err != nil {
		return inodes, err
	}
	for _, i := range d {
		if !config.Contains(config.Settings.ExcludeFsTypes, i.Fstype) {
			u, _ := disk.Usage(i.Mountpoint)
			inodes = append(inodes, Inodes{
				Path:        i.Mountpoint,
				Device:      i.Device,
				Fstype:      i.Fstype,
				Total:       float64(u.InodesTotal),
				Free:        float64(u.InodesFree),
				Used:        float64(u.InodesUsed),
				UsedPercent: u.InodesUsedPercent,
			})
		}
	}
	return inodes, nil
}
