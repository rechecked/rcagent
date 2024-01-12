package sender

import (
	"fmt"
	"time"

	"github.com/rechecked/rcagent/internal/config"
	"github.com/rechecked/rcagent/internal/manager"
	"github.com/rechecked/rcagent/internal/server"
	"github.com/rechecked/rcagent/internal/status"
)

// Set up passive related loop
func Run() {

	manager.Init()

	// Verify we have checks to run...
	if len(config.CfgData.Checks) == 0 {
		config.LogDebug("SENDER: No checks configured")
	} else {
		config.LogDebugf("SENDER: %d checks configured", len(config.CfgData.Checks))
	}

	// Verify we have senders to send to...
	if len(config.CfgData.Senders) == 0 {
		config.LogDebug("SENDER: No senders configured")
	} else {
		config.LogDebugf("SENDER: %d senders configured", len(config.CfgData.Senders))
	}

	// Tick every second and check for any passive checks we need
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for range t.C {
		config.CfgData.RLock()
		runChecks()
		config.CfgData.RUnlock()
	}
}

func runChecks() {
	now := time.Now()
	for i, check := range config.CfgData.Checks {
		if check.Disabled || check.NextRun.After(now) {
			continue
		}

		// Parse check interval duration and error if it is bad and disable,
		// then set next run time if it isn't disabled
		dur, err := time.ParseDuration(check.Interval)
		if err != nil {
			config.LogDebugf("Interval Error: %s\n", err)
			config.Log.Infof("The interval for '%s' is invalid, disabling", check.Name())
			config.CfgData.Checks[i].Disabled = true
		}
		config.CfgData.Checks[i].NextRun = now.Add(dur)

		// Run the check and get the value data back, try to send it off if we can
		data, err := server.GetDataFromEndpoint(check.Endpoint, check.Options)
		if err != nil {
			config.Log.Infof("Check Error: %s\n", err)
		}
		chk, ok := data.(status.CheckResult)
		if ok {
			go sendToSenders(chk, check)
			config.LogDebugf("%s\n", chk.String())
		} else {
			config.LogDebug(data)
			config.Log.Infof("The check for '%s' is invalid, check endpoints and options, disabling",
				check.Name())
			config.CfgData.Checks[i].Disabled = true
		}
	}
}

func sendToSenders(chk status.CheckResult, cfg config.CheckCfg) {

	config.LogDebugf("Sending check: %s\n", chk.String())

	// Get all senders
	config.CfgData.RLock()
	senders := config.CfgData.Senders
	config.CfgData.RUnlock()

	for _, sender := range senders {
		// We only have NRDP for now but more later?
		if sender.Type == "nrdp" {
			s := new(NRDPServer)
			err := s.SetConn(sender.Url, sender.Token)
			if err != nil {
				config.Log.Errorf("Error: sendToSenders: %s", err)
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
				config.Log.Errorf("Error: sendToSenders: %s", err)
			}
			config.LogDebugf("Sender NRDP repsonse: %s", resp.String())
		}
	}
}
