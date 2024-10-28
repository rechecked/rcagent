package sender

import (
	"testing"
)

func TestRCManagerConnect(t *testing.T) {

	s := RCManagerSender{}
	err := s.SetConn("testfail/rcmanager", "")
	if err == nil {
		t.Log("RCManagerSender is not properly validating host/token")
		t.Fail()
	}

}

func TestRCManagerCheckResults(t *testing.T) {

	s := RCManagerSender{}
	if err := s.SetConn("http://192.168.1.100/api/", "TestToken"); err != nil {
		t.Log(err)
		t.Fail()
	}

}
