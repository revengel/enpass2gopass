package gopassstore

import (
	"fmt"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/secrets"
)

// GopassSecret -
type GopassSecret struct {
	secret    *secrets.AKV
	data      string
	finalized bool
}

// Set -
func (gs *GopassSecret) Set(k, v string, multiline bool) (err error) {
	if k == "" || v == "" {
		return
	}

	if multiline {
		data := []byte(fmt.Sprintf("%s\n%s\n", k, v))
		_, err = gs.secret.Write(data)
		return err
	}
	return gs.secret.Set(k, v)
}

// SetPassword -
func (gs *GopassSecret) SetPassword(v string) bool {
	if gs.secret.Password() != "" {
		return false
	}
	gs.secret.SetPassword(v)
	return true
}

// AddYamlData -
func (gs *GopassSecret) AddYamlData(k, v string) {
	if k == "" || v == "" {
		return
	}

	if gs.data == "" {
		gs.data += "---\n"
	}

	gs.data += fmt.Sprintf("%s\n\n%s\n", k, v)
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
	return
}

// GetSecret -
func (gs *GopassSecret) GetSecret() (s gopass.Byter, err error) {
	err = gs.finalize()
	if err != nil {
		return nil, err
	}

	return gs.secret, nil
}

// Bytes -
func (gs *GopassSecret) Bytes() ([]byte, error) {
	s, err := gs.GetSecret()
	if err != nil {
		return nil, err
	}

	return s.Bytes(), nil
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
		false)
	if err != nil {
		return
	}

	err = s.Set("Content-Transfer-Encoding", "Base64", false)
	if err != nil {
		return
	}

	err = s.write(base64data)
	if err != nil {
		return
	}

	return
}
