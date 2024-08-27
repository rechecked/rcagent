package sender

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RCManagerServer struct {
	Url    string
	APIKey string
}

type RCManagerResponse struct {
}

// Create a new RCManagerServer and verify the url
func (n *RCManagerServer) SetConn(u, token string) error {

	if _, err := url.ParseRequestURI(u); err != nil {
		return err
	}
	n.Url = u

	return nil
}

// Send a request to the NRDP server with any check data we want to pass
func (r *RCManagerServer) Send(checks []NRDPCheckResult) error {

	// Create list of
	res := "[]"
	if len(checks) > 0 {
		results, err := json.Marshal(checks)
		if err != nil {
			return err
		}
		res = string(results)
	}

	data := url.Values{
		"cmd":   {"submitcheck"},
		"token": {r.APIKey},
		"json":  {fmt.Sprintf(`{"checkresults":%s}`, res)},
	}

	err := sendToRCManager(r.Url, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *RCManagerServer) String() string {
	return r.Url
}

func sendToRCManager(url string, data url.Values) error {

	resp, err := http.PostForm(url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(body)

	return nil
}
