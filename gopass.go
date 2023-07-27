package main

import (
	"fmt"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/secrets"
	"github.com/revengel/enpass2gopass/enpass"
	"github.com/revengel/enpass2gopass/gopassstore"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNotFound is returned if an entry was not found.
	ErrNotFound = fmt.Errorf("entry is not in the password store")
)

func getGopassItemSecrets(prefix string, item enpass.DataItem) (map[string]gopass.Byter, error) {
	var (
		err    error
		data   string
		o      = make(map[string]gopass.Byter)
		s      = secrets.NewAKV()
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

	err = gopassstore.SecretSet(s, "title", item.GetTitle(), false)
	if err != nil {
		return o, err
	}

	err = gopassstore.SecretSet(s, "subtitle", item.GetSubtitle(), false)
	if err != nil {
		return o, err
	}

	err = gopassstore.SecretSet(s, "category", item.GetCategoryPath(), false)
	if err != nil {
		return o, err
	}

	data = gopassstore.SecretAddYamlData(data, "note", item.GetNote())

	err = gopassstore.SecretSet(s, "tags", item.GetFoldersStr(foldersMap), false)
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

		if s.Password() == "" && field.CheckType("password") {
			passwd := field.GetValue()
			s.SetPassword(passwd)
			continue
		}

		var labelName = field.GetLabel()
		if field.CheckType("totp") {
			labelName = "totp"
		}

		if labelName == "e_mail" {
			labelName = "email"
		}

		if field.IsMultiline() {
			data = gopassstore.SecretAddYamlData(data, labelName, field.GetValue())
			continue
		}

		err = gopassstore.SecretSet(s, labelName, field.GetValue(), field.IsMultiline())
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
			data = gopassstore.SecretAddYamlData(data, labelName, val)
			continue
		}

		gopassAttachPath, err := getGopassAttachPath(gopassPath, name)
		if err != nil {
			return o, err
		}

		attachSec, err := getGopassItemAttachSecret(attach)
		if err != nil {
			return o, err
		}

		o[gopassAttachPath] = attachSec
	}

	_, err = s.Write([]byte(data))
	if err != nil {
		return o, err
	}

	o[gopassDataPath] = s
	return o, err
}

func getGopassItemAttachSecret(attach enpass.Attachment) (o gopass.Byter, err error) {
	var dataB64 = attach.GetDataBase64Encoded()
	var flName = attach.GetNameOriginal()
	return gopassstore.GetItemAttachSecret(flName, dataB64)
}

func gopassSaveSecret(s gopass.Byter, gp *gopassstore.Gopass, p string, dryrun bool, l *log.Entry) error {
	if dryrun {
		return nil
	}

	saved, err := gopassstore.SaveSecret(s, gp, p, dryrun)
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
