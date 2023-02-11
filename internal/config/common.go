
package config

import (
    "os"
    "errors"
)

func FileExists(file string) bool {
    _, err := os.Stat(file)
    return !errors.Is(err, os.ErrNotExist)
}

func Contains(s []string, val string) bool {
    for _, v := range s {
        if val == v {
            return true
        }
    }
    return false
}