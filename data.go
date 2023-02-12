package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
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
		out[folder.UUID] = transliterate(folder.Title)
	}
	return out
}

// loadind data from json file
func loadData(dataPath string) (d Data, err error) {
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
