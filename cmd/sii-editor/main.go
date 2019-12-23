package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Config         string `flag:"config,c" vardefault:"config" description:"Optional configuration file"`
		Listen         string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	router     = mux.NewRouter()
	userConfig *configFile

	version = "dev"
)

func init() {
	rconfig.SetVariableDefaults(map[string]string{
		"config": userConfigPath,
	})

	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("sii-editor %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	var err error

	if userConfig, err = loadUserConfig(cfg.Config); err != nil && err != errUserConfigNotFound {
		log.WithError(err).Fatal("Unable to load user config")
	}

	if err = userConfig.loadDefaults(); err != nil {
		log.WithError(err).Fatal("Unable to load missing defaults for user config")
	}

	if err := http.ListenAndServe(cfg.Listen, router); err != nil {
		log.WithError(err).Fatal("HTTP server caused an error")
	}
}
