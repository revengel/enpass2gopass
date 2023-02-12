package main

import (
	"fmt"
	"strings"
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
	return transliterate(i.Category)
}

// GetTitle -
func (i DataItem) GetTitle() string {
	return i.Title
}

// GetTitlePath -
func (i DataItem) GetTitlePath() string {
	return transliterate(i.Title)
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
	return foldersMap.GetFolders(i.Folders)
}

// GetFirstFolder -
func (i DataItem) GetFirstFolder() string {
	if fs := i.GetFolders(); len(fs) > 0 {
		return fs[0]
	}
	return ""
}

// GetFoldersStr -
func (i DataItem) GetFoldersStr() string {
	return fmt.Sprintf("[%s]", strings.Join(i.GetFolders(), ", "))
}
