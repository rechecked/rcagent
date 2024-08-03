package server

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/rechecked/rcagent/internal/config"
	"github.com/rechecked/rcagent/internal/endpoints"
	"github.com/rechecked/rcagent/internal/status"
)

type serverError struct {
	Message   string   `json:"message"`
	Status    string   `json:"status"`
	Endpoints []string `json:"endpoints,omitempty"`
}

func Setup() {
	config.LogDebug("Setting up endpoints")
	setupEndpoints()

	// Check if we are using adhoc certs, if we are, generate them
	if config.Settings.TLS.Cert == "adhoc" && config.Settings.TLS.Key == "adhoc" {
		config.Settings.TLS.Cert = config.GetConfigFilePath("rcagent.pem")
		config.Settings.TLS.Key = config.GetConfigFilePath("rcagent.key")
		if !config.FileExists(config.Settings.TLS.Cert) {
			config.LogDebug("Setting up certificates")
			err := GenerateCert(config.Settings.TLS.Cert, config.Settings.TLS.Key)
			if err != nil {
				config.Log.Error(err)
			}
		}
	}
}

func Run(restart chan struct{}) {

	// Get details for server
	hostname, _ := os.Hostname()
	host := config.Settings.GetServerHost()

	// Add handlers and run server with config
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleMain)
	mux.HandleFunc("/status/", handleStatusAPI)

	// Create server with config so we can restart it later
	srv := &http.Server{
		Addr:    host,
		Handler: mux,
	}

	config.Log.Infof("Starting server: %s (%s)", hostname, host)
	go serve(srv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
	case <-restart:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		defer func(r chan struct{}) {
			r <- struct{}{}
		}(restart)
		if err := srv.Shutdown(ctx); err != nil {
			config.Log.Errorf("Shutdown error: %v", err)
		}
	}
}

func serve(srv *http.Server) {
	var err error
	if config.Settings.TLS.Cert != "" && config.Settings.TLS.Key != "" {
		err = srv.ListenAndServeTLS(config.Settings.TLS.Cert, config.Settings.TLS.Key)
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil && err != http.ErrServerClosed {
		config.Log.Error(err)
	}
}

func setupEndpoints() {
	// Set up saved values for network counters
	status.Setup()

	// Set up endpoints
	endpointFunc("memory/virtual", status.HandleMemory)
	endpointFunc("memory/swap", status.HandleSwap)
	endpointFunc("cpu/percent", status.HandleCPU)
	endpointFunc("disk", status.HandleDisks)
	endpointFunc("services", status.HandleServices)
	endpointFunc("processes", status.HandleProcesses)
	endpointFunc("plugins", status.HandlePlugins)
	endpointFunc("network", status.HandleNetworks)
	endpointFunc("system", status.HandleSystem)
	endpointFunc("system/users", status.HandleUsers)
	endpointFunc("system/version", status.HandleVersion)

	// Unix only
	if runtime.GOOS != "windows" {
		endpointFunc("load", status.HandleLoad)
		endpointFunc("disk/inodes", status.HandleInodes)
		//endpointFunc("docker", status.HandleDocker)
	}

	// Windows only
	// TODO: add counters
}

func endpointFunc(path string, endpoint config.Endpoint) {
	config.Endpoints[path] = endpoint
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, http.StatusNotFound)
}

func handleStatusAPI(w http.ResponseWriter, r *http.Request) {
	setupHeader(&w)

	var jsonData []byte
	var err error

	if err = r.ParseForm(); err != nil {
		config.Log.Error(err)
	}

	// Parse config values and build thresholds
	values := config.ParseValues(r)

	// Validate token
	if err = validateToken(r); err != nil {
		error := serverError{
			Message: "Could not authenticate: invalid token given",
			Status:  "error",
		}
		jsonData, _ = ConvertToJson(error, values.Pretty)
		w.Write(jsonData)
		return
	}

	// Get status API endpoint and path from url
	fullpath := strings.TrimPrefix(r.URL.Path, "/status/")

	data, err := endpoints.GetDataFromEndpoint(fullpath, values)
	if err != nil {
		// Find out if we have any endpoints we can give out a list
		// of accessible endpoints to the output
		var e []string
		prefix, _, _ := strings.Cut(fullpath, "/")
		for ep := range config.Endpoints {
			if strings.Contains(ep, prefix) {
				e = append(e, ep)
			}
		}
		sort.Strings(e)

		// Endpoint doesn't exist, give a generic invalid path error
		data = serverError{
			Message:   "Invalid API endpoint path given",
			Status:    "error",
			Endpoints: e,
		}
	}
	jsonData, err = ConvertToJson(data, values.Pretty)

	if err != nil {
		config.Log.Errorf("Error getting data. Err: %s", err)
	}
	w.Write(jsonData)
}

func errorHandler(w http.ResponseWriter, status int) {
	setupHeader(&w)
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		// Custom 404 style error
		error := serverError{
			Message: "Could not get data: invalid URL path",
			Status:  "error",
		}
		jsonData, err := json.Marshal(error)
		if err != nil {
			config.Log.Error(err)
		}
		w.Write(jsonData)
	}
}

func validateToken(r *http.Request) error {
	token := []byte(r.FormValue("token"))
	configToken := []byte(config.Settings.Token)
	if subtle.ConstantTimeCompare(token, configToken) == 1 {
		return nil
	}
	return errors.New("validateToken: Could not authenticate token")
}

func setupHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
