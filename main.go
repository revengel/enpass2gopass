// main package
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/revengel/enpass2gopass/store"
	"github.com/revengel/enpass2gopass/store/enpass"
	"github.com/revengel/enpass2gopass/store/gopass"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "0201-150405"
	customFormatter.FullTimestamp = false

	logger.SetFormatter(customFormatter)
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.WarnLevel)
}

func main() {
	var (
		prefix        string
		logLevel      string
		dryrun, debug bool
		gp            store.StoreDestination
		err           error
	)

	flag.StringVar(&prefix, "prefix", "enpass", "gopass path prefix")
	flag.StringVar(&logLevel, "log-level", logrus.InfoLevel.String(), "log level")
	flag.BoolVar(&dryrun, "dry-run", false, "do not write changes to gopass")
	flag.BoolVar(&debug, "debug", false, "enable debug log level")
	flag.Parse()

	err = setLogLevel(logLevel, debug)
	if err != nil {
		logger.Fatalf("Cannot set log level format '%s': %s", logLevel, err.Error())
	}

	values := flag.Args()
	if len(values) == 0 {
		logger.Fatal("Need to set path to json file with data")
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

	gp, err = gopass.NewStore(ctx, prefix, dryrun, logger)
	if err != nil {
		logger.Fatalf("Failed to connect gopass: %s", err)
	}

	defer gp.Close()

	data, err := enpass.LoadData(values[0])
	if err != nil {
		logger.Fatalf("Cannot load data from json file: %s", err.Error())
	}

	for _, item := range data.Items {
		var (
			err error
			l   = logger.WithField("type", "item")
		)

		gopassPath, err := gopass.GetSecretPath(item)
		if err != nil {
			logger.WithError(err).Fatal("cannot get gopass path")
		}

		var ll = l.WithField("gopassKey", gopassPath)
		fields, err := item.GetFields()
		if err != nil {
			ll.WithError(err).Fatal("cannot get item fields")
		}

		_, err = gp.Save(fields, gopassPath)
		if err != nil {
			ll.WithError(err).Fatal("cannot save secret")
		}
	}

	_, err = gp.Cleanup()
	if err != nil {
		logger.WithError(err).Fatal("cannot cleanup passwords storage")
	}
}
