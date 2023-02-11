
package api

import (
	"github.com/shirou/gopsutil/v3/net"
	"github.com/rechecked/rcagent/internal/config"
)

// Get a whole list of interfaces
func HandleNetworks(cv config.Values) interface{} { 
	ifs, err := net.Interfaces()
	if err != nil {
		return nil
	}
	return ifs
}

