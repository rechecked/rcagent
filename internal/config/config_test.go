package config

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestReplaceVariables(t *testing.T) {

	// Set a list of secrets to use
	CfgData.Secrets = map[string]string{
		"TEST_SECRET": "123456",
		"HOST":        "host123",
	}

	cfgFile := []byte(`some config
	test config: hello
	test config host: $HOST
	test config secret: $TEST_SECRET
	test config no secret: TEST_SECRET
	testboth: $TEST_SECRET $HOST`)

	replaceVariables(&cfgFile)

	assert.Equal(t, `some config
	test config: hello
	test config host: host123
	test config secret: 123456
	test config no secret: TEST_SECRET
	testboth: 123456 host123`,
		string(cfgFile),
		"Config variable replacements should match")
}
