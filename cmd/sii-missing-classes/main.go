package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/go_helpers/v2/str"
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
		fmt.Printf("sii-missing-blocks %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	if len(rconfig.Args()) != 2 {
		log.Fatal("Missing filename to analyze")
	}

	filename := rconfig.Args()[1]

	log.WithField("filename", filename).Info("Loading file")
	unit, err := sii.ReadUnitFile(filename)
	if err != nil {
		log.Fatalf("Unable to parse input file: %s", err)
	}

	var unknownClasses []string

	for _, e := range unit.Entries {
		if reflect.TypeOf(e).Elem() == reflect.TypeOf(sii.RawBlock{}) {
			log.WithFields(log.Fields{
				"class": e.Class(),
				"name":  e.Name(),
			}).Debug("Found raw block")

			if !str.StringInSlice(e.Class(), unknownClasses) {
				unknownClasses = append(unknownClasses, e.Class())
			}
		}
	}

	if len(unknownClasses) == 0 {
		log.Info("No unknown classes found")
		return
	}

	sort.Strings(unknownClasses)

	log.WithField("count", len(unknownClasses)).Info("Found unknown classes")

	fmt.Println(strings.Join(unknownClasses, "\n"))
}
