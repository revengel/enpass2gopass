package main

import (
	"fmt"

	"github.com/revengel/enpass2gopass/enpass"
	"github.com/revengel/enpass2gopass/gopassstore"
	"github.com/revengel/enpass2gopass/store"
)

func getGopassItemSecret(item enpass.DataItem) (p string, s store.Secret, err error) {
	var folder = item.GetFirstFolder(foldersMap)
	s = gopassstore.NewEmptyGopassSecret()

	p, err = getGopassPath(folder, item)
	if err != nil {
		return
	}

	err = s.Set("title", item.GetTitle(), store.SecretTitleField, false)
	if err != nil {
		return
	}

	err = s.Set("subtitle", item.GetSubtitle(), store.SecretUsernameField, false)
	if err != nil {
		return
	}

	err = s.Set("category", item.GetCategoryPath(), store.SecretSimpleField, false)
	if err != nil {
		return
	}

	err = s.Set("note", item.GetNote(), store.SecretMultilineField, false)
	if err != nil {
		return
	}

	err = s.Set("tags", item.GetFoldersStr(foldersMap), store.SecretTagsField, false)
	if err != nil {
		return
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
		switch {
		case field.CheckType("password"):
			fieldType = store.SecretPasswordField
		case field.CheckType("url"):
			fieldType = store.SecretURLField
		}

		if field.IsMultiline() {
			fieldType = store.SecretMultilineField
		}

		err = s.Set(labelName, field.GetValue(), fieldType, false)
		if err != nil {
			return
		}
	}

	for _, attach := range item.Attachments {
		isText, err := attach.IsTextData()
		if err != nil {
			return p, s, err
		}

		if isText {
			labelName := fmt.Sprintf("attachment - %s", attach.GetName())
			val, err := attach.GetDataString()
			if err != nil {
				return p, s, err
			}

			err = s.Set(labelName, val, store.SecretMultilineField, false)
			if err != nil {
				return p, s, err
			}
			continue
		}

		var dataB64 = attach.GetDataBase64Encoded()
		var flName = attach.GetNameOriginal()
		err = s.Set(flName, dataB64, store.SecretAttachmentField, false)
		if err != nil {
			return p, s, err
		}
	}

	return p, s, err
}
