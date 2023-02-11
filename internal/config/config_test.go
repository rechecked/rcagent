
package config

import (
    "testing"
    "net/http"
    "net/url"
)

func TestOverrides(t *testing.T) {

    // Fake request with units set to GB
    uri := "status/?"
    param := make(url.Values)
    param["units"] = []string{"GB"}
    req, err := http.NewRequest("GET", uri+param.Encode(), nil)
    if err != nil {
        t.Log(err)
        t.Fail()
    }

    // Settings with default units set
    Settings.Units = "B"

    v := ParseValues(req)

    // If units does not show up as GB then failed

    if v.Units() != "GB" {
        t.Log("Units were not properly applied to config values.")
        t.Fail()
    }

}
