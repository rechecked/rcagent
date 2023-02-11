
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
    "github.com/rechecked/rcagent/internal/api"
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
        config.Settings.TLS.Cert = "rcagent.pem"
        config.Settings.TLS.Key = "rcagent.key"
        if !config.FileExists(config.Settings.TLS.Cert) {
            err := GenerateCert()
            if err != nil {
                log.Error(err)
            }
        }
    }

    // Set up endpoints
    EndpointFunc("memory/virtual", api.HandleMemory)
    EndpointFunc("memory/swap", api.HandleSwap)
    EndpointFunc("cpu/percent", api.HandleCPU)
    EndpointFunc("services", api.HandleServices)
    EndpointFunc("processes", api.HandleProcesses)
    EndpointFunc("disks", api.HandleDisks)
    EndpointFunc("disks/inodes", api.HandleInodes)
    EndpointFunc("plugins", api.HandlePlugins)
    EndpointFunc("system", api.HandleSystem)

    // Unix only
    // TODO: add linux logs
    if runtime.GOOS != "windows" {
        EndpointFunc("load", api.HandleLoad)
    }

    // Windows only
    // TODO: add windows logs and counters

    // Add handlers and run server with config
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/api/", handleAPI)

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

func EndpointFunc(path string, endpoint Endpoint) {
    endpoints[path] = endpoint
}

func handleMain(w http.ResponseWriter, r *http.Request) {
    errorHandler(w, r, http.StatusNotFound)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {

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

    // Get API endpoint and path from url
    fullpath := strings.TrimPrefix(r.URL.Path, "/api/")

    endpoint := endpoints[fullpath]
    if endpoint != nil {

        // Get the data back from the endpoint
        e := endpoint(values)

        // Check if we are checkable type
        chk, ok := e.(api.Checkable)
        if values.Check && ok {
            check := api.GetCheckResult(chk, values.Warning, values.Critical)
            jsonData, err = ConvertToJson(check, values.Pretty)
        }

        // Check if we are a checkable against type
        chk2, ok2 := e.(api.CheckableAgainst)
        if values.Check && ok2 {
            check := api.GetCheckAgainstResult(chk2, values.Expected)
            jsonData, err = ConvertToJson(check, values.Pretty)
        }

        // If we aren't doing a check, convert endpoint return to JSON
        if len(jsonData) == 0 && err == nil {
            jsonData, err = ConvertToJson(e, values.Pretty)
        }

    } else {

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
        error := serverError{
            Message: "Invalid API endpoint path given",
            Status: "error",
            Endpoints: e,
        }
        jsonData, err = ConvertToJson(error, values.Pretty)
    }

    /*
    data := struct{
        Thresholds *api.Thresholds `json:"thresholds"`
    }{
        Thresholds: thresholds,
    }
    */

    if err != nil {
        log.Errorf("Error getting data. Err: %s", err)
    }
    w.Write(jsonData)
}

func ConvertToJson(i interface{}, pretty bool) ([]byte, error) {
    if pretty {
        return json.MarshalIndent(i, "", "    ")
    }
    return json.Marshal(i)
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
