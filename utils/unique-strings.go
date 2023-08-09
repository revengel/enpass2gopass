package utils

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// UniqueStrings -
type UniqueStrings struct {
	sync.Mutex
	counts map[string]uint64
	unique map[string]uint64
}

// Unique -
func (u *UniqueStrings) Unique(s string) string {
	u.Lock()
	defer u.Unlock()

	var out = s
	var val uint64 = 1
	var hash = GetHash(s)

	if v, ok := u.counts[hash]; ok {
		val = v + 1
		out = fmt.Sprintf("%s_%d", s, val)
		log.Warnf("string '%s' will be rename to '%s'", s, out)
	}

	u.counts[hash] = val
	u.unique[GetHash(out)] = 1
	return out
}

// Has -
func (u *UniqueStrings) Has(s string) bool {
	u.Lock()
	defer u.Unlock()

	_, ok := u.unique[GetHash(s)]
	return ok
}

// NewUniqueStrings -
func NewUniqueStrings() *UniqueStrings {
	return &UniqueStrings{
		counts: make(map[string]uint64),
		unique: make(map[string]uint64),
	}
}
