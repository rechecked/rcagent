package manager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
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
	NeedsConfigUpdate bool `json:"needsConfigUpdate"`
	NeedsSecretsUpdate bool `json:"needsSecretsUpdate"`
}

// Set up the manager connection
func Run(restart chan<- struct{}) {

	// Check if we should try to connect
	if config.Settings.Manager.APIKey == "" {
		return
	}

	c := time.Tick(1 * time.Minute)
	for range c {
		checkin()
		validateCert(restart)
	}

}

// Do inital registration when the agent starts up... send basic data and if we
// need to get a certificate we do that now.
func Register() {

	i := getHostInfo()
	data := map[string]string{
		"hostname":  i.Hostname,
		"machineId": i.MachineId,
		"address":   getOutboundIP(),
		"version":   config.Version,
		"os":        i.OS,
		"platform":  i.Platform,
		"token": config.Settings.Token,
	}

	_, err := sendPost("agents/register", data)
	if err != nil {
		fmt.Println(err)
	}
}

// Send some basic data to the manager to "check in" with it, indicating
// that the agent is running, accessible, and provides feedback on current status
func checkin() {

	i := getHostInfo()
	data := map[string]string{
		"machineId": i.MachineId,
	}

	b, err := sendPost("agents/checkin", data)
	if err != nil {
		fmt.Println(err)
	}

	c := CheckInStatus{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		fmt.Println(err)
	}

	// Get secrets if they need updating
	if c.NeedsSecretsUpdate {
		updateSecrets()
	}

	if c.NeedsConfigUpdate {
		updateConfigs()
	}

}

func updateSecrets() {

	

}

func updateConfigs() {

}

// Send a POST request
func sendPost(path string, data map[string]string) ([]byte, error) {

	// Make sure we have a proper url, default to manage.rechecked.io if empty
	url := config.Settings.Manager.Url
	if url == "" {
		url = "https://manage.rechecked.io"
	}

	if url[len(url)-1:] != "/" {
		url += "/"
	}
	url += path

	// Set up an HTTP client, ignore invalid certs if we have the config option set
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Settings.Manager.IgnoreCert},
	}
	client := &http.Client{Transport: tr}

	postBody, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-API-Key", config.Settings.Manager.APIKey)

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
		fmt.Println(err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
