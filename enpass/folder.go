package enpass

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
		if fv := f.GetFolder(id); id != "" && fv != "" {
			out = append(out, fv)
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
