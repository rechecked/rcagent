package config

import (
	"net/http"
	"net/url"
	"testing"
)

func TestOverrides(t *testing.T) {

	// Settings with default units set
	Settings.Units = "B"

	// Fake request with units set to GB
	uri := "status/?"
	param := make(url.Values)
	param["units"] = []string{"GB"}
	req, err := http.NewRequest("GET", uri+param.Encode(), nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// If units does not show up as GB then failed
	v := ParseValues(req)
	if v.Units() != "GB" {
		t.Log("Units were not properly applied to config values")
		t.Fail()
	}

	// Test default is always set to "B" even if not configured
	Settings.Units = ""
	req, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	v = ParseValues(req)
	if v.Units() != "B" {
		t.Log("Units should default to B")
		t.Fail()
	}

}
