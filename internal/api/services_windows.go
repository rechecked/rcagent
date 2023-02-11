// +build windows

package api

import (
    "github.com/shirou/gopsutil/v3/winservices"
)

func getServices() ([]Service, error) {
    var srvs []Service
    services, err := winservices.ListServices()
    for _, s := range services {
        status := "stopped"
        if s.Status.Pid > 0 {
            status = "running"
        }
        srvs = append(srvs, Service{
            Name: s.Name,
            Status: status,
        })
    }
    return srvs, err
}

