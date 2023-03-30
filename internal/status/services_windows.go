//go:build windows
// +build windows

package status

import (
	"github.com/shirou/gopsutil/v3/winservices"
	"golang.org/x/sys/windows"
)

func getServices() ([]Service, error) {
	var srvs []Service
	services, err := winservices.ListServices()
	for _, s := range services {
		status := "unknown"
		// Need to create a new service before querying for status
		service, err := winservices.NewService(s.Name)
		if err != nil {
			continue
		}
		qs, err := service.QueryStatus()
		if err != nil {
			continue
		}
		// Check status is running or not from returned ServiceStatus
		if qs.State == windows.SERVICE_RUNNING {
			status = "running"
		} else if qs.State == windows.SERVICE_STOPPED {
			status = "stopped"
		}
		srvs = append(srvs, Service{
			Name:   s.Name,
			Status: status,
		})
	}
	return srvs, err
}
