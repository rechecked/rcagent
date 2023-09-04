//go:generate goversioninfo

package main

import (
	"embed"
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/kardianos/service"

	"github.com/rechecked/rcagent/internal/config"
	"github.com/rechecked/rcagent/internal/manager"
	"github.com/rechecked/rcagent/internal/sender"
	"github.com/rechecked/rcagent/internal/server"
)

type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {

	// Register with the manager on startup, since we may need a certificate
	go manager.Register()

	// Set up server configuration and run
	c := make(chan struct{})
	go runServer(c)

	// If we have a sender (passive checks)
	go sender.Run()

	// Connect to the manager for sync
	go manager.Run(c)

	return nil
}

func runServer(c chan struct{}) {
	restart := make(chan struct{})
	go server.Run(logger, restart)
	<- c
	restart <- struct{}{}
	<- restart
	go runServer(c)
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

var logger service.Logger

//go:embed build/package/config.yml
var defaultConfigFile embed.FS

//go:embed VERSION
var defaultVersion embed.FS

func main() {

	// All actions the service can perform
	action := flag.String("a", "run", "Service action to run: 'install', 'uninstall', or 'run'. Default is 'run'.")
	configFile := flag.String("f", "", "Config file location")
	version := flag.Bool("v", false, "Show version of rcagent")
	flag.Parse()

	// Parse/set version then show if someone does -v
	config.ParseVersion(defaultVersion)
	if *version {
		log.Printf("ReChecked Agent, version: %s\n", config.Version)
		os.Exit(0)
	}

	var deps []string
	name := "rcagent"
	if runtime.GOOS == "linux" {
		deps = []string{
			"Requires=network.target",
			"After=network-online.target syslog.target",
		}
	}
	// Change name on macos to conform to macos
	if runtime.GOOS == "darwin" {
		name = "io.rechecked.rcagent"
	}

	svcConfig := &service.Config{
		Name:         name,
		DisplayName:  "RCAgent",
		Description:  "ReChecked system status and monitoring agent",
		Dependencies: deps,
	}

	// Initialize config settings (no config.yml on install)
	if *action == "run" {
		err := config.ParseFile(*configFile, defaultConfigFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize service
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize service logger
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Run actions for the service (run, install, uninstall)
	switch *action {
	case "install":
		err = s.Install()
	case "uninstall":
		err = s.Uninstall()
	default:
		err = s.Run()
	}

	// Exit with error if we hit one
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

}
