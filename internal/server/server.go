
package server

import (
    //"fmt"
    "net/http"
    "encoding/json"
    "os"
    "strings"
    "crypto/subtle"
    "errors"
    "sort"
    "runtime"
    "github.com/kardianos/service"
    "github.com/rechecked/rcagent/internal/config"
    "github.com/rechecked/rcagent/internal/status"
)

type Endpoint func(cv config.Values) interface{}

type serverError struct {
    Message   string   `json:"message"`
    Status    string   `json:"status"`
    Endpoints []string `json:"endpoints,omitempty"`
}

var log service.Logger
var endpoints = make(map[string]Endpoint)

func Run(l service.Logger) {
    log = l

    // Get details for server
    hostname, _ := os.Hostname()
    host := config.Settings.GetServerHost()
    log.Infof("Starting server: %s (%s)", hostname, host)

    // Check if we are using adhoc certs, if we are, generate them
    if config.Settings.TLS.Cert == "adhoc" && config.Settings.TLS.Key == "adhoc" {
        config.Settings.TLS.Cert = config.ConfigDir+"rcagent.pem"
        config.Settings.TLS.Key = config.ConfigDir+"rcagent.key"
        if !config.FileExists(config.Settings.TLS.Cert) {
            err := GenerateCert()
            if err != nil {
                log.Error(err)
            }
        }
    }

    SetupEndpoints();

    // Add handlers and run server with config
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/status/", handleStatusAPI)

    var err error
    if config.Settings.TLS.Cert != "" && config.Settings.TLS.Key != "" {
        err = http.ListenAndServeTLS(host, config.Settings.TLS.Cert, config.Settings.TLS.Key, nil)
    } else {
        err = http.ListenAndServe(host, nil)
    }
    if err != nil {
        log.Error(err)
    }
}

func SetupEndpoints() {
    // Set up saved values for network counters
    status.Setup()

    // Set up endpoints
    endpointFunc("memory/virtual", status.HandleMemory)
    endpointFunc("memory/swap", status.HandleSwap)
    endpointFunc("cpu/percent", status.HandleCPU)
    endpointFunc("disk", status.HandleDisks)
    //endpointFunc("docker", status.HandleDocker)
    endpointFunc("services", status.HandleServices)
    endpointFunc("processes", status.HandleProcesses)
    endpointFunc("plugins", status.HandlePlugins)
    endpointFunc("network", status.HandleNetworks)
    endpointFunc("system", status.HandleSystem)
    endpointFunc("system/users", status.HandleUsers)

    // Unix only
    if runtime.GOOS != "windows" {
        endpointFunc("load", status.HandleLoad)
        endpointFunc("disk/inodes", status.HandleInodes)
    }

    // Windows only
    // TODO: add counters
}

func GetDataFromEndpoint(path string, values config.Values) (interface{}, error) {
    endpoint := endpoints[path]
    if endpoint != nil {

        // Get the data back from the endpoint
        e := endpoint(values)

        // Check if we are checkable type
        chk, ok := e.(status.Checkable)
        if values.Check && ok {
            check := status.GetCheckResult(chk, values.Warning, values.Critical)
            return check, nil
        }

        // Check if we are a checkable against type
        chk2, ok2 := e.(status.CheckableAgainst)
        if values.Check && ok2 {
            check := status.GetCheckAgainstResult(chk2, values.Expected)
            return check, nil
        }

        // If we aren't doing a check, convert endpoint return to JSON
        return e, nil
    }
    return nil, errors.New("GetDataFromEndpoint: Endpoint does not exist")
}

func endpointFunc(path string, endpoint Endpoint) {
    endpoints[path] = endpoint
}

func handleMain(w http.ResponseWriter, r *http.Request) {
    errorHandler(w, r, http.StatusNotFound)
}

func handleStatusAPI(w http.ResponseWriter, r *http.Request) {

    var jsonData []byte
    var err error

    defer w.Header().Set("Content-Type", "application/json")

    if err = r.ParseForm(); err != nil {
        log.Error(err)
    }

    // Parse config values and build thresholds
    values := config.ParseValues(r)

    // Validate token
    if err = validateToken(w, r); err != nil {
        error := serverError{
            Message: "Could not authenticate: invalid token given",
            Status: "error",
        }
        jsonData, _ = ConvertToJson(error, values.Pretty)
        w.Write(jsonData)
        return
    }

    // Get status API endpoint and path from url
    fullpath := strings.TrimPrefix(r.URL.Path, "/status/")

    data, err := GetDataFromEndpoint(fullpath, values)
    if err != nil {
        // Find out if we have any endpoints we can give out a list
        // of accessible endpoints to the output
        var e []string
        prefix, _, _ := strings.Cut(fullpath, "/")
        for ep, _ := range endpoints {
            if strings.Contains(ep, prefix) {
                e = append(e, ep)
            }
        }
        sort.Strings(e)

        // Endpoint doesn't exist, give a generic invalid path error
        data = serverError{
            Message: "Invalid API endpoint path given",
            Status: "error",
            Endpoints: e,
        }
    }
    jsonData, err = ConvertToJson(data, values.Pretty)

    if err != nil {
        log.Errorf("Error getting data. Err: %s", err)
    }
    w.Write(jsonData)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
    defer w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if status == http.StatusNotFound {
        // Custom 404 style error
        error := serverError{
            Message: "Could not get data: invalid URL path",
            Status: "error",
        }
        jsonData, err := json.Marshal(error)
        if err != nil {
            log.Error(err)
        }
        w.Write(jsonData)
    }
}

func validateToken(w http.ResponseWriter, r *http.Request) error {
    token := []byte(r.FormValue("token"))
    configToken := []byte(config.Settings.Token)
    if subtle.ConstantTimeCompare(token, configToken) == 1 {
        return nil
    }
    return errors.New("validateToken: Could not authenticate token")
}
