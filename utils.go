package main

import (
	"errors"
	"path/filepath"

	"github.com/revengel/enpass2gopass/enpass"
	log "github.com/sirupsen/logrus"
)

func setLogLevel(level string, debug bool) (err error) {
	if debug {
		log.SetLevel(log.DebugLevel)
		return
	}

	lvl, err := log.ParseLevel(level)
	if err != nil {
		return
	}

	log.SetLevel(lvl)
	return
}

func getGopassPath(folder string, item enpass.DataItem) (out string, err error) {
	switch {
	case item.IsTrashed():
		out = filepath.Join(out, "trash")
	case item.IsArchived():
		out = filepath.Join(out, "archive")
	case item.IsFavorite():
		out = filepath.Join(out, "favorite")
	}

	var cat = item.GetCategoryPath()
	if cat == "" {
		return "", errors.New("category cannot be empty")
	}
	out = filepath.Join(out, cat)

	if folder != "" {
		out = filepath.Join(out, folder)
	}

	var title = item.GetTitlePath()
	if title == "" {
		return "", errors.New("title cannot be empty")
	}
	out = filepath.Join(out, title)

	return out, nil
}
