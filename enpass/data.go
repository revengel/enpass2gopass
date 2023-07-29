package enpass

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

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

// LoadData - loadind data from json file
func LoadData(dataPath string) (d Data, err error) {
	absPath, err := filepath.Abs(dataPath)
	if err != nil {
		return
	}

	jsonFile, err := os.Open(absPath)
	if err != nil {
		return
	}

	defer jsonFile.Close()

	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &d)
	if err != nil {
		return
	}

	return
}
