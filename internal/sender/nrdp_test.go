package sender

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/rechecked/rcagent/internal/config"
	"github.com/rechecked/rcagent/internal/status"
)

func TestNRDPBadCreate(t *testing.T) {

	n := NRDPSender{}
	err := n.SetConn("testfail/nrdp", "")
	if err == nil {
		t.Log("NRDPSender is not properly validating host/token")
		t.Fail()
	}

}

func TestNRDPConnect(t *testing.T) {

	n := NRDPSender{}
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

func TestNRDPFormatCheckResults(t *testing.T) {

	hostCfg := config.CheckCfg{
		Hostname: "Test Host",
	}
	hostStatus := status.CheckResult{
		Exitcode:   0,
		Output:     "Test output",
		Perfdata:   "x=10",
		LongOutput: "extra data",
	}
	host := []NRDPCheckResult{
		{
			Checkresult: NRDPObjectType{
				Type: "host",
			},
			Hostname: "Test Host",
			State:    0,
			Output:   "Test output\nextra data | x=10",
		},
	}

	n := NRDPSender{}
	n.Format(hostCfg, hostStatus)

	if !reflect.DeepEqual(n.Checks, host) {
		t.Fail()
	}

	srvCfg := config.CheckCfg{
		Hostname:    "Test Host",
		Servicename: "Test Service",
	}
	srvStatus := status.CheckResult{
		Exitcode: 0,
		Output:   "Test Output",
	}
	srv := []NRDPCheckResult{
		{
			Checkresult: NRDPObjectType{
				Type: "service",
			},
			Hostname:    "Test Host",
			Servicename: "Test Service",
			State:       0,
			Output:      "Test output",
		},
	}

	n = NRDPSender{}
	n.Format(srvCfg, srvStatus)

	if !reflect.DeepEqual(n.Checks, srv) {
		t.Fail()
	}

	fmt.Println(n.Checks)

}

func TestNRDPSendCheckResults(t *testing.T) {

	mockResponse := `{"result":{"status":0,"message":"OK","meta":{"output":"2 checks processed"}}}`

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://192.168.1.100/nrdp/",
		httpmock.NewStringResponder(200, mockResponse))

	n := NRDPSender{}
	n.Checks = []NRDPCheckResult{
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

	if err := n.SetConn("http://192.168.1.100/nrdp/", "TestToken"); err != nil {
		t.Log(err)
		t.Fail()
	}

	err := n.Send()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Test sender string
	if n.String() != "http://192.168.1.100/nrdp/" {
		t.Log("String mismatch: ", n.String())
		t.Fail()
	}

}
