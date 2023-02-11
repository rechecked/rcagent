// +build !windows

package status

import (
    "strings"
    "github.com/go-cmd/cmd"
)

func getServices() ([]Service, error) {

    // Parse systemctl (or whatever other type of system service manager they have)
    // systemctl list-units --type=service --all --plain

    svcs := []Service{}
    c := cmd.NewCmd("systemctl", "list-units", "--type=service", "-all", "--plain")
    s := <-c.Start()

    if len(s.Stdout) > 0 {
        for _, l := range s.Stdout {
            if !strings.Contains(l, ".service") {
                continue
            }
            tmp := strings.Fields(l)
            if len(tmp) >= 4 {
                svcs = append(svcs, Service{
                    Name: strings.ReplaceAll(tmp[0], ".service", ""),
                    Status: tmp[3],
                })
            }
        }
    }

    return svcs, nil

}