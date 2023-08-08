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

// Set -
func (s Secret) Set(k, v string, fieldType string, sensitivity bool) error {
	if k == "" || v == "" {
		return nil
	}

	switch fieldType {
	case store.SecretTitileField:
		if t := s.entry.GetTitle(); t != "" {
			return s.setKey("Title", v, false)
		}
		return s.setKey(k, v, sensitivity)
	case store.SecretPasswordField:
		if p := s.entry.GetPassword(); p != "" {
			return s.setKey("Password", v, true)
		}
		return s.setKey(k, v, sensitivity)
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

// Bytes -
func (s Secret) Bytes() []byte {
	return nil
}

// NewSecret -
func NewSecret(db *gokeepasslib.Database) *Secret {
	return &Secret{
		entry: gokeepasslib.NewEntry(),
		db:    db,
	}
}
