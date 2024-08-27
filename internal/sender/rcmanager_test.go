package sender

import (
	"testing"
)

func TestRCManagerConnect(t *testing.T) {

	s := new(RCManagerServer)
	err := s.SetConn("testfail/rcmanager", "")
	if err == nil {
		t.Log("RCManager is not properly validating host/token")
		t.Fail()
	}

}
