package gopassstore

import (
	"fmt"

	"github.com/gopasspw/gopass/pkg/gopass/secrets"
	"github.com/revengel/enpass2gopass/store"
)

// GopassSecret -
type GopassSecret struct {
	secret      *secrets.AKV
	attachments map[string]*secrets.AKV
	data        string
	finalized   bool
}

// Set -
func (gs *GopassSecret) Set(k, v string, fieldType string, sensitivity bool) error {
	_ = sensitivity
	if k == "" || v == "" {
		return nil
	}

	switch fieldType {
	case store.SecretPasswordField:
		// SetPassword -
		if gs.secret.Password() == "" {
			gs.secret.SetPassword(v)
			return nil
		}
		fallthrough
	case store.SecretTitleField, store.SecretUsernameField, store.SecretURLField, store.SecretTagsField, store.SecretSimpleField:
		return gs.secret.Set(k, v)
	case store.SecretMultilineField:
		if gs.data == "" {
			gs.data += "---\n"
		}
		gs.data += fmt.Sprintf("%s\n\n%s\n", k, v)
		return nil
	case store.SecretAttachmentField:
		var err error
		var secret = secrets.NewAKV()
		err = secret.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", k))
		if err != nil {
			return err
		}

		err = secret.Set("Content-Transfer-Encoding", "Base64")
		if err != nil {
			return err
		}

		_, err = secret.Write([]byte(v))
		if err != nil {
			return err
		}

		gs.attachments[k] = secret
		return nil
	default:
		return store.ErrSecretFieldInvalidType
	}
}

func (gs *GopassSecret) write(data string) (err error) {
	_, err = gs.secret.Write([]byte(data))
	if err != nil {
		return err
	}
	return
}

func (gs *GopassSecret) finalize() (err error) {
	if gs.finalized {
		return
	}

	err = gs.write(gs.data)
	if err != nil {
		return err
	}

	gs.finalized = true
	return nil
}

// NewEmptyGopassSecret -
func NewEmptyGopassSecret() *GopassSecret {
	return &GopassSecret{
		secret:      secrets.NewAKV(),
		attachments: make(map[string]*secrets.AKV),
	}
}
