package keepassstore

import (
	"encoding/base64"
	"fmt"

	"github.com/revengel/enpass2gopass/store"
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

// Secret -
type Secret struct {
	entry gokeepasslib.Entry
	db    *gokeepasslib.Database
}

func (s Secret) setKey(k, v string, sensitivity bool) error {
	s.entry.Values = append(s.entry.Values, gokeepasslib.ValueData{
		Key: k,
		Value: gokeepasslib.V{
			Content:   v,
			Protected: wrappers.NewBoolWrapper(sensitivity),
		},
	})
	return nil
}

func (s Secret) setKeyOrAlt(k, altK, v string, sensitivity bool) error {
	if t := s.entry.GetContent(k); t == "" {
		return s.setKey(k, v, sensitivity)
	}
	return s.setKey(altK, v, sensitivity)
}

// Set -
func (s Secret) Set(k, v string, fieldType string, sensitivity bool) error {
	if k == "" || v == "" {
		return nil
	}

	switch fieldType {
	case store.SecretTitleField:
		return s.setKeyOrAlt("Title", k, v, false)
	case store.SecretUsernameField:
		return s.setKeyOrAlt("Username", k, v, sensitivity)
	case store.SecretPasswordField:
		return s.setKeyOrAlt("Password", k, v, true)
	case store.SecretURLField:
		return s.setKeyOrAlt("URL", k, v, false)
	case store.SecretTagsField:
		return s.setKeyOrAlt("Tags", k, v, false)
	case store.SecretMultilineField:
		var notes = s.entry.GetContent("Notes")
		notes += fmt.Sprintf("%s\n\n%s\n", k, v)
		return s.setKey("Notes", notes, false)
	case store.SecretSimpleField:
		return s.setKey(k, v, sensitivity)
	case store.SecretAttachmentField:
		var dataBytes []byte
		dataBytes, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return err
		}

		var bin = s.db.AddBinary(dataBytes)
		s.entry.Binaries = append(s.entry.Binaries, bin.CreateReference(k))
		return nil
	default:
		return store.ErrSecretFieldInvalidType
	}
}

// NewSecret -
func NewSecret(db *gokeepasslib.Database) *Secret {
	return &Secret{
		entry: gokeepasslib.NewEntry(),
		db:    db,
	}
}
