package sender

import (
	"net/url"

	"github.com/rechecked/rcagent/internal/config"
	"github.com/rechecked/rcagent/internal/manager"
	"github.com/rechecked/rcagent/internal/status"
)

type RCManagerSender struct {
	Url    string
	APIKey string
	Checks []RCManagerCheck
}

type RCManagerCheck struct {
	CheckId    int64  `json:"checkId"`
	ExitCode   int64  `json:"exitCode"`
	Output     string `json:"output"`
	LongOutput string `json:"longOutput"`
	PerfData   string `json:"perfData"`
}

type RCManagerResponse struct {
}

// Create a new RCManagerSender and verify the url
func (r *RCManagerSender) SetConn(u, t string) error {

	if _, err := url.ParseRequestURI(u); err != nil {
		return err
	}
	r.Url = u
	r.APIKey = t

	return nil
}

func (r *RCManagerSender) TestConn() error {

	params := url.Values{}
	res, err := manager.SendGet("", params)
	if string(res) == "" {
		return nil
	}

	return err
}

func (r *RCManagerSender) Format(cfg config.CheckCfg, chk status.CheckResult) {
	var checks []RCManagerCheck

	check := RCManagerCheck{
		CheckId:    cfg.CheckId,
		ExitCode:   int64(chk.Exitcode),
		Output:     chk.Output,
		LongOutput: chk.LongOutput,
		PerfData:   chk.Perfdata,
	}

	checks = append(checks, check)

	r.Checks = checks
}

// Send a request to the NRDP server with any check data we want to pass
func (r *RCManagerSender) Send() error {

	if len(r.Checks) == 0 {
		return nil
	}

	_, err := manager.SendPost("agents/checks/submit", r.Checks)
	if err != nil {
		return err
	}

	return nil
}

func (r *RCManagerSender) String() string {
	return r.Url
}
