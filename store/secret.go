package store

import "errors"

const (
	// SecretTitileField - secret field password type
	SecretTitileField = "title"
	// SecretPasswordField - secret field password type
	SecretPasswordField = "password"
	// SecretSimpleField - secret field simple text type
	SecretSimpleField = "simple"
	// SecretMultilineField - secret field multiline type
	SecretMultilineField = "multiline"
	// SecretAttachmentField - secret field multiline type
	SecretAttachmentField = "attachment"
)

var (
	// ErrSecretFieldInvalidType - error when invalid secret field type
	ErrSecretFieldInvalidType = errors.New("secret field has invalid type")
)

// Secret -
type Secret interface {
	Set(k, v string, fieldType string, sensitivity bool) error
	Bytes() []byte
}
