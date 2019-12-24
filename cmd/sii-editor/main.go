package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
	"github.com/Luzifer/sii"
)

var (
	cfg = struct {
		Config         string `flag:"config,c" vardefault:"config" description:"Optional configuration file"`
		DecryptKey     string `flag:"decrypt-key" default:"" description:"Hex formated decryption key" validate:"nonzero"`
		Listen         string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	baseGameUnit *sii.Unit
	router       = mux.NewRouter()
	userConfig   *configFile

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
	decryptKey, err := hex.DecodeString(cfg.DecryptKey)
	if err != nil {
		log.WithError(err).Fatal("Unable to read encryption key")
	}

	sii.SetEncryptionKey(decryptKey)

	if userConfig, err = loadUserConfig(cfg.Config); err != nil && err != errUserConfigNotFound {
		log.WithError(err).Fatal("Unable to load user config")
	}

	if err = userConfig.loadDefaults(); err != nil {
		log.WithError(err).Fatal("Unable to load missing defaults for user config")
	}

	log.Info("Loading game base data...")

	if baseGameUnit, err = readBaseData(); err != nil {
		log.WithError(err).Fatal("Unable to load game definitions")
	}

	log.WithFields(log.Fields{
		"companies": len(baseGameUnit.BlocksByClass("company_permanent")),
		"cargos":    len(baseGameUnit.BlocksByClass("cargo_data")),
	}).Info("Game base data loaded")

	log.WithField("addr", cfg.Listen).Info("Starting API server...")

	if err := http.ListenAndServe(cfg.Listen, router); err != nil {
		log.WithError(err).Fatal("HTTP server caused an error")
	}
}
