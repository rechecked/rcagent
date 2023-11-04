package sender

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type NRDPServer struct {
	Name  string
	Url   string
	Token string
}

type NRDPResponse struct {
	Status  int
	Message string
	Output  string
}

type NRDPObjectType struct {
	Type string `json:"type"`
}

type NRDPCheckResult struct {
	Checkresult NRDPObjectType `json:"checkresult"`
	Hostname    string         `json:"hostname"`
	Servicename string         `json:"servicename,omitempty"`
	State       int            `json:"state"`
	Output      string         `json:"output"`
}

// Create a new NRDPServer and verify the url
func (n *NRDPServer) SetConn(u, token string) error {

	if _, err := url.ParseRequestURI(u); err != nil {
		return err
	}
	n.Url = u
	n.Token = token

	return nil
}

// Send a request to the NRDP server with any check data we want to pass
func (n *NRDPServer) Send(checks []NRDPCheckResult) (NRDPResponse, error) {

	// Create string of json for NRDP
	res := "[]"
	if len(checks) > 0 {
		results, err := json.Marshal(checks)
		if err != nil {
			return NRDPResponse{}, err
		}
		res = string(results)
	}

	data := url.Values{
		"cmd":   {"submitcheck"},
		"token": {n.Token},
		"json":  {fmt.Sprintf(`{"checkresults":%s}`, res)},
	}

	resp, err := sendToNRDP(n.Url, data)
	if err != nil {
		return NRDPResponse{}, err
	}

	return resp, nil
}

// Check if NRDP server and creds are valid and return and error
// if they are not...
func (n *NRDPServer) TestConn() error {

	data := url.Values{
		"cmd":   {"submitcheck"},
		"token": {n.Token},
		"json":  {`{"checkresults":[]}`},
	}

	resp, err := http.PostForm(n.Url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	status := gjson.GetBytes(body, "result.status")
	message := gjson.GetBytes(body, "result.message")
	if status.Int() == 0 && message.String() == "OK" {
		return nil
	}

	return errors.New("could not validate message")
}

func (n *NRDPServer) String() string {
	return n.Url
}

func (n *NRDPResponse) String() string {
	return fmt.Sprintf("Status: %d | Message: %s | Meta Output: %s", n.Status, n.Message, n.Output)
}

func sendToNRDP(url string, data url.Values) (NRDPResponse, error) {

	resp, err := http.PostForm(url, data)
	if err != nil {
		return NRDPResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NRDPResponse{}, err
	}

	status := gjson.GetBytes(body, "result.status")
	message := gjson.GetBytes(body, "result.message")
	output := gjson.GetBytes(body, "result.meta.output")

	n := NRDPResponse{
		Status:  int(status.Int()),
		Message: message.String(),
		Output:  output.String(),
	}

	return n, nil
}
