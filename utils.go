package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/revengel/enpass2gopass/enpass"
	log "github.com/sirupsen/logrus"
)

func getHashFromBytes(in []byte) string {
	h := sha256.New()
	h.Write(in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func getHash(in string) string {
	return getHashFromBytes([]byte(in))
}

func truncStr(in string, maxLen int) string {
	if len(in) <= maxLen {
		return in
	}
	return strings.TrimSuffix(in[:maxLen], "_")
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

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
	uniqPath := insertedPaths.GetUniquePath(out)

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
	out = insertedPaths.GetUniquePath(out)
	insertedPaths.Register(out)
	return out, err
}
