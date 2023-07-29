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
	var err error
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
	case store.SecretSimpleField:
		return gs.secret.Set(k, v)
	case store.SecretMultilineField:
		data := []byte(fmt.Sprintf("%s\n%s\n", k, v))
		_, err = gs.secret.Write(data)
		return err
	case store.SecretYamlField:
		// AddYamlData -
		if gs.data == "" {
			gs.data += "---\n"
		}
		gs.data += fmt.Sprintf("%s\n\n%s\n", k, v)
	default:
		return store.ErrSecretFieldInvalidType
	}

	return nil
}

func (gs *GopassSecret) write(data string) (err error) {
	_, err = gs.secret.Write([]byte(data))
	if err != nil {
		return err
	}
	return
}

// Finalize -
func (gs *GopassSecret) Finalize() error {
	var err error
	if gs.finalized {
		return nil
	}

	err = gs.write(gs.data)
	if err != nil {
		return err
	}

	gs.finalized = true
	return nil
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
