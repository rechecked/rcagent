
package config

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "os"
    "embed"
    "gopkg.in/yaml.v3"
)

type settings struct {
    Address string `yaml:"address"`
    Port int `yaml:"port"`
    Token string `yaml:"token"`
    TLS struct {
        Cert string `yaml:"cert"`
        Key string `yaml:"key"`
    }
    Units string `yaml:"defaultUnits"`
    PluginDir string `yaml:"pluginDir"`
    PluginTypes map[string][]string `yaml:"pluginTypes"`
    ExcludeFsTypes []string `yaml:"excludeFsTypes"`
    RunPluginsAsRoot bool `yaml:"runPluginsAsRoot"`
    Debug bool `yaml:"debug"`
}

type Values struct {
    Check    bool
    Pretty   bool
    Plugin   string
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
    if (v.units != "") {
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
        Check: check,
        Pretty: pretty,
        Plugin: r.FormValue("plugin"),
        Path: r.FormValue("path"),
        Args: r.Form["arg"],
        Against: r.FormValue("against"),
        Expected: r.FormValue("expected"),
        Warning: r.FormValue("warning"),
        Critical: r.FormValue("critical"),
        Delta: delta,
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

    return nil
}

// Function to look for the config file in the normal locations
// that it could be on standard systems
func findConfig() (string, error) {
    var err error
    if ConfigDir != "" {
        return ConfigDir+"/config.yml", nil
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
