
package status

import (
    "fmt"
    "sync"
    "time"
    "strings"
    "github.com/shirou/gopsutil/v3/net"
    "github.com/rechecked/rcagent/internal/config"
)

var ifmux = &sync.Mutex{}
var interfaces = make(map[string]Interface)

type Interface struct {
    HardwareAddr       string                `json:"hardwareAddr"`
    Addrs              net.InterfaceAddrList `json:"addrs"`
    stored             time.Time
    net.IOCountersStat
}

type InterfaceDelta struct {
    Name       string  `json:"name"`
    OutTotal   float64 `json:"outTotal"`
    OutPerSec  float64 `json:"outPerSec"`
    InTotal    float64 `json:"inTotal"`
    InPerSec   float64 `json:"inPerSec"`
    Units      string  `json:"units"`
    checkType  string
    checkValue float64
}

func (i InterfaceDelta) String() string {
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

func (i InterfaceDelta) CheckValue() float64 {
    return i.checkValue
}

func (i InterfaceDelta) PerfData(warn, crit string) string {
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
        if cv.Delta > 0 {
            tmpItr := getInterfaceStats(cv.Name)
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

            deltaItr := InterfaceDelta{
                Name: itr.Name,
                OutTotal: ConvertToUnit(dOut, cv.Units()),
                OutPerSec: outPs,
                InTotal: ConvertToUnit(dIn, cv.Units()),
                InPerSec: inPs,
                Units: cv.Units(),
                checkType: cv.Against,
                checkValue: cVal,
            }

            // Save the current itr for later
            setInterfaceStats(itr)

            return deltaItr
        }

        return itr
    }

    return ifs
}

func getNetworkIfs() ([]Interface, error) {
    var ifList []Interface
    ifs, err := net.Interfaces()
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
