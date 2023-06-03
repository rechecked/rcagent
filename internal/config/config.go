package config

import (
	"embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

var Version string
var ConfigDir string
var PluginDir string
var Settings = new(settings)
var AllowedUnits = []string{"B", "kB", "MB", "GB", "TB", "PB", "KiB", "MiB", "GiB", "TiB", "PiB"}

func (s *settings) GetServerHost() string {
	return fmt.Sprintf("%s:%s", s.Address, strconv.Itoa(s.Port))
}

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

func ParseFile(file string, defaultFile embed.FS) error {

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

	// If we don't pass a config, look for a config
	if file == "" {
		file, err = findConfig()
		if err != nil {
			return err
		}
	}

	yamlData, err = ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, Settings)
	if err != nil {
		return err
	}

	// Set defaults if we need to
	if Settings.PluginDir == "" {
		Settings.PluginDir = getPluginDir()
	}

	// Make sure any checks are set to their options for check=1
	hostname, _ := os.Hostname()
	for i, c := range Settings.PassiveChecks {
		Settings.PassiveChecks[i].Options.Check = true
		if strings.Contains(c.Hostname, "$HOST") {
			Settings.PassiveChecks[i].Hostname = strings.Replace(c.Hostname, "$HOST", hostname, -1)
		}
	}

	return nil
}

// Set the version to the version in the VERSION file if it doesn't
// already exist during build time (this is mostly for DEV)
func ParseVersion(versionFile embed.FS) {
	v, _ := versionFile.ReadFile("VERSION")
	if Version == "" {
		Version = string(v)
	}
}

// Function to look for the config file in the normal locations
// that it could be on standard systems
func findConfig() (string, error) {
	var err error
	if ConfigDir != "" {
		return ConfigDir + "/config.yml", nil
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
