package config

import (
	"errors"
	"os"

	"github.com/kardianos/service"
)

var Log service.Logger

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

func UsingManager() bool {
	return Settings.Manager.APIKey != ""
}

func LogDebug(v ...interface{}) {
	if DebugMode {
		Log.Info(v...)
	}
}

func LogDebugf(format string, a ...interface{}) {
	if DebugMode {
		Log.Infof(format, a...)
	}
}
