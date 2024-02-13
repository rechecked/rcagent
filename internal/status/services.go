package status

import (
	"fmt"

	"github.com/rechecked/rcagent/internal/config"
)

type Service struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	expected string
}

func (s Service) String() string {
	return fmt.Sprintf("%s is [%s] (expected value is [%s])", s.Name, s.Status, s.expected)
}

func (s Service) CheckValue() string {
	return s.Status
}

func HandleServices(cv config.Values) interface{} {
	svcs, err := getServices()
	if err != nil {
		return []string{}
	}

	// Filter services
	if cv.Check && cv.Against != "" {
		for _, s := range svcs {
			if s.Name == cv.Against {
				s.expected = cv.Expected
				return s
			}
		}
	}

	return svcs
}
