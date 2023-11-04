package sender

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestNRDPBadCreate(t *testing.T) {

	n := new(NRDPServer)
	err := n.SetConn("testfail/nrdp", "")
	if err == nil {
		t.Log("NRDPServer is not properly validating host/token")
		t.Fail()
	}

}

func TestNRDPConnect(t *testing.T) {

	n := new(NRDPServer)
	err := n.SetConn("http://192.168.1.100/nrdp/", "TestToken")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Test connection that should work
	mockResponse := `{"result":{"status":0,"message":"OK","meta":{"output":"0 checks processed"}}}`
	httpmock.RegisterResponder("POST", "http://192.168.1.100/nrdp/",
		httpmock.NewStringResponder(200, mockResponse))
	if err := n.TestConn(); err != nil {
		t.Log(err)
		t.Fail()
	}

	// Test connection that cannot validate response
	mockResponse = `{}`
	httpmock.Reset()
	httpmock.RegisterResponder("POST", "http://192.168.1.100/nrdp/",
		httpmock.NewStringResponder(404, mockResponse))
	if err := n.TestConn(); err == nil {
		t.Fail()
	}

	// Test connection that has bad JSON repsonse
	mockResponse = ``
	httpmock.Reset()
	httpmock.RegisterResponder("POST", "http://192.168.1.100/nrdp/",
		httpmock.NewStringResponder(404, mockResponse))
	if err := n.TestConn(); err == nil {
		t.Fail()
	}

}

func TestNRDPSendCheckResults(t *testing.T) {

	mockResponse := `{"result":{"status":0,"message":"OK","meta":{"output":"2 checks processed"}}}`
	mockChecks := []NRDPCheckResult{
		{
			Checkresult: NRDPObjectType{
				Type: "host",
			},
			Hostname: "Test Host",
			State:    0,
			Output:   "Test output for test check 1",
		},
		{
			Checkresult: NRDPObjectType{
				Type: "service",
			},
			Hostname:    "Test Host",
			Servicename: "Test Service",
			State:       0,
			Output:      "Test output for test check 2",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://192.168.1.100/nrdp/",
		httpmock.NewStringResponder(200, mockResponse))

	n := new(NRDPServer)
	if err := n.SetConn("http://192.168.1.100/nrdp/", "TestToken"); err != nil {
		t.Log(err)
		t.Fail()
	}

	resp, err := n.Send(mockChecks)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Test sender string
	if n.String() != "http://192.168.1.100/nrdp/" {
		t.Log("String mismatch: ", n.String())
		t.Fail()
	}

	// Test response string
	if resp.String() != "Status: 0 | Message: OK | Meta Output: 2 checks processed" {
		t.Log("String mismatch: ", resp.String())
		t.Fail()
	}

}
