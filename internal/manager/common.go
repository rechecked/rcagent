package manager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/rechecked/rcagent/internal/config"
)

const (
	DEFAULT_ACTIVIATION_WAIT_TIME   = 60  // 1 Minute
	DEFAULT_LIMIT_REACHED_WAIT_TIME = 600 // 10 Minutes
)

type HostInfo struct {
	Hostname  string
	MachineId string
	OS        string
	Platform  string
}

type RegisterStatus struct {
	Activated    bool `json:"activated"`
	LimitReached bool `json:"limitReached"`
}

type CheckInStatus struct {
	NeedsConfigUpdate  bool `json:"needsConfigUpdate"`
	NeedsSecretsUpdate bool `json:"needsSecretsUpdate"`
}

type ConfigsData struct {
	Configs map[string]string
	Senders string
	Plugins map[string]string
}

// Set up the manager connection
func Run(restart chan<- struct{}) {

	// Check if we should try to connect
	if !config.UsingManager() {
		return
	}

	t := time.NewTicker(20 * time.Second)
	defer t.Stop()
	for range t.C {
		checkin()
		validateCert(restart)
	}

}

// If we are connected to the manager, we need to do config parsing once before
// we start this to ensure we have the most up to date configs - the manager sync will
// automatically update configs while it's running
func Init() {
	if config.UsingManager() {
		if Sync(true, true) {
			config.ParseConfigDir()
		}
	}
}

func Sync(s, c bool) bool {
	nu := false
	if s {
		if updateSecrets() {
			nu = true
		}
	}
	if c {
		if updateConfigs() {
			nu = true
		}
	}
	return nu
}

// Do inital registration when the agent starts up... send basic data and if we
// need to get a certificate we do that now.
func Register() {

	// Skip registration if we aren't using the manager
	if !config.UsingManager() {
		return
	}

	if config.DebugMode {
		url, _ := getManagerUrl("", nil)
		config.LogDebugf("Registering with RCM (%s)", url)
	}

	i := getHostInfo()
	data := map[string]string{
		"hostname":  i.Hostname,
		"machineId": i.MachineId,
		"address":   getOutboundIP(),
		"version":   config.Version,
		"os":        i.OS,
		"platform":  i.Platform,
		"token":     config.Settings.Token,
	}

	resp, err := sendPost("agents/register", data)
	if err != nil {
		config.Log.Error(err)
	}

	// Validate that it was registered, or mention that it is waiting approval
	if len(resp) > 0 {
		r := RegisterStatus{}
		err = json.Unmarshal(resp, &r)
		if err != nil {
			config.Log.Error(err)
			return
		}

		// If we are waiting for approval, then wait to try again
		if !r.Activated {
			config.Log.Info("ReChecked Manager Activation Required: Waiting for registration approval in ReChecked Manager")
			time.Sleep(time.Second * DEFAULT_ACTIVIATION_WAIT_TIME)
			Register()
		}

		// If the org's agent limit is reached, we can do the same as above but with limit message
		if r.LimitReached {
			config.Log.Info("ReChecked Manager Agents Limit Reached: Cannot register with ReChecked Manager")
			time.Sleep(time.Second * DEFAULT_LIMIT_REACHED_WAIT_TIME)
			Register()
		}
	}
}

func GetMachineId() string {
	return getHostInfo().MachineId
}

// Send some basic data to the manager to "check in" with it, indicating
// that the agent is running, accessible, and provides feedback on current status
func checkin() {

	data := map[string]string{
		"machineId": getHostInfo().MachineId,
	}

	b, err := sendPost("agents/checkin", data)
	if err != nil {
		config.Log.Error(err)
		return
	}

	c := CheckInStatus{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		config.Log.Error(err)
		return
	}

	// Sync certain things if they need to be synced
	if Sync(c.NeedsSecretsUpdate, c.NeedsConfigUpdate) {
		config.ParseConfigDir()
	}
}

// Send a POST request
func sendPost(path string, data map[string]string) ([]byte, error) {

	cfgUrl, err := getManagerUrl(path, nil)
	if err != nil {
		return []byte{}, err
	}

	config.LogDebugf("Sending POST: %s", cfgUrl)

	postBody, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", cfgUrl, bytes.NewBuffer(postBody))
	if err != nil {
		return []byte{}, err
	}

	return getRequest(req)
}

// Send a GET request
func sendGet(path string, params url.Values) ([]byte, error) {

	cfgUrl, err := getManagerUrl(path, params)
	if err != nil {
		return []byte{}, err
	}

	config.LogDebugf("Sending GET: %s", cfgUrl)

	req, err := http.NewRequest("GET", cfgUrl, nil)
	if err != nil {
		return []byte{}, err
	}

	return getRequest(req)
}

func downloadFile(name string, url string) error {

	// Create or truncate file
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	if config.DebugMode {
		config.LogDebugf("Downloading: %s", url)
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func getRequest(req *http.Request) ([]byte, error) {

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-API-Key", config.Settings.Manager.APIKey)

	// Set up an HTTP client, ignore invalid certs if we have the config option set
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Settings.Manager.IgnoreCert},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode == http.StatusOK {
		return bodyBytes, err
	}

	return []byte{}, nil
}

func getManagerUrl(path string, params url.Values) (string, error) {

	// Make sure we have a proper url, default to cloud RCM if no manager URL specified
	cfgUrl := config.Settings.Manager.Url
	if cfgUrl == "" {
		cfgUrl = "https://manage.rechecked.io/api"
	}

	url, err := url.Parse(cfgUrl)
	if err != nil {
		return "", err
	}

	url = url.JoinPath(path)

	// Add params to url if we need to and make sure they are url encoded
	if len(params) > 0 {
		url.RawQuery = params.Encode()
	}

	return url.String(), nil
}

func getHostInfo() HostInfo {

	hostname, _ := os.Hostname()
	machineId, _ := machineid.ProtectedID("rcagent")
	host, _ := host.Info()

	i := HostInfo{
		Hostname:  hostname,
		MachineId: machineId,
		OS:        host.OS,
		Platform:  host.Platform,
	}

	return i
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		config.Log.Error(err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
