package field

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
	// SecretAttachmentField - secret field multiline type
	SecretAttachmentField = "attachment"
)

type FieldType string

type Field struct {
	Key       string
	Value     []byte
	Type      FieldType
	Multiline bool
	Sensitive bool
}

func (self Field) GetKey() string {
	return self.Key
}

func (self Field) GetValue() []byte {
	return self.Value
}

func (self Field) GetValueString() string {
	return string(self.Value)
}

func (self Field) GetType() FieldType {
	if self.Type == "" {
		return SecretSimpleField
	}
	return self.Type
}

func (self Field) IsType(t FieldType) bool {
	return self.GetType() == t
}

func (self Field) IsMultiline() bool {
	return self.Multiline
}

func (self Field) IsSensitive() bool {
	return self.Sensitive
}

func NewField(k string, v []byte, t FieldType, multi, sens bool) FieldInterface {
	return &Field{
		Key:       k,
		Value:     v,
		Type:      t,
		Multiline: multi,
		Sensitive: sens,
	}
}

func vOrDef(in, def string) string {
	if in != "" {
		return in
	}
	return def
}

func NewTitleField(k, v string) FieldInterface {
	return NewField(vOrDef(k, "title"), []byte(v), SecretTitleField, false, false)
}

func NewSimpleField(k string, v []byte, multi, sens bool) FieldInterface {
	return NewField(k, v, SecretSimpleField, multi, sens)
}

func NewUsernameField(k, v string) FieldInterface {
	return NewField(vOrDef(k, "username"), []byte(v), SecretUsernameField, false, false)
}

func NewUrlField(k, v string) FieldInterface {
	return NewField(vOrDef(k, "url"), []byte(v), SecretUsernameField, false, false)
}

func NewTagsField(k, v string) FieldInterface {
	return NewField(vOrDef(k, "tags"), []byte(v), SecretTagsField, false, false)
}

func NewAttachmentField(k string, v []byte) FieldInterface {
	return NewField(k, v, SecretAttachmentField, false, false)
}

func NewPasswordField(k, v string) FieldInterface {
	return NewField(vOrDef(k, "password"), []byte(v), SecretPasswordField, false, true)
}
