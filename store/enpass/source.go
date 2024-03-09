package enpass

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/revengel/enpass2gopass/store"
)

type EnpassSource struct {
	path string
}

func (self EnpassSource) LoadData() (o []store.StoreSourceItem, err error) {
	var d Data
	jsonFile, err := os.Open(self.path)
	if err != nil {
		return
	}

	defer jsonFile.Close()

	b, err := io.ReadAll(jsonFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &d)
	if err != nil {
		return
	}

	var foldersMap = d.GetFoldersMap()
	var items []store.StoreSourceItem
	for _, item := range d.Items {
		var folders = foldersMap.GetFolders(item.Folders)
		item.Folders = folders
		items = append(items, item)
	}

	return items, nil
}

func NewEnpassJsonSource(dataPath string) (o *EnpassSource, err error) {
	absPath, err := filepath.Abs(dataPath)
	if err != nil {
		return
	}

	return &EnpassSource{
		path: absPath,
	}, err
}
