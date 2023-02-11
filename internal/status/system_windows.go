// +build windows

package status

import (
    wapi "github.com/iamacarpet/go-win64api"
    so "github.com/iamacarpet/go-win64api/shared"
)

func getUsers() ([]so.SessionDetails, error) {
    users, err := wapi.ListLoggedInUsers()
    if err != nil {
        return users, err
    }
    return users, nil
}