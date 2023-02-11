
package api

import (
    "io/ioutil"
    "strings"
    "os"
    "os/user"
    "os/exec"
    "runtime"
    "errors"
    "path/filepath"
    "fmt"
    "github.com/go-cmd/cmd"
    "github.com/rechecked/rcagent/internal/config"
)

type Plugin struct {
    cmd  []string
    args []string
    name string
    path string
}

type PluginResults struct {
    Output   string `json:"output"`
    ExitCode int    `json:"exitcode"`
}

func (p *Plugin) CreateCmd() error {

    var n, a bool = false, false

    file := filepath.Clean(p.path + "/" + p.name)
    ext := filepath.Ext(p.name)

    // Check plugin exists
    if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
        return err
    }

    // Check if we have an extension in the config
    for e, tmpl := range config.Settings.PluginTypes {
        if e == ext {
            // Replace the template with values
            for _, f := range tmpl {
                if strings.Contains(f, "$pluginName") {
                    p.cmd = append(p.cmd, file)
                    n = true
                    continue
                }
                if strings.Contains(f, "$pluginArgs") {
                    p.cmd = append(p.cmd, p.args...)
                    a = true
                    continue
                }
                p.cmd = append(p.cmd, f)
            }
            break
        }
    }

    // Append the file and args if we haven't yet
    if (!n) {
        p.cmd = append(p.cmd, file)
    }
    if (!a) {
        p.cmd = append(p.cmd, p.args...)
    }

    if config.Settings.Debug {
        fmt.Printf("Command: %s\n", p.cmd)
    }

    return nil
}

func (p *Plugin) Run() PluginResults {

    var out string
    c := new(cmd.Cmd)

    // Validate user exists and can run before we start
    if !isValidUser() {
        return PluginResults{
            Output: "The user 'rcagent' does not seem to exist on the system. To run plugins as root set runAsRoot: true in the config.",
            ExitCode: 3,
        }
    }
    
    // Create options
    options := cmd.Options{
        Buffered: true,
        BeforeExec: []func(*exec.Cmd){setUser},
        LineBufferSize: cmd.DEFAULT_LINE_BUFFER_SIZE,
    }

    if len(p.cmd) >= 2 {
        c = cmd.NewCmdOptions(options, p.cmd[0], p.cmd[1:]...)
    } else if len(p.cmd) == 1 {
        c = cmd.NewCmdOptions(options, p.cmd[0])
    } else {
        return PluginResults{
            Output: fmt.Sprintf("Error running plugin, command: %s.", p.cmd),
            ExitCode: 3,
        }
    }

    s := <-c.Start()
    if s.Error != nil {
        out = fmt.Sprintf("%s", s.Error)
        s.Exit = 1
    } else {
       if s.Exit == 0 {
            out = strings.Join(s.Stdout, "\n")
        } else {
            out = strings.Join(s.Stderr, "\n")
        }
    }

    return PluginResults{
        Output: out,
        ExitCode: s.Exit,
    }
}

func HandlePlugins(cv config.Values) interface{} {

    var res interface{}
    data, err := getPlugins()
    if err != nil {
        return nil
    }

    if cv.Plugin != "" {
        plugin := Plugin{
            name: cv.Plugin,
            path: config.Settings.PluginDir,
            args: cv.Args,
        }
        err = plugin.CreateCmd()
        if err == nil {
            res = plugin.Run()
        } else {
            res = PluginResults{
                Output: "Plugin does not exist",
                ExitCode: 1,
            }
        }
    } else {
        res = struct{
            Plugins []string `json:"plugins"`
        }{
            Plugins: data,
        }
    }

    return res
}

func getPlugins() ([]string, error) {

    var plugins []string

    files, err := ioutil.ReadDir(config.Settings.PluginDir)
    if err != nil {
        return plugins, err
    }

    for _, file := range files {
        if !file.IsDir() {
            plugins = append(plugins, file.Name())
        }
    }

    return plugins, nil
}

func isValidUser() bool {
    if runtime.GOOS != "windows" {
        u, err := user.Lookup("rcagent")
        if err != nil {
            return false
        }
        if u.Uid == "" {
            return false
        }
    }
    return true
}
