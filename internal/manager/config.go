package manager

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rechecked/rcagent/internal/config"
)

var sanRegex = regexp.MustCompile(`[^a-z0-9_]+`)

func updateSecrets() bool {

	params := url.Values{}
	params.Add("machineId", getHostInfo().MachineId)

	json, err := sendGet("agents/update/secrets", params)
	if err != nil {
		config.Log.Error(err)
		return false
	}

	f := config.GetConfigDirFilePath("manager/secrets.json")

	// If there are no secrets then return now and don't add a secrets file
	if len(json) == 0 {
		os.Remove(f)
		return true
	}

	// Make sure the directory exists
	os.MkdirAll(config.GetConfigDirFilePath("manager"), 0755)

	if err := os.WriteFile(f, json, 0666); err != nil {
		config.Log.Error(err)
		return false
	}

	data := map[string]string{
		"machineId": getHostInfo().MachineId,
	}

	_, err = sendPost("agents/update/secrets", data)
	if err != nil {
		config.Log.Error(err)
		return false
	}

	return true
}

func updateConfigs() bool {

	params := url.Values{}
	params.Add("machineId", getHostInfo().MachineId)

	rawJSON, err := sendGet("agents/update/configs", params)
	if err != nil {
		config.Log.Error(err)
		return false
	}

	c := ConfigsData{}
	err = json.Unmarshal(rawJSON, &c)
	if err != nil {
		config.Log.Error(err)
		return false
	}

	// Make sure the directory exists
	os.MkdirAll(config.GetConfigDirFilePath("manager"), 0755)

	// Remove all config files from manager directory, this is to make sure
	// we are synced with the manager's configs
	files, err := filepath.Glob(config.GetConfigDirFilePath("manager/*.yml"))
	if err != nil {
		config.Log.Error(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			config.Log.Error(err)
		}
	}

	if c.Senders != "" {
		f := config.GetConfigDirFilePath("manager/senders.yml")
		if err := os.WriteFile(f, []byte(c.Senders), 0666); err != nil {
			config.Log.Error(err)
		}
	}

	if len(c.Configs) > 0 {
		for cfg, cfgStr := range c.Configs {
			f := config.GetConfigDirFilePath("manager/" + SanatizeFilename(cfg) + ".yml")
			if err := os.WriteFile(f, []byte(cfgStr), 0666); err != nil {
				config.Log.Error(err)
			}
		}
	}

	data := map[string]string{
		"machineId": getHostInfo().MachineId,
	}

	_, err = sendPost("agents/update/configs", data)
	if err != nil {
		config.Log.Error(err)
		return false
	}

	// Download plugins
	if len(c.Plugins) > 0 {

		// Remove all current files
		err = os.RemoveAll(config.GetPluginDirFilePath("manager"))
		if err != nil {
			config.Log.Error(err)
			return false
		}

		// Make sure the directory exists
		err = os.MkdirAll(config.GetPluginDirFilePath("manager"), 0755)
		if err != nil {
			config.Log.Error(err)
			return false
		}

		for name, url := range c.Plugins {

			if url == "" {
				continue
			}

			// Download actual file
			f := config.GetPluginDirFilePath("manager/" + name)
			err := downloadFile(f, url)
			if err != nil {
				config.Log.Error(err)
				continue
			}

		}
	}

	return true
}

func SanatizeFilename(f string) string {
	sf := strings.ToLower(f)
	sf = strings.ReplaceAll(sf, " ", "_")
	sf = sanRegex.ReplaceAllString(sf, "")
	return sf
}
