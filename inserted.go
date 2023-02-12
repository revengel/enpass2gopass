package main

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
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
	var hash = getHash(p)
	if v, ok := i.data[hash]; ok {
		val = v + 1
	}
	i.data[hash] = val
}

// Check -
func (i *InsertedPaths) Check(p string) uint64 {
	i.Lock()
	defer i.Unlock()

	var hash = getHash(p)
	if v, ok := i.data[hash]; ok {
		return v
	}
	return 0
}

// GetUniquePath -
func (i *InsertedPaths) GetUniquePath(p string) string {
	if c := i.Check(p); c > 0 {
		newPath := fmt.Sprintf("%s_%d", p, c+1)
		log.Warnf("gopass path '%s' will be rename to '%s'", p, newPath)
		return newPath
	}
	return p
}

func newInsertedPaths() *InsertedPaths {
	return &InsertedPaths{
		data: make(map[string]uint64),
	}
}
