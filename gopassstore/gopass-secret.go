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
		var data string
		if gs.data == "" {
			data += "---\n"
		}
		data += fmt.Sprintf("%s\n\n%s\n", k, v)
		gs.data += data
		return gs.write(data)
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

// NewEmptyGopassSecret -
func NewEmptyGopassSecret() *GopassSecret {
	return &GopassSecret{
		secret:      secrets.NewAKV(),
		attachments: make(map[string]*secrets.AKV),
	}
}
