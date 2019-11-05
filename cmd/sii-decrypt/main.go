package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/Luzifer/rconfig/v2"
	"github.com/Luzifer/sii"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	cfg = struct {
		Backup         bool   `flag:"backup,b" default:"true" description:"Bakcup encrypted file when in-place mode is used"`
		DecryptKey     string `flag:"decrypt-key" default:"" description:"Hex formated decryption key" validate:"nonzero"`
		InPlace        bool   `flag:"in-place,i" default:"false" description:"Write the decrypted version to the same file the encrypted version was read from"`
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

	var filename = rconfig.Args()[1]

	stat, err := os.Stat(filename)
	if err != nil {
		log.WithError(err).Fatal("Unable to read info of save-file")
	}

	f, err := os.Open(filename)
	if err != nil {
		log.WithError(err).Fatal("Unable to open encrypted file")
	}
	defer f.Close()

	sii.SetEncryptionKey(decryptKey)

	contentReader, err := sii.DecryptRaw(f, stat.Size())
	if err != nil {
		log.WithError(err).Fatal("Unable to decrypt file")
	}

	var output io.Writer = os.Stdout

	if cfg.InPlace {
		if cfg.Backup {
			if err := copyFile(filename, filename+".bak"); err != nil {
				log.WithError(err).Fatal("Unable to create file backup")
			}
		}

		// Close input file right now as we need to re-create it
		f.Close()

		tf, err := os.Create(filename)
		if err != nil {
			log.WithError(err).Fatal("Unable to create output file")
		}
		defer tf.Close()

		output = tf
	}

	if _, err := io.Copy(output, contentReader); err != nil {
		log.WithError(err).Fatal("Unable to copy content")
	}
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "Unable to open source file")
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, "Unable to open dest file")
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}
