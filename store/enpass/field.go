package enpass

import "github.com/revengel/enpass2gopass/utils"

// Field -
type Field struct {
	Deleted   uint8  `json:"deleted"`
	Type      string `json:"type"`
	Sensitive uint8  `json:"sensitive"`
	Label     string `json:"label"`
	Value     string `json:"value"`
}

// IsDeleted -
func (f Field) IsDeleted() bool {
	return f.Deleted == 1
}

// IsSensitive -
func (f Field) IsSensitive() bool {
	return f.Sensitive == 1
}

// IsMultiline -
func (f Field) IsMultiline() bool {
	return f.CheckType("multiline")
}

// CheckType -
func (f Field) CheckType(in string) bool {
	return f.Type == in
}

// CheckTypes -
func (f Field) CheckTypes(in []string) bool {
	for _, i := range in {
		if res := f.CheckType(i); res {
			return true
		}
	}
	return false
}

// GetLabel -
func (f Field) GetLabel() string {
	return utils.Transliterate(f.Label)
}

// GetValue -
func (f Field) GetValue() string {
	return f.Value
}
