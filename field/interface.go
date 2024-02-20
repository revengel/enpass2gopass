package field

type FieldInterface interface {
	GetKey() string
	GetValue() []byte
	GetValueString() string
	GetType() FieldType
	IsType(t FieldType) bool
	IsMultiline() bool
	IsSensitive() bool
}
