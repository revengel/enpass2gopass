package store

import "errors"

const (
	// SecretPasswordField - secret field password type
	SecretPasswordField = "password"
	// SecretSimpleField - secret field simple text type
	SecretSimpleField = "simple"
	// SecretMultilineField - secret field multiline type
	SecretMultilineField = "multiline"
	// SecretYamlField - secret field yaml type
	SecretYamlField = "yaml"
)

var (
	// ErrSecretFieldInvalidType - error when invalid secret field type
	ErrSecretFieldInvalidType = errors.New("secret field has invalid type")
)

// Secret -
type Secret interface {
	Set(k, v string, fieldType string, sensitivity bool) error
	Finalize() error
	Bytes() []byte
}
