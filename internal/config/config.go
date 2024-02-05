package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type settings struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Token   string `yaml:"token"`
	TLS     struct {
		Cert string `yaml:"cert"`
		Key  string `yaml:"key"`
	}
	Units            string              `yaml:"defaultUnits"`
	PluginDir        string              `yaml:"pluginDir"`
	PluginTypes      map[string][]string `yaml:"pluginTypes"`
	ExcludeFsTypes   []string            `yaml:"excludeFsTypes"`
	RunPluginsAsRoot bool                `yaml:"runPluginsAsRoot"`
	Debug            bool                `yaml:"debug"`
	Senders          []SenderCfg         `yaml:"senders"`
	PassiveChecks    []CheckCfg          `yaml:"checks"`
	Manager          ManagerCfg          `yaml:"manager"`
}

type SenderCfg struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
	Type  string `yaml:"type"`
}

type CheckCfg struct {
	Hostname    string `yaml:"hostname"`
	Servicename string `yaml:"servicename"`
	Interval    string `yaml:"interval"`
	Endpoint    string `yaml:"endpoint"`
	Options     Values `yaml:"options"`
	Disabled    bool
	NextRun     time.Time
}

type ManagerCfg struct {
	Url        string `yaml:"url"`
	APIKey     string `yaml:"apikey"`
	IgnoreCert bool   `yaml:"ignoreCert"`
}

type Values struct {
	Check    bool
	Pretty   bool
	Plugin   string
	Name     string
	Path     string
	Args     []string
	Against  string
	Expected string
	Warning  string
	Critical string
	Delta    int
	units    string
}

type Data struct {
	Checks  []CheckCfg
	Senders []SenderCfg
	Secrets map[string]string
	sync.RWMutex
}

var Version string
var DebugMode bool
var ConfigFile string
var ConfigDefaultFile embed.FS
var ConfigDir string
var PluginDir string
var Settings = new(settings)
var AllowedUnits = []string{"B", "kB", "MB", "GB", "TB", "PB", "KiB", "MiB", "GiB", "TiB", "PiB"}

// Actual config data (checks, senders, secrets) after parse
var CfgData Data

func (v *Values) Units() string {
	units := Settings.Units
	if v.units != "" {
		units = v.units
	}
	if units == "" || !Contains(AllowedUnits, units) {
		units = "B"
	}
	return units
}

func (c *CheckCfg) Name() string {
	if c.Servicename != "" {
		return fmt.Sprintf("%s - %s", c.Hostname, c.Servicename)
	}
	return c.Hostname
}

func (c *CheckCfg) isEmpty() bool {
	return c.Hostname == ""
}

func (c *CheckCfg) apply(nc CheckCfg) {

}

// ================================
// Actual config settings
// ================================

func (s *settings) GetServerHost() string {
	return fmt.Sprintf("%s:%s", s.Address, strconv.Itoa(s.Port))
}

// ================================
// Other functions
// ================================

func ParseValues(r *http.Request) Values {

	pretty, _ := strconv.ParseBool(r.FormValue("pretty"))
	check, _ := strconv.ParseBool(r.FormValue("check"))
	delta, _ := strconv.Atoi(r.FormValue("delta"))

	var v = Values{
		Check:    check,
		Pretty:   pretty,
		Plugin:   r.FormValue("plugin"),
		Name:     r.FormValue("name"),
		Path:     r.FormValue("path"),
		Args:     r.Form["arg"],
		Against:  r.FormValue("against"),
		Expected: r.FormValue("expected"),
		Warning:  r.FormValue("warning"),
		Critical: r.FormValue("critical"),
		Delta:    delta,
	}

	// Override units (empty string is allowed to clear old value)
	units := r.FormValue("units")
	validUnit := Contains(AllowedUnits, units)
	if validUnit || units == "" {
		v.units = units
	}

	return v
}

func InitConfig(file string, defaultFile embed.FS) error {

	// If we don't pass a config, look for a config
	if file == "" {
		cfgFile, err := findConfig()
		if err != nil {
			return err
		}
		ConfigFile = cfgFile
	} else {
		ConfigFile = file
	}

	ConfigDefaultFile = defaultFile
	err := ParseConfig()

	return err
}

