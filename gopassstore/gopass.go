package gopassstore

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/gopasspw/gopass/pkg/gopass/secrets"
	"github.com/revengel/enpass2gopass/utils"
)

var (
	// ErrNotFound is returned if an entry was not found.
	ErrNotFound = fmt.Errorf("entry is not in the password store")
)

// Gopass -
type Gopass struct {
	ctx context.Context
	api *api.Gopass
}

// Get -
func (g Gopass) Get(p string) (o gopass.Secret, err error) {
	o, err = g.api.Get(g.ctx, p, "latest")
	if err != nil {
		return
	}
	return
}

// Set -
func (g Gopass) Set(s gopass.Byter, p string) (err error) {
	err = g.api.Set(g.ctx, p, s)
	if err != nil {
		return
	}
	return
}

// List -
func (g Gopass) List(keyRe string) (keys []string, err error) {
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
func (g Gopass) Remove(p string) (err error) {
	err = g.api.Remove(g.ctx, p)
	if err != nil {
		return
	}
	return
}

// Close -
func (g *Gopass) Close() error {
	return g.api.Close(g.ctx)
}

// Diff -
func (g Gopass) Diff(a, b gopass.Byter) bool {
	ahash := utils.GetHashFromBytes(a.Bytes())
	bhash := utils.GetHashFromBytes(b.Bytes())
	return ahash == bhash
}

// DiffWithStorage -
func (g Gopass) DiffWithStorage(s gopass.Byter, p string) (bool, error) {
	rSec, err := g.Get(p)
	if err != nil {
		// TODO: need be refactoring
		if err.Error() == ErrNotFound.Error() {
			return false, nil
		}
		return false, err
	}

	return g.Diff(s, rSec), nil
}

// SetIfChanged -
func (g Gopass) SetIfChanged(s gopass.Byter, p string) (bool, error) {
	same, err := g.DiffWithStorage(s, p)
	if err != nil {
		return false, err
	}

	if same {
		return false, nil
	}

	err = g.Set(s, p)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Sync -
func (g Gopass) Sync() error {
	var err = g.api.Sync(g.ctx)
	return err
}

// NewGopass -
func NewGopass(ctx context.Context) (g *Gopass, err error) {
	var gp *api.Gopass
	gp, err = api.New(ctx)
	if err != nil {
		return g, fmt.Errorf("failed to initialize gopass API: %s", err.Error())
	}

	return &Gopass{
		ctx: ctx,
		api: gp,
	}, nil
}

// SecretSet -
func SecretSet(s *secrets.AKV, k, v string, multiline bool) (err error) {
	if k == "" || v == "" {
		return
	}

	if multiline {
		data := []byte(fmt.Sprintf("%s\n%s\n", k, v))
		_, err = s.Write(data)
		return err
	}

	return s.Set(k, v)
}

// SecretAddYamlData -
func SecretAddYamlData(data, k, v string) string {
	if k == "" || v == "" {
		return data
	}

	if data == "" {
		data += "---\n"
	}

	data += fmt.Sprintf("%s\n\n%s\n", k, v)
	return data
}

// GetItemAttachSecret -
func GetItemAttachSecret(flname, dataBase64 string) (o gopass.Byter, err error) {
	var s = secrets.NewAKV()

	err = SecretSet(s, "Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", flname), false)
	if err != nil {
		return s, err
	}

	err = SecretSet(s, "Content-Transfer-Encoding", "Base64", false)
	if err != nil {
		return s, err
	}

	_, err = s.Write([]byte(dataBase64))
	if err != nil {
		return
	}

	return s, nil
}

// SaveSecret -
func SaveSecret(s gopass.Byter, gp *Gopass, p string, dryrun bool) (saved bool, err error) {
	if dryrun {
		return
	}

	saved, err = gp.SetIfChanged(s, p)
	if err != nil {
		return false, err
	}

	return true, nil
}
