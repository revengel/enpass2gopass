package gopassstore

import (
	"fmt"

	"github.com/gopasspw/gopass/pkg/gopass/secrets"
	"github.com/revengel/enpass2gopass/store"
)

// GopassSecret -
type GopassSecret struct {
	secret    *secrets.AKV
	data      string
	finalized bool
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
	case store.SecretTitileField, store.SecretSimpleField:
		return gs.secret.Set(k, v)
	case store.SecretMultilineField:
		var data string
		if gs.data == "" {
			data += "---\n"
		}
		data += fmt.Sprintf("%s\n\n%s\n", k, v)
		gs.data += data
		return gs.write(data)
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

// Bytes -
func (gs *GopassSecret) Bytes() []byte {
	return gs.secret.Bytes()
}

// NewEmptyGopassSecret -
func NewEmptyGopassSecret() *GopassSecret {
	return &GopassSecret{
		secret: secrets.NewAKV(),
		data:   "",
	}
}

// NewAttachmentGopassSecret -
func NewAttachmentGopassSecret(filename, base64data string) (s *GopassSecret, err error) {
	s = NewEmptyGopassSecret()

	err = s.Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", filename),
		store.SecretSimpleField, false)
	if err != nil {
		return
	}

	err = s.Set("Content-Transfer-Encoding", "Base64", store.SecretSimpleField, false)
	if err != nil {
		return
	}

	err = s.write(base64data)
	if err != nil {
		return
	}

	return
}
