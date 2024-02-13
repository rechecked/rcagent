package status

import (
	"testing"
)

func TestConversions(t *testing.T) {

	ret := ConvertToUnit(1000, "MB")
	if ret != 0.001 {
		t.Fail()
	}

	ret = ConvertToUnit(1024, "MiB")
	if ret != 0.0009765625 {
		t.Fail()
	}

	ret = ConvertToUnit(100, "B")
	if ret != 100 {
		t.Fail()
	}

	ret = ConvertToUnit(100, "HSDeB")
	if ret != 1e-16 {
		t.Fail()
	}

}
