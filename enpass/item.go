package enpass

import (
	"fmt"
	"strings"

	"github.com/revengel/enpass2gopass/field"
	"github.com/revengel/enpass2gopass/utils"
)

// DataItem -
type DataItem struct {
	Trashed  uint8 `json:"trashed"`
	Archived uint8 `json:"archived"`
	Favorite uint8 `json:"favorite"`

	Category string `json:"category"`

	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Note     string `json:"note"`

	Fields      []Field      `json:"fields"`
	Folders     []string     `json:"folders"`
	Attachments []Attachment `json:"attachments"`
}

// IsTrashed -
func (i DataItem) IsTrashed() bool {
	return i.Trashed == 1
}

// IsArchived -
func (i DataItem) IsArchived() bool {
	return i.Archived == 1
}

// IsFavorite -
func (i DataItem) IsFavorite() bool {
	return i.Favorite == 1
}

// GetCategory -
func (i DataItem) GetCategory() string {
	return i.Category
}

// GetCategoryPath -
func (i DataItem) GetCategoryPath() string {
	return utils.Transliterate(i.Category)
}

// GetTitle -
func (i DataItem) GetTitle() string {
	return i.Title
}

// GetTitlePath -
func (i DataItem) GetTitlePath() string {
	return utils.Transliterate(i.Title)
}

// GetSubtitle -
func (i DataItem) GetSubtitle() string {
	return i.Subtitle
}

// GetNote -
func (i DataItem) GetNote() string {
	return i.Note
}

// GetFolders -
func (i DataItem) GetFolders() []string {
	return i.Folders
}

// GetFirstFolder -
func (i DataItem) GetFirstFolder() string {
	if len(i.Folders) == 0 {
		return ""
	}
	return i.Folders[0]
}

// GetFoldersStr -
func (i DataItem) GetFoldersStr() string {
	return fmt.Sprintf("[%s]", strings.Join(i.GetFolders(), ", "))
}

// GetFields -
func (i DataItem) GetFields() (out []field.FieldInterface, err error) {
	out = append(out, field.NewTitleField("", i.GetTitle()))
	out = append(out, field.NewUsernameField("subtitle", i.GetSubtitle()))

	out = append(out, field.NewSimpleField("category", []byte(i.GetCategoryPath()), false, false))
	out = append(out, field.NewSimpleField("note", []byte(i.GetCategoryPath()), true, false))

	out = append(out, field.NewTagsField("", i.GetFoldersStr()))

	for _, f := range i.Fields {
		var ignoreTypes = []string{"section", ".Android#"}
		if f.IsDeleted() || f.CheckTypes(ignoreTypes) {
			continue
		}

		if f.GetLabel() == "" || f.GetValue() == "" {
			continue
		}

		var labelName = f.GetLabel()
		if f.CheckType("totp") {
			labelName = "totp"
		}

		if labelName == "e_mail" {
			labelName = "email"
		}

		var fieldType field.FieldType = field.SecretSimpleField
		switch {
		case f.CheckType("password"):
			fieldType = field.SecretPasswordField
		case f.CheckType("url"):
			fieldType = field.SecretURLField
		}

		ff := field.NewField(labelName, []byte(f.GetValue()), fieldType, f.IsMultiline(), false)
		out = append(out, ff)
	}

	for _, attach := range i.Attachments {
		isText, err := attach.IsTextData()
		if err != nil {
			return nil, err
		}

		if isText {
			labelName := fmt.Sprintf("attachment - %s", attach.GetName())
			val, err := attach.GetDataBytes()
			if err != nil {
				return nil, err
			}

			ff := field.NewSimpleField(labelName, val, true, false)
			out = append(out, ff)
			continue
		}

		var flName = attach.GetNameOriginal()
		val, err := attach.GetDataBytes()
		if err != nil {
			return nil, err
		}

		ff := field.NewAttachmentField(flName, val)
		out = append(out, ff)
	}

	return
}
