// +build darwin

package services

import (
    "strings"
    "github.com/go-cmd/cmd"
)

func getServices() ([]Service, error) {

    // Parse launchctl for macOS systems
    // launchctl list

    svcs := []Service{}
    c := cmd.NewCmd("launchctl", "list")
    s := <-c.Start()

    if len(s.Stdout) > 0 {
        for _, l := range s.Stdout {
            tmp := strings.Fields(l)
            status := "stopped"
            if tmp[1] == "-" {
                status = "running"
            }
            svcs = append(svcs, Service{
                Name: tmp[3],
                Status: status,
            })
        }
    }

    return svcs, nil

}