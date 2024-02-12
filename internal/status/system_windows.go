//go:build windows
// +build windows

package status

import (
	wapi "github.com/iamacarpet/go-win64api"
	so "github.com/iamacarpet/go-win64api/shared"

	"github.com/rechecked/rcagent/internal/config"
)

func getUsers() ([]so.SessionDetails, error) {
	users, err := wapi.ListLoggedInUsers()
	if err != nil {
		config.Log.Error(err)
		return users, err
	}
	return users, nil
}
