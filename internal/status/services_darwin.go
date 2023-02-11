// +build darwin

package status

import (
    "fmt"
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
            fmt.Printf("%s\n", tmp)
            status := "stopped"
            if tmp[0] != "-" {
                status = "running"
            }
            if len(tmp) >= 3 {
                svcs = append(svcs, Service{
                    Name: tmp[2],
                    Status: status,
                })
            }
        }
    }

    return svcs, nil

}