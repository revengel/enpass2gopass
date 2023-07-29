package utils

import (
	"fmt"
	"sync"
)

// InsertedPaths -
type InsertedPaths struct {
	sync.Mutex
	data map[string]uint64
}

// Register -
func (i *InsertedPaths) Register(p string) {
	i.Lock()
	defer i.Unlock()

	var val uint64 = 1
	var hash = GetHash(p)
	if v, ok := i.data[hash]; ok {
		val = v + 1
	}
	i.data[hash] = val
}

// Check -
func (i *InsertedPaths) Check(p string) uint64 {
	i.Lock()
	defer i.Unlock()

	var hash = GetHash(p)
	if v, ok := i.data[hash]; ok {
		return v
	}
	return 0
}

// GetUniquePath -
func (i *InsertedPaths) GetUniquePath(p string) (bool, string) {
	if c := i.Check(p); c > 0 {
		newPath := fmt.Sprintf("%s_%d", p, c+1)
		return true, newPath
	}
	return false, p
}

// NewInsertedPaths -
func NewInsertedPaths() *InsertedPaths {
	return &InsertedPaths{
		data: make(map[string]uint64),
	}
}
