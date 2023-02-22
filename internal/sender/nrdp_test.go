package sender

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
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

	mockResponse := `{"result":{"status":0,"message":"OK","meta":{"output":"0 checks processed"}}}`

	n := new(NRDPServer)
	err := n.SetConn("http://192.168.1.100/nrdp/", "TestToken")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Test connection that shouldn't work
	/*
	   if err := n.TestConn(); err == nil {
	       t.Log("TestConn() should error if it cannot connect")
	       t.Fail()
	   }
	*/

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://192.168.1.100/nrdp/",
		httpmock.NewStringResponder(200, mockResponse))

	// Test connection that should work
	if err := n.TestConn(); err != nil {
		t.Log(err)
		t.Fail()
	}

}

func TestNRDPSendCheckResults(t *testing.T) {

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
		func(req *http.Request) (*http.Response, error) {
			// Do something
			return httpmock.NewStringResponse(404, "Hello"), nil
		})

	n := new(NRDPServer)
	err := n.SetConn("http://192.168.1.100/nrdp/", "TestToken")
	resp, err := n.Send(mockChecks)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(resp.String())

}
