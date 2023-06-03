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

	"github.com/rechecked/rcagent/internal/config"
)

// Set up the manager connection
func Run() {

	// Check if we should try to connect
	if config.Settings.Manager.APIKey == "" {
		return
	}

	checkin()
	c := time.Tick(1 * time.Minute)
	for range c {
		checkin()
	}

}

// Send some basic data to the manager to "check in" with it, indicating
// that the agent is running, accessible, and provides feedback on current status
func checkin() {

	hostname, _ := os.Hostname()
	machineId, _ := machineid.ProtectedID("rcagent")

	data := map[string]string{
		"hostname": hostname,
		"machineId": machineId,
	}

	fmt.Println(data)

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
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(postBody))
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
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
