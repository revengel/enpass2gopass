package main

import (
	"fmt"

	"github.com/revengel/enpass2gopass/enpass"
	"github.com/revengel/enpass2gopass/gopassstore"
	log "github.com/sirupsen/logrus"
)

func getGopassItemSecrets(prefix string, item enpass.DataItem) (map[string]*gopassstore.GopassSecret, error) {
	var (
		err    error
		o      = make(map[string]*gopassstore.GopassSecret)
		s      = gopassstore.NewEmptyGopassSecret()
		folder = item.GetFirstFolder(foldersMap)
	)

	gopassPath, err := getGopassPath(prefix, folder, item)
	if err != nil {
		return o, err
	}

	gopassDataPath, err := getGopassDataPath(gopassPath)
	if err != nil {
		return o, err
	}

	err = s.Set("title", item.GetTitle(), false)
	if err != nil {
		return o, err
	}

	err = s.Set("subtitle", item.GetSubtitle(), false)
	if err != nil {
		return o, err
	}

	err = s.Set("category", item.GetCategoryPath(), false)
	if err != nil {
		return o, err
	}

	s.AddYamlData("note", item.GetNote())

	err = s.Set("tags", item.GetFoldersStr(foldersMap), false)
	if err != nil {
		return o, err
	}

	for _, field := range item.Fields {
		var ignoreTypes = []string{"section", ".Android#"}
		if field.IsDeleted() || field.CheckTypes(ignoreTypes) {
			continue
		}

		if field.GetLabel() == "" || field.GetValue() == "" {
			continue
		}

		if field.CheckType("password") {
			passwd := field.GetValue()
			if ok := s.SetPassword(passwd); ok {
				continue
			}
		}

		var labelName = field.GetLabel()
		if field.CheckType("totp") {
			labelName = "totp"
		}

		if labelName == "e_mail" {
			labelName = "email"
		}

		if field.IsMultiline() {
			s.AddYamlData(labelName, field.GetValue())
			continue
		}

		err = s.Set(labelName, field.GetValue(), field.IsMultiline())
		if err != nil {
			return o, err
		}
	}

	for _, attach := range item.Attachments {
		var name = attach.GetLabelName()
		isText, err := attach.IsTextData()
		if err != nil {
			return o, err
		}

		if isText {
			labelName := fmt.Sprintf("attachment - %s", attach.GetName())
			val, err := attach.GetDataString()
			if err != nil {
				return o, err
			}
			s.AddYamlData(labelName, val)
			continue
		}

		gopassAttachPath, err := getGopassAttachPath(gopassPath, name)
		if err != nil {
			return o, err
		}

		var dataB64 = attach.GetDataBase64Encoded()
		var flName = attach.GetNameOriginal()
		attachSec, err := gopassstore.NewAttachmentGopassSecret(flName, dataB64)
		if err != nil {
			return o, err
		}

		o[gopassAttachPath] = attachSec
	}

	o[gopassDataPath] = s
	return o, err
}

func gopassSaveSecret(s *gopassstore.GopassSecret, gp *gopassstore.Gopass, p string, dryrun bool, l *log.Entry) error {
	if dryrun {
		return nil
	}

	secret, err := s.GetSecret()
	if err != nil {
		return err
	}

	saved, err := gp.SetIfChanged(secret, p)
	if err != nil {
		return err
	}

	if saved {
		l.Info("secret has been updated")
	} else {
		l.Debug("secret already in actual state")
	}
	return nil
}
