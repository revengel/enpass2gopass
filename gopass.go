package main

import (
	"fmt"

	"github.com/revengel/enpass2gopass/enpass"
	"github.com/revengel/enpass2gopass/gopassstore"
	"github.com/revengel/enpass2gopass/store"
	log "github.com/sirupsen/logrus"
)

func getGopassItemSecrets(prefix string, item enpass.DataItem) (map[string]store.Secret, error) {
	var (
		err    error
		o      = make(map[string]store.Secret)
		s      store.Secret
		folder = item.GetFirstFolder(foldersMap)
	)

	s = gopassstore.NewEmptyGopassSecret()

	gopassPath, err := getGopassPath(prefix, folder, item)
	if err != nil {
		return o, err
	}

	gopassDataPath, err := getGopassDataPath(gopassPath)
	if err != nil {
		return o, err
	}

	err = s.Set("title", item.GetTitle(), store.SecretTitileField, false)
	if err != nil {
		return o, err
	}

	err = s.Set("subtitle", item.GetSubtitle(), store.SecretSimpleField, false)
	if err != nil {
		return o, err
	}

	err = s.Set("category", item.GetCategoryPath(), store.SecretSimpleField, false)
	if err != nil {
		return o, err
	}

	err = s.Set("note", item.GetNote(), store.SecretMultilineField, false)
	if err != nil {
		return o, err
	}

	err = s.Set("tags", item.GetFoldersStr(foldersMap), store.SecretSimpleField, false)
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

		var labelName = field.GetLabel()
		if field.CheckType("totp") {
			labelName = "totp"
		}

		if labelName == "e_mail" {
			labelName = "email"
		}

		var fieldType = store.SecretSimpleField
		if field.CheckType("password") {
			fieldType = store.SecretPasswordField
		}

		if field.IsMultiline() {
			fieldType = store.SecretMultilineField
		}

		err = s.Set(labelName, field.GetValue(), fieldType, false)
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

			err = s.Set(labelName, val, store.SecretMultilineField, false)
			if err != nil {
				return o, err
			}
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

func gopassSaveSecret(s store.Secret, gp store.Store, p string, dryrun bool, l *log.Entry) error {
	if dryrun {
		return nil
	}

	saved, err := gp.SetIfChanged(s, p)
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
