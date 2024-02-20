package main

import (
	"github.com/sirupsen/logrus"
)

func setLogLevel(level string, debug bool) (err error) {
	if debug {
		logger.SetLevel(logrus.DebugLevel)
		return
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return
	}

	logger.SetLevel(lvl)
	return
}
