package gopass

import (
	"errors"
	"path/filepath"
)

type SecretPathItemInterface interface {
	IsTrashed() bool
	IsArchived() bool
	IsFavorite() bool

	GetCategoryPath() string
	GetTitlePath() string
	GetFirstFolder() string
}

func GetSecretPath(item SecretPathItemInterface) (out string, err error) {
	var folder = item.GetFirstFolder()
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
