package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/sii"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func init() {
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("sii-test-read-write %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	if len(rconfig.Args()) != 3 {
		log.Fatal("Usage: sii-test-read-write <input-file> <output-file>")
	}

	log.Info("Loading input file")
	unit, err := sii.ReadUnitFile(rconfig.Args()[1])
	if err != nil {
		log.Fatalf("Unit-file parsing failed: %s", err)
	}

	log.Info("Writing output file")
	if err = sii.WriteUnitFile(rconfig.Args()[2], unit); err != nil {
		log.Fatalf("Unit-file writing failed: %s", err)
	}

	log.Info("Done")
}
