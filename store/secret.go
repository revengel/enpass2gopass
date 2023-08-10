package store

import "errors"

const (
	// SecretTitleField - secret field password type
	SecretTitleField = "title"
	// SecretUsernameField -
	SecretUsernameField = "username"
	// SecretPasswordField - secret field password type
	SecretPasswordField = "password"
	// SecretURLField -
	SecretURLField = "url"
	// SecretTagsField -
	SecretTagsField = "tags"
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
}
