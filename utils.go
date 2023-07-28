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

func getUniquePath(in string) (out string) {
	var changed bool
	changed, out = insertedPaths.GetUniquePath(in)
	if changed {
		log.Warnf("gopass path '%s' will be rename to '%s'", in, out)
	}
	return out
}

func getGopassPath(prefix, folder string, item enpass.DataItem) (out string, err error) {
	out = prefix
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
	uniqPath := getUniquePath(out)

	insertedPaths.Register(out)
	return uniqPath, nil
}

func getGopassDataPath(prefix string) (out string, err error) {
	out = filepath.Join(prefix, "data")
	insertedPaths.Register(out)
	return out, err
}

func getGopassAttachPath(prefix, attachName string) (out string, err error) {
	if prefix == "" {
		return "", errors.New("prefix cannot be empty")
	}

	if attachName == "" {
		return "", errors.New("attachName cannot be empty")
	}

	out = filepath.Join(prefix, "attachments", attachName)
	out = getUniquePath(out)
	insertedPaths.Register(out)
	return out, err
}
