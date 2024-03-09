package enpass

import (
	"github.com/revengel/enpass2gopass/utils"
)

// Data -
type Data struct {
	Folders []FolderItem `json:"folders"`
	Items   []DataItem   `json:"items"`
}

// FolderItem -
type FolderItem struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

// GetFoldersMap -
func (d Data) GetFoldersMap() FoldersMap {
	out := make(map[string]string)
	for _, folder := range d.Folders {
		out[folder.UUID] = utils.Transliterate(folder.Title)
	}
	return out
}
