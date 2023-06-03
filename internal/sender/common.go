package sender

import (
	"fmt"
	"github.com/rechecked/rcagent/internal/config"
	"github.com/rechecked/rcagent/internal/server"
	"github.com/rechecked/rcagent/internal/status"
	"time"
)

// Save passive checks where we can edit them
var Checks []config.CheckCfg

// Set up passive related loop
func Run() {

	// Verify we have checks to run...
	Checks = config.Settings.PassiveChecks
	if len(Checks) == 0 {
		if config.Settings.Debug {
			fmt.Print("Stopping senders: No checks configured\n")
		}
		return
	}

	// Verify we have senders to send to...
	if len(config.Settings.Senders) == 0 {
		if config.Settings.Debug {
			fmt.Print("Stopping senders: No senders configured\n")
		}
		return
	}

	// Tick every second and check for any passive checks we need
	// to send, then process them and continue... until program exit
	c := time.Tick(1 * time.Second)
	for range c {
		now := time.Now()
		for i, check := range Checks {
			if check.Disabled || check.NextRun.After(now) {
				continue
			}

			// Parse check interval duration and error if it is bad and disable,
			// then set next run time if it isn't disabled
			dur, err := time.ParseDuration(check.Interval)
			if err != nil {
				if config.Settings.Debug {
					fmt.Printf("Interval Error: %s\n", err)
				}
				fmt.Printf("The interval for '%s - %s' is invalid, disabling",
					check.Hostname, check.Servicename)
				Checks[i].Disabled = true
			}
			Checks[i].NextRun = now.Add(dur)

			// Run the check and get the value data back, try to send it off if we can
			data, err := server.GetDataFromEndpoint(check.Endpoint, check.Options)
			if err != nil {
				fmt.Printf("Check Error: %s\n", err)
			}
			chk, ok := data.(status.CheckResult)
			if ok {
				go sendToSenders(chk, check)
				if config.Settings.Debug {
					fmt.Printf("%s\n", chk.String())
				}
			} else {
				fmt.Printf("The check for '%s - %s' is invalid, check endpoints and options, disabling",
					check.Hostname, check.Servicename)
				Checks[i].Disabled = true
			}
		}
	}
}

func sendToSenders(chk status.CheckResult, cfg config.CheckCfg) {
	if config.Settings.Debug {
		fmt.Printf("Sending check: %s\n", chk.String())
	}

	// Get all senders
	senders := config.Settings.Senders

	for _, sender := range senders {
		// We only have NRDP for now but more later?
		if sender.Type == "nrdp" {
			s := new(NRDPServer)
			err := s.SetConn(sender.Url, sender.Token)
			if err != nil {
				fmt.Printf("Error: sendToSenders: %s", err)
			}

			// Set output
			output := chk.Output
			if chk.LongOutput != "" {
				output = fmt.Sprintf("%s\n%s", output, chk.LongOutput)
			}
			if chk.Perfdata != "" {
				output = fmt.Sprintf("%s | %s", output, chk.Perfdata)
			}

			// Create the nrdp result
			var checks []NRDPCheckResult
			if cfg.Servicename != "" {
				checks = []NRDPCheckResult{
					{
						Checkresult: NRDPObjectType{
							Type: "service",
						},
						Hostname:    cfg.Hostname,
						Servicename: cfg.Servicename,
						State:       chk.Exitcode,
						Output:      output,
					},
				}
			} else {
				checks = []NRDPCheckResult{
					{
						Checkresult: NRDPObjectType{
							Type: "host",
						},
						Hostname: cfg.Hostname,
						State:    chk.Exitcode,
						Output:   output,
					},
				}
			}
			resp, err := s.Send(checks)
			if err != nil {
				fmt.Printf("Error: sendToSenders: %s", err)
			}
			if config.Settings.Debug {
				fmt.Printf("Sender NRDP repsonse: %s", resp.String())
			}
		}
	}
}
