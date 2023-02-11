//go:generate goversioninfo

package main

import (
    "flag"
    "log"
    "os"
    "embed"
    "github.com/kardianos/service"
    "github.com/rechecked/rcagent/internal/server"
    "github.com/rechecked/rcagent/internal/config"
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

    // Set up server configuration and run
    go server.Run(logger)

    // If we have a sender (passive checks)
    //go sender.NDRP()

    // Do work here
    return nil
}

func (p *program) Stop(s service.Service) error {
    // Stop should not block. Return with a few seconds.
    return nil
}

var logger service.Logger

//go:embed build/package/config.yml
var defaultConfigFile embed.FS

func main() {

    // All actions the service can perform
    action := flag.String("a", "run", "Service action to run: 'install', 'uninstall', or 'run'. Default is 'run'.")
    configFile := flag.String("f", "", "Config file location")
    version := flag.Bool("v", false, "Show version of rcagent")
    flag.Parse()

    // Show version and quit
    if *version {
        log.Printf("ReChecked Agent, version %s\n", config.Version)
        os.Exit(0)
    }

    svcConfig := &service.Config{
        Name:        "rcagent",
        DisplayName: "RCAgent",
        Description: "ReChecked system status and monitoring agent",
        Dependencies: []string{
            "Requires=network.target",
            "After=network-online.target syslog.target",
        },
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
