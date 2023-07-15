package manager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/rechecked/rcagent/internal/config"
)

type HostInfo struct {
	Hostname string
	MachineId string
	OS string
	Platform string
}

// Set up the manager connection
func Run() {

	// Check if we should try to connect
	if config.Settings.Manager.APIKey == "" {
		return
	}

	c := time.Tick(1 * time.Minute)
	for range c {
		checkin()
	}

}

// Do inital registration when the agent starts up... send basic data and if we
// need to get a certificate we do that now.
func Register() {

	i := getHostInfo()
	data := map[string]string{
		"hostname": i.Hostname,
		"machineId": i.MachineId,
		"address": getOutboundIP(),
		"version": config.Version,
		"os": i.OS,
		"platform": i.Platform,
	}

	fmt.Println(data)

	err := sendPost("agents/register", data)
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

	fmt.Println(data)

	err := sendPost("agents/checkin", data)
	if err != nil {
		fmt.Println(err)
	}
}

// Send a POST request
func sendPost(path string, data map[string]string) error {

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
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-API-Key", config.Settings.Manager.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func getHostInfo() HostInfo {

	hostname, _ := os.Hostname()
	machineId, _ := machineid.ProtectedID("rcagent")
	host, _ := host.Info()

	i := HostInfo{
		Hostname: hostname,
		MachineId: machineId,
		OS: host.OS,
		Platform: host.Platform,
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