func ParseConfig() error {

	// Load the secrets file
	var err error
	CfgData.Secrets, err = ParseSecretsFile()
	if err != nil {
		Log.Error(err)
	}

	// Parse main config file
	err = ParseFile(ConfigDefaultFile)
	if err != nil {
		return err
	}

	// Parse checks/sender configs in directory
	ParseConfigDir()

	LogDebug("Configuration:")
	LogDebugf(" - Checks: %d", len(CfgData.Checks))
	LogDebugf(" - Senders: %d", len(CfgData.Senders))

	return nil
}

// Run through all YAML config files in the config directory if it exists
// and import any configurations for checks and senders. Does not return
// actual error but logs issues.
func ParseConfigDir() {

	p := GetConfigDirFilePath("")
	i, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			// We don't care if the file doesn't exist, only log in debug mode
			LogDebugf("Config dir does not exist: %s", p)
		} else {
			LogDebug(err)
		}
		return
	}

	if !i.IsDir() {
		Log.Errorf("Config dir should be a directory: %s", p)
		return
	}

	// Look for both .yml and .yaml
	exts := []string{"yml", "yaml"}
	files := []string{}
	for _, ext := range exts {
		fs, err := filepath.Glob(filepath.Join(p, "**/*."+ext))
		if err != nil {
			Log.Error(err)
			return
		}
		files = append(files, fs...)
	}

	LogDebugf("CONFIG: Config dir files found: %v", files)

	// Reset the checks and senders to the orginal config file values
	CfgData.Lock()
	defer CfgData.Unlock()
	CfgData.Checks = Settings.PassiveChecks
	CfgData.Senders = Settings.Senders

	// Load the secrets file
	CfgData.Secrets, err = ParseSecretsFile()
	if err != nil {
		Log.Error(err)
	}

	// Load up the senders and checks from each file
	for _, file := range files {

		yamlData, err := os.ReadFile(file)
		if err != nil {
			Log.Error(err)
			continue
		}

		replaceVariables(&yamlData)

		// Replace secrets in file before parsing
		if len(CfgData.Secrets) > 0 {
			for k, v := range CfgData.Secrets {
				// Build regex something like "\$VARIABLE(\s|\r|\n)" and append the group onto variable
				re, err := regexp.Compile(fmt.Sprintf("\\$%s(\\s|\\r|\\n)", k))
				if err == nil {
					yamlData = re.ReplaceAll(yamlData, []byte(v+"$1"))
				}
			}
		}

		var tmp settings
		err = yaml.Unmarshal(yamlData, &tmp)
		if err != nil {
			Log.Error(err)
			continue
		}

		// Senders specifically are not additive, only one of the same URL can exist
		if len(tmp.Senders) > 0 {
			if len(CfgData.Senders) > 0 {
				for i, oSender := range CfgData.Senders {
					for _, nSender := range tmp.Senders {
						if oSender.Url == nSender.Url {
							CfgData.Senders[i] = nSender
						}
					}
				}
			} else {
				CfgData.Senders = tmp.Senders
			}
		}

		// Checks are additive, we need to add a new check to existing checks
		if len(tmp.PassiveChecks) > 0 {
			if len(CfgData.Checks) > 0 {

				// Check if a new check with same host/service name exists... if it does we need to
				// do an additive addition to the current passive checks
				for _, nCheck := range tmp.PassiveChecks {

					if nCheck.isEmpty() {
						continue
					}

					i := findPassiveCheck(nCheck)
					if i != -1 {
						fmt.Printf("found check: %v\n", CfgData.Checks[i])
						CfgData.Checks[i].apply(nCheck)
					} else {
						CfgData.Checks = append(CfgData.Checks, nCheck)
					}
				}

			} else {
				CfgData.Checks = tmp.PassiveChecks
			}
		}

	}

	// Make sure any checks are set to their options for check=1
	for i := range CfgData.Checks {
		CfgData.Checks[i].Options.Check = true
	}

}

