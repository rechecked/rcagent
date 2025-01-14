package endpoints

import (
	"errors"
	"fmt"

	"github.com/rechecked/rcagent/internal/config"
	"github.com/rechecked/rcagent/internal/status"
)

func GetDataFromEndpoint(path string, values config.Values) (interface{}, error) {
	endpoint := config.Endpoints[path]
	if endpoint != nil {

		// Get the data back from the endpoint
		e := endpoint(values)

		// Check if we are checkable type
		chk, ok := e.(status.Checkable)
		if values.Check && ok {
			check := status.GetCheckResult(chk, values.Warning, values.Critical)
			return check, nil
		}

		// Check if we are a checkable against type
		chk2, ok2 := e.(status.CheckableAgainst)
		if values.Check && ok2 {
			check := status.GetCheckAgainstResult(chk2, values.Expected)
			return check, nil
		}

		// If we have an error we need to format it for output
		if err, ok := e.(error); ok {
			if values.Check {
				e = status.CheckResult{
					Output:   fmt.Sprintf("Error: %s", err.Error()),
					Exitcode: status.STATUS_ERROR,
				}
			} else {
				e = struct {
					Error string `json:"error"`
				}{
					Error: err.Error(),
				}
			}
		}

		// If we aren't doing a check, just return value
		return e, nil
	}
	return nil, errors.New("GetDataFromEndpoint: Endpoint does not exist")
}
