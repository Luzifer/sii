package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		DryRun         bool   `flag:"dry-run,n" default:"true" description:"Do NOT apply any destructive action"`
		GameDir        string `flag:"game-dir" default:"" description:"Directory ETS/ATS is installed in" validate:"nonzero"`
		SavePath       string `flag:"save-path" default:"" description:"Path to game.sii file to cleanup" validate:"nonzero"`
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
		fmt.Printf("save-clean  %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	log.Info("Loading base gamedata...")

	baseGameUnit, err := readBaseData()
	if err != nil {
		log.WithError(err).Fatal("Unable to load base gamedata")
	}

	log.WithFields(log.Fields{
		"cargos":    len(baseGameUnit.BlocksByClass("cargo_data")),
		"cities":    len(baseGameUnit.BlocksByClass("city_data")),
		"companies": len(baseGameUnit.BlocksByClass("company_permanent")),
	}).Info("Game base data loaded")

	log.Info("Loading savegame...")

	// Read savegame for cleaning
	game, err := loadSave(cfg.SavePath)
	if err != nil {
		log.WithError(err).Fatal("Unable to load savegame")
	}

	log.WithField("blocks", len(game.Entries)).Info("Savegame loaded")

	log.Info("Cleaning savegame...")

	// Execute cleaner
	if err = cleanSaveGame(baseGameUnit, game); err != nil {
		log.WithError(err).Fatal("Unable to clean savegame")
	}

	log.WithField("blocks", len(game.Entries)).Info("Savegame cleaned")

	if cfg.DryRun {
		log.Warn("Dry-Run enabled, skipping save!")
		return
	}

	// Store back savegame after cleaning
	if err = storeSave(cfg.SavePath, game); err != nil {
		log.WithError(err).Fatal("Unable to store savegame")
	}
}