func findPassiveCheck(check CheckCfg) int {
	for i, c := range CfgData.Checks {
		if check.Hostname == c.Hostname && check.Servicename == c.Servicename {
			return i
		}
	}
	return -1
}

func ParseSecretsFile() (map[string]string, error) {
	secrets := make(map[string]string)

	f := GetConfigDirFilePath("manager/secrets.json")
	if !FileExists(f) {
		return secrets, nil
	}

	data, err := os.ReadFile(f)
	if err != nil {
		return secrets, err
	}

	err = json.Unmarshal(data, &secrets)
	if err != nil {
		return secrets, err
	}

	return secrets, nil
}

func ParseFile(defaultFile embed.FS) error {

	var yamlData []byte
	var err error

	// Set defaults with the default config file so that we don't
	// have to worry about new config files having new options that
	// the user doesn't have in theirs in the future
	yamlData, err = defaultFile.ReadFile("build/package/config.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, Settings)
	if err != nil {
		return err
	}

	yamlData, err = os.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	replaceVariables(&yamlData)

	err = yaml.Unmarshal(yamlData, Settings)
	if err != nil {
		return err
	}

	// Set defaults if we need to
	if Settings.PluginDir == "" {
		Settings.PluginDir = getPluginDir()
	}

	// Check if we should set debug mode on, it may already be forced
	// on during initialization in main.go
	if !DebugMode {
		if Settings.Debug {
			DebugMode = true
		}
	}

	return nil
}

// Set the version to the version in the VERSION file if it doesn't
// already exist during build time (this is mostly for DEV)
func ParseVersion(versionFile embed.FS) {
	v, _ := versionFile.ReadFile("VERSION")
	if Version == "" {
		Version = strings.TrimSpace(string(v))
	}
}

func GetConfigFilePath(name string) string {
	return filepath.Join(ConfigDir, name)
}

func GetConfigDirFilePath(path string) string {
	return GetConfigFilePath(filepath.Join("conf.d", path))
}

func GetPluginDirFilePath(name string) string {
	return filepath.Join(Settings.PluginDir, name)
}

// Function to look for the config file in the normal locations
// that it could be on standard systems
func findConfig() (string, error) {
	var err error
	if ConfigDir != "" {
		return GetConfigFilePath("config.yml"), nil
	}
	paths := []string{
		"config.yml",
		"build/package/config.yml",
		"/etc/rcagent/config.yml",
		"/usr/local/rcagent/config.yml",
		"C:\\Program Files\\rcagent\\config.yml",
		"C:\\Program Files (x86)\\rcagent\\config.yml",
	}
	for _, p := range paths {
		_, err = os.Stat(p)
		if os.IsNotExist(err) {
			continue
		}
		return p, nil
	}
	return "", err
}

// Gets plugin directory (global PluginDir is set during compilation and packaging)
func getPluginDir() string {
	if PluginDir != "" {
		return PluginDir
	}
	// Try to set to a default value if we didn't get one
	paths := []string{
		"plugins/",
		"/usr/lib64/rcagent/plugins",
		"/usr/lib/rcagent/plugins",
		"C:\\Program Files\\rcagent\\plugins",
	}
	for _, p := range paths {
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			continue
		}
		return p
	}
	return ""
}

// Replaces variables in a file with values using regex to validate line ends/spaces after variables.
func replaceVariables(data *[]byte) {

	hostname, _ := os.Hostname()

	// Replace $LOCAL_HOSTNAME with local agent hostname
	if hostname != "" {
		re := regexp.MustCompile(`\$LOCAL_HOSTNAME(\s|\z)`)
		*data = re.ReplaceAll(*data, []byte(hostname+"$1"))
	}

	// Replace all secrets if found
	if len(CfgData.Secrets) > 0 {
		for k, v := range CfgData.Secrets {
			// Build regex something like "\$VARIABLE(\s|\z)" and append the group onto variable
			re, err := regexp.Compile(fmt.Sprintf("\\$%s(\\s|\\z)", k))
			if err == nil {
				*data = re.ReplaceAll(*data, []byte(v+"$1"))
			}
		}
	}
}
