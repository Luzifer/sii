package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/Luzifer/rconfig/v2"
	"github.com/Luzifer/sii"
	log "github.com/sirupsen/logrus"
)

var (
	cfg = struct {
		DecryptKey     string `flag:"decrypt-key" default:"" description:"Hex formated decryption key" validate:"nonzero"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	decryptKey []byte
	version    = "dev"
)

func init() {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("sii-decrypt %s\n", version)
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
	decryptKey, err = hex.DecodeString(cfg.DecryptKey)
	if err != nil {
		log.WithError(err).Fatal("Unable to read encryption key")
	}

	if len(rconfig.Args()) != 2 {
		log.Fatal("Expecting exactly one SII file as an argument")
	}

	stat, err := os.Stat(rconfig.Args()[1])
	if err != nil {
		log.WithError(err).Fatal("Unable to read info of save-file")
	}

	f, err := os.Open(rconfig.Args()[1])
	if err != nil {
		log.WithError(err).Fatal("Unable to open encrypted file")
	}
	defer f.Close()

	sii.SetEncryptionKey(decryptKey)

	contentReader, err := sii.DecryptRaw(f, stat.Size())
	if err != nil {
		log.WithError(err).Fatal("Unable to decrypt file")
	}

	tf, err := os.Create("/tmp/output")
	if err != nil {
		log.WithError(err).Fatal("Unable to create output file")
	}
	defer tf.Close()

	if _, err := io.Copy(tf, contentReader); err != nil {
		log.WithError(err).Fatal("Unable to copy content")
	}
}
