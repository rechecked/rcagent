package sender

import (
	"testing"
)

func TestRCManagerConnect(t *testing.T) {

	s := new(RCManagerServer)
	err := s.SetConn("testfail/nrdp", "")
	if err == nil {
		t.Log("NRDPServer is not properly validating host/token")
		t.Fail()
	}

}
