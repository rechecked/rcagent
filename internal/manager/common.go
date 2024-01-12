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

type HostInfo struct {
	Hostname  string
	MachineId string
	OS        string
	Platform  string
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

	t := time.NewTicker(1 * time.Minute)
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

	_, err := sendPost("agents/register", data)
	if err != nil {
		config.Log.Error(err)
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

	if config.DebugMode {
		config.LogDebugf("Sending POST: %s", cfgUrl)
	}

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

	if config.DebugMode {
		config.LogDebugf("Sending GET: %s", cfgUrl)
	}

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

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
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

	// Add params to url if we need to
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
