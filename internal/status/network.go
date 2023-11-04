package status

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rechecked/rcagent/internal/config"
	"github.com/shirou/gopsutil/v3/net"
)

var ifmux = &sync.Mutex{}
var interfaces = make(map[string]Interface)

type Interface struct {
	HardwareAddr string                `json:"hardwareAddr"`
	Addrs        net.InterfaceAddrList `json:"addrs"`
	stored       time.Time
	net.IOCountersStat
}

type InterfaceDelta struct {
	HardwareAddr string                `json:"hardwareAddr"`
	Addrs        net.InterfaceAddrList `json:"addrs"`
	net.IOCountersStat
	InterfaceDeltaStat
}

type InterfaceDeltaStat struct {
	OutTotal   float64 `json:"outTotal"`
	OutPerSec  float64 `json:"outPerSec"`
	InTotal    float64 `json:"inTotal"`
	InPerSec   float64 `json:"inPerSec"`
	Units      string  `json:"units"`
	checkType  string
	checkValue float64
}

func (i InterfaceDeltaStat) String() string {
	if i.checkType == "in" {
		return fmt.Sprintf("Current inbound traffic %0.2f %s/s",
			i.checkValue, i.Units)
	}
	if i.checkType == "out" {
		return fmt.Sprintf("Current outbound traffic %0.2f %s/s",
			i.checkValue, i.Units)
	}
	return fmt.Sprintf("Current network traffic %0.2f %s/s",
		i.checkValue, i.Units)
}

func (i InterfaceDeltaStat) CheckValue() float64 {
	return i.checkValue
}

func (i InterfaceDeltaStat) PerfData(warn, crit string) string {
	var perfdata []string
	var data string

	if i.checkType == "in" || i.checkType == "total" {
		data = fmt.Sprintf("'in'=%0.2f%s/s", i.InPerSec, i.Units)
		perfdata = append(perfdata, createPerfData(data, warn, crit))
	}
	if i.checkType == "out" || i.checkType == "total" {
		data = fmt.Sprintf("'out'=%0.2f%s/s", i.OutPerSec, i.Units)
		perfdata = append(perfdata, createPerfData(data, warn, crit))
	}
	return strings.Join(perfdata, " ")
}

// Save initial counters for per second values
func Setup() {
	ifs, _ := getNetworkIfs()
	for _, i := range ifs {
		setInterfaceStats(i)
	}
}

func setInterfaceStats(itr Interface) {
	ifmux.Lock()
	interfaces[itr.Name] = itr
	ifmux.Unlock()
}

func getInterfaceStats(itr string) Interface {
	ifmux.Lock()
	defer ifmux.Unlock()
	return interfaces[itr]
}

// Get a whole list of interfaces
func HandleNetworks(cv config.Values) interface{} {
	ifs, _ := getNetworkIfs()

	if cv.Name != "" {
		var itr Interface
		for _, i := range ifs {
			if i.Name == cv.Name {
				itr = i
			}
		}

		// If delta is passed, get the old value and adjust it
		// based on amount of seconds passed
		if cv.Delta > 0 || cv.Check {
			deltaItr := getInterfaceDeltaStat(itr, cv)
			return deltaItr
		}

		return itr
	}

	// If delta, we need the InterfaceDeltaStat added to the interfaces
	if cv.Delta > 0 {
		var ifDeltaList []InterfaceDelta
		for _, i := range ifs {
			deltaItr := getInterfaceDeltaStat(i, cv)
			ifDeltaList = append(ifDeltaList, InterfaceDelta{
				i.HardwareAddr,
				i.Addrs,
				net.IOCountersStat{
					Name:        i.Name,
					BytesSent:   i.BytesSent,
					BytesRecv:   i.BytesRecv,
					PacketsSent: i.PacketsSent,
					PacketsRecv: i.PacketsRecv,
					Errin:       i.Errin,
					Errout:      i.Errout,
					Dropin:      i.Dropin,
					Dropout:     i.Dropout,
					Fifoin:      i.Fifoin,
					Fifoout:     i.Fifoout,
				},
				deltaItr,
			})
		}

		return ifDeltaList
	}

	return ifs
}

func getNetworkIfs() ([]Interface, error) {
	var ifList []Interface

	ifs, err := net.Interfaces()
	if err != nil {
		return ifList, err
	}

	ifsCounters, err := net.IOCounters(true)
	if err != nil {
		return ifList, err
	}

	for _, i := range ifs {
		// Append the counter data onto each of the interfaces
		for _, x := range ifsCounters {
			if x.Name == i.Name {
				ifList = append(ifList, Interface{
					i.HardwareAddr,
					i.Addrs,
					time.Now(),
					x,
				})
			}
		}
	}

	return ifList, nil
}

func getInterfaceDeltaStat(itr Interface, cv config.Values) InterfaceDeltaStat {
	tmpItr := getInterfaceStats(itr.Name)
	timeSince := itr.stored.Sub(tmpItr.stored)

	dOut := itr.BytesRecv - tmpItr.BytesRecv
	in := float64(dOut) / timeSince.Seconds()
	inPs := ConvertToUnit(uint64(in), cv.Units())
	dIn := itr.BytesSent - tmpItr.BytesSent
	out := float64(dIn) / timeSince.Seconds()
	outPs := ConvertToUnit(uint64(out), cv.Units())

	// Get check value
	var cVal float64
	if cv.Against == "in" {
		cVal = inPs
	} else if cv.Against == "out" {
		cVal = outPs
	} else {
		cv.Against = "total"
		cVal = inPs + outPs
	}

	deltaItr := InterfaceDeltaStat{
		OutTotal:   ConvertToUnit(dOut, cv.Units()),
		OutPerSec:  outPs,
		InTotal:    ConvertToUnit(dIn, cv.Units()),
		InPerSec:   inPs,
		Units:      cv.Units(),
		checkType:  cv.Against,
		checkValue: cVal,
	}

	// Save the current itr for later
	setInterfaceStats(itr)

	return deltaItr
}
