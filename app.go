package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/revengel/enpass2gopass/store"
	"github.com/revengel/enpass2gopass/store/enpass"
	"github.com/revengel/enpass2gopass/store/gopass"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	EnpassJsonSourceType   = "enpassJsonSource"
	GopassDestinationType  = "gopassDestination"
	KeepassDestinationType = "keepassDestination"
)

type app struct {
	ctx         context.Context
	logger      *logrus.Logger
	source      store.StoreSource
	destination store.StoreDestination
}

func (a *app) Close() error {
	var err error
	if a.destination != nil {
		err = a.destination.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *app) SetLogLevel(cmd *cobra.Command, args []string) error {
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		a.logger.SetLevel(logrus.DebugLevel)
		return nil
	}

	logLevelStr, _ := cmd.Flags().GetString("log-level")
	if logLevelStr == "" {
		return nil
	}

	lvl, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		a.logger.Warnf("cannot parse log level '%s': %s", logLevelStr, err.Error())
		return nil
	}

	a.logger.SetLevel(lvl)
	return nil
}

func (a *app) Before(cmd *cobra.Command, args []string) error {
	var err error
	err = a.SetLogLevel(cmd, args)
	if err != nil {
		return err
	}

	prefix, _ := cmd.Flags().GetString("prefix")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	sourceProvider, _ := cmd.Flags().GetString("source-provider")
	destProvider, _ := cmd.Flags().GetString("destination-provider")

	switch sourceProvider {
	case EnpassJsonSourceType:
		enpassJsonPath, _ := cmd.Flags().GetString("source-enpass-json-path")
		if enpassJsonPath == "" {
			return errors.New("source enpass json file is not set")
		}
		a.source, err = enpass.NewEnpassJsonSource(enpassJsonPath)
	default:
		return fmt.Errorf("invalid source provider: %s", sourceProvider)
	}

	if err != nil {
		return fmt.Errorf("failed to connect source: %s", err)
	}

	switch destProvider {
	case GopassDestinationType:
		a.destination, err = gopass.NewStore(a.ctx, prefix, dryRun, a.logger)
	default:
		return fmt.Errorf("invalid destination provider: %s", destProvider)
	}

	if err != nil {
		return fmt.Errorf("failed to connect destination: %s", err)
	}

	return nil
}

func (a *app) Import(cmd *cobra.Command, args []string) error {
	items, err := a.source.LoadData()
	if err != nil {
		return fmt.Errorf("Cannot load data from source: %s", err.Error())
	}

	for _, item := range items {
		secretPath, err := item.GetSecretPath()
		if err != nil {
			return fmt.Errorf("cannot get seceret path: %s", err.Error())
		}

		fields, err := item.GetFields()
		if err != nil {
			return fmt.Errorf("cannot get item fields; secret key - '%s': %s", secretPath, err.Error())
		}

		_, err = a.destination.Save(fields, secretPath)
		if err != nil {
			return fmt.Errorf("cannot save secret; secret key - '%s': %s", secretPath, err.Error())
		}
	}

	_, err = a.destination.Cleanup()
	if err != nil {
		return fmt.Errorf("cannot cleanup passwords storage: %s", err.Error())
	}

	return nil
}
