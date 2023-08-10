package gopassstore

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/revengel/enpass2gopass/store"
	"github.com/revengel/enpass2gopass/utils"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNotFound is returned if an entry was not found.
	ErrNotFound = fmt.Errorf("entry is not in the password store")
)

// Gopass -
type Gopass struct {
	ctx    context.Context
	api    *api.Gopass
	prefix string
	items  *utils.UniqueStrings
}

// Get -
func (g Gopass) get(p string) (o gopass.Secret, err error) {
	return g.api.Get(g.ctx, p, "latest")
}

// Set -
func (g Gopass) set(s gopass.Byter, p string) (err error) {
	return g.api.Set(g.ctx, p, s)
}

func (g Gopass) list(keyRe string) (keys []string, err error) {
	keys, err = g.api.List(g.ctx)
	if err != nil {
		return
	}

	if keyRe == "" {
		return
	}

	gopassKeyRe := regexp.MustCompile(keyRe)
	var filteredKeys []string
	for _, k := range keys {
		if !gopassKeyRe.MatchString(k) {
			continue
		}
		filteredKeys = append(filteredKeys, k)
	}

	return filteredKeys, nil
}

// Remove -
func (g Gopass) remove(p string) error {
	return g.api.Remove(g.ctx, p)
}

// Close -
func (g *Gopass) Close() error {
	return g.api.Close(g.ctx)
}

// Diff -
func (g Gopass) diff(a, b gopass.Byter) bool {
	ahash := utils.GetHashFromBytes(a.Bytes())
	bhash := utils.GetHashFromBytes(b.Bytes())
	return ahash == bhash
}

// DiffWithStorage -
func (g Gopass) diffWithStorage(s gopass.Byter, p string) (bool, error) {
	rSec, err := g.get(p)
	if err != nil {
		// TODO: need be refactoring
		if err.Error() == ErrNotFound.Error() {
			return false, nil
		}
		return false, err
	}

	return g.diff(s, rSec), nil
}

func (g Gopass) saveSecret(s gopass.Byter, p string) (bool, error) {
	p = g.items.Unique(p)
	var l = log.WithField("gopasskey", p)

	same, err := g.diffWithStorage(s, p)
	if err != nil {
		return false, err
	}

	if same {
		l.Debug("gopass secret already in actual state")
		return false, nil
	}

	err = g.set(s, p)
	if err != nil {
		return false, err
	}

	l.Info("secret has been updated")
	return true, nil
}

// Cleanup -
func (g Gopass) Cleanup() (bool, error) {
	var deletesCount = 0
	ll, err := g.list(`^` + g.prefix + `/`)
	if err != nil {
		return false, err
	}

	for _, k := range ll {
		if g.items.Has(k) {
			continue
		}

		var lc = log.WithField("type", "cleaner").
			WithField("gopasskey", k)

		lc.Info("gopass key will be deleted")
		err = g.remove(k)
		if err != nil {
			return false, err
		}

		deletesCount++
	}

	return deletesCount > 0, nil
}

func (g Gopass) getMainSecretPath(p string) string {
	return filepath.Join(g.prefix, p, "data")
}

func (g Gopass) getAttachmentSecretPath(p, attachmentName string) string {
	return filepath.Join(g.prefix, p, "attachments", attachmentName)
}

// Save -
func (g Gopass) Save(s store.Secret, p string) (bool, error) {
	var keyPath = g.getMainSecretPath(p)
	var out bool
	var secret *GopassSecret

	switch v := s.(type) {
	case *GopassSecret:
		secret = v
	default:
		return false, errors.New("secret must be *GopassSecret")
	}

	same, err := g.saveSecret(secret.secret, keyPath)
	if err != nil {
		return out, err
	}

	out = out || same

	for attachName, secret := range secret.attachments {
		var keyPath = g.getAttachmentSecretPath(p, attachName)
		same, err := g.saveSecret(secret, keyPath)
		if err != nil {
			return out, err
		}

		out = out || same
	}

	return out, nil
}

// Sync -
func (g Gopass) Sync() error {
	return g.api.Sync(g.ctx)
}

// NewStore -
func NewStore(ctx context.Context, prefix string) (g *Gopass, err error) {
	var gp *api.Gopass
	gp, err = api.New(ctx)
	if err != nil {
		return g, fmt.Errorf("failed to initialize gopass API: %s", err.Error())
	}

	if prefix == "" {
		prefix = "enpass"
	}

	return &Gopass{
		ctx:    ctx,
		api:    gp,
		prefix: prefix,
		items:  utils.NewUniqueStrings(),
	}, nil
}
