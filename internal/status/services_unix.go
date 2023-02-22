//go:build !windows && !darwin
// +build !windows,!darwin

package status

import (
	"github.com/go-cmd/cmd"
	"strings"
)

func getServices() ([]Service, error) {

	// Parse systemctl (or whatever other type of system service manager they have)
	// systemctl list-units --type=service --all --plain --no-pager --no-legend

	svcs := []Service{}
	c := cmd.NewCmd("systemctl", "list-units", "--type=service", "--all", "--plain", "--no-pager", "--no-legend")
	s := <-c.Start()

	if len(s.Stdout) > 0 {
		for _, l := range s.Stdout {
			tmp := strings.Fields(l)
			if len(tmp) >= 4 {
				svcs = append(svcs, Service{
					Name:   strings.ReplaceAll(tmp[0], ".service", ""),
					Status: tmp[3],
				})
			}
		}
	}

	return svcs, nil

}
