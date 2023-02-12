package main

// FoldersMap -
type FoldersMap map[string]string

// GetFolder -
func (f FoldersMap) GetFolder(id string) string {
	if v, ok := f[id]; ok {
		return v
	}
	return ""
}

// GetFolders -
func (f FoldersMap) GetFolders(ids []string) (out []string) {
	for _, id := range ids {
		if f := f.GetFolder(id); id != "" {
			out = append(out, f)
		}
	}
	return out
}

// GetFirstFolder -
func (f FoldersMap) GetFirstFolder(ids []string) string {
	if fs := f.GetFolders(ids); len(fs) > 0 {
		return fs[0]
	}
	return ""
}
