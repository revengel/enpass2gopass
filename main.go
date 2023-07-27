package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/revengel/enpass2gopass/enpass"
	"github.com/revengel/enpass2gopass/gopassstore"
	log "github.com/sirupsen/logrus"
)

var (
	insertedPaths = newInsertedPaths()
	foldersMap    enpass.FoldersMap
)

func init() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "0201-150405"
	customFormatter.FullTimestamp = false

	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	var (
		prefix        string
		logLevel      string
		dryrun, debug bool
		gp            *gopassstore.Gopass
		err           error
	)

	flag.StringVar(&prefix, "prefix", "enpass", "gopass path prefix")
	flag.StringVar(&logLevel, "log-level", log.InfoLevel.String(), "log level")
	flag.BoolVar(&dryrun, "dry-run", false, "do not write changes to gopass")
	flag.BoolVar(&debug, "debug", false, "enable debug log level")
	flag.Parse()

	err = setLogLevel(logLevel, debug)
	if err != nil {
		log.Fatalf("Cannot set log level format '%s': %s", logLevel, err.Error())
	}

	values := flag.Args()
	if len(values) == 0 {
		log.Fatal("Need to set path to json file with data")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()
	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	gp, err = gopassstore.NewGopass(ctx)
	if err != nil {
		log.Fatalf("Failed to connect gopass: %s", err)
	}

	defer gp.Close()

	data, err := enpass.LoadData(values[0])
	if err != nil {
		log.Fatalf("Cannot load data from json file: %s", err.Error())
	}

	foldersMap = data.GetFoldersMap()

	for _, item := range data.Items {
		var (
			err error
			l   = log.WithField("type", "item")
		)

		itemSecrets, err := getGopassItemSecrets(prefix, item)
		if err != nil {
			l.WithError(err).Fatal("cannot generate gopass secrets")
		}

		for key, s := range itemSecrets {
			var ll = l.WithField("gopassKey", key)
			ll.Debug("saving secret to gopass store")

			if log.GetLevel() >= log.DebugLevel {
				// output data secrets
				fmt.Println()
				reader := bytes.NewReader(s.Bytes())
				io.Copy(os.Stdout, reader)
			}

			err = gopassSaveSecret(s, gp, key, dryrun, ll)
			if err != nil {
				ll.WithError(err).Fatal("cannot save secret")
			}
		}
	}

	var lc = log.WithField("type", "cleaner")
	ll, err := gp.List(`^` + prefix + `/`)
	if err != nil {
		lc.WithError(err).Fatal("cannot get gopass keys list")
	}

	for _, k := range ll {
		if c := insertedPaths.Check(k); c > 0 {
			continue
		}

		var lck = lc.WithField("gopassPath", k)
		lck.Info("gopass key will be deleted")
		err = gp.Remove(k)
		if err != nil {
			lck.Fatal("cannot delete gopass key")
		}
	}
}
